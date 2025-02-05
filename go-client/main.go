package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	pb "movie-recommender/go-client/pb/proto"
	"net/http"
	"net/url"
	"os"
	"strings"

	"google.golang.org/grpc"
)

var apiKey = os.Getenv("TMDB_API_KEY")

type Keyword struct {
	ID      int    `json:"id"`
	Keyword string `json:"name"`
}

type KeywordResponse struct {
	ID       int       `json:"id"`
	Keywords []Keyword `json:"keywords"`
}

type MovieSearchResponse struct {
	Results []MovieResult `json:"results"`
}

type TopRatedResponse struct {
	Results    []MovieResult `json:"results"`
	TotalPages int           `json:"total_pages"`
}

type MovieResult struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
}

type Movie struct {
	Title    string `json:"title"`
	Overview string `json:"overview"`
	Keywords string `json:"keywords"`
}

// getMovieID performs a search for a movie title on TMDb, and returns the ID of the
// first result. If no results are found, it returns an error.
func getMovieID(title string) (int, error) {
	baseURL := "https://api.themoviedb.org/3/search/movie"
	query := url.QueryEscape(title)
	apiURL := fmt.Sprintf("%s?api_key=%s&query=%s", baseURL, apiKey, query)

	resp, err := http.Get(apiURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var searchResponse MovieSearchResponse
	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		return 0, err
	}

	if len(searchResponse.Results) == 0 {
		return 0, fmt.Errorf("no results found for movie title: %s", title)
	}

	// Return the first result's ID
	return searchResponse.Results[0].ID, nil
}

// getJSON is a helper function that performs a GET request to the specified URL,
// reads the response, and unmarshals the JSON into the target interface.
func getJSON(url string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}

// getMovieInfo retrieves the movie information and its keywords,
// then sets the Movie.Keywords field to the top 5 keywords joined by commas.
func getMovieInfo(movie_id int) (Movie, error) {
	var movie Movie
	// Fetch movie details.
	movieURL := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d?api_key=%s", movie_id, apiKey)
	if err := getJSON(movieURL, &movie); err != nil {
		return Movie{}, err
	}

	// Fetch keywords.
	var keywordResp KeywordResponse
	keywordURL := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d/keywords?api_key=%s", movie_id, apiKey)
	if err := getJSON(keywordURL, &keywordResp); err != nil {
		return movie, err
	}

	// Extract the top 5 keywords.
	topCount := 5
	if len(keywordResp.Keywords) < topCount {
		topCount = len(keywordResp.Keywords)
	}
	topKeywords := make([]string, 0, topCount)
	for i := 0; i < topCount; i++ {
		topKeywords = append(topKeywords, keywordResp.Keywords[i].Keyword)
	}

	// Join the keywords into a comma-separated string.
	movie.Keywords = strings.Join(topKeywords, ", ")

	return movie, nil
}

func getTopRatedMovies() ([]int, error) {

	// First, get the first page to determine the total number of pages.
	movieURL := fmt.Sprintf("https://api.themoviedb.org/3/movie/top_rated?api_key=%s", apiKey)
	var topRated TopRatedResponse

	if err := getJSON(movieURL, &topRated); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total Pages: %d\n", topRated.TotalPages)
	// Initialize a slice to hold all movie IDs.
	allMovieIDs := []int{}

	// Process the first page.
	for _, movie := range topRated.Results {
		allMovieIDs = append(allMovieIDs, movie.ID)
	}

	// Now iterate through pages 2 to TotalPages.
	for page := 2; page <= topRated.TotalPages; page++ {
		pageURL := fmt.Sprintf("https://api.themoviedb.org/3/movie/top_rated?api_key=%s&page=%d", apiKey, page)
		var pageResp TopRatedResponse

		if err := getJSON(pageURL, &pageResp); err != nil {
			log.Fatalf("Error fetching page %d: %v", page, err)
		}

		for _, movie := range pageResp.Results {
			allMovieIDs = append(allMovieIDs, movie.ID)
		}
	}

	return allMovieIDs, nil
}

func buildIndex(client pb.EmbeddingServiceClient) error {
	var allMovieIDs []int
	allMovieIDs, err := getTopRatedMovies()

	if err != nil {
		log.Fatal(err)
	}
	// Now print out all collected movie IDs.
	fmt.Println("Collected Top Rated Movie IDs:")

	for _, movie_id := range allMovieIDs {
		movie, err := getMovieInfo(movie_id)

		if err != nil {
			log.Fatal(err)
		}

		req := &pb.MovieRequest{
			MovieId:  int32(movie_id),
			Title:    movie.Title,
			Overview: movie.Overview,
			Keywords: movie.Keywords,
		}

		log.Printf("Sending MovieRequest: MovieId=%d, Title=%q, Overview=%q, Keywords=%q",
			req.MovieId, req.Title, req.Overview, req.Keywords)

		_, err = client.GetMovieEmbedding(context.Background(), req)

		if err != nil {
			log.Fatalf("Error calling gRPC: %v", err)
		}
	}
	return nil
}

var (
	grpcClient pb.EmbeddingServiceClient
)

func main() {
	// Optional flag for your existing "build the index" operation
	buildFlag := flag.Bool("build", false, "Build the index in the vector DB")
	flag.Parse()

	// Connect to gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	grpcClient = pb.NewEmbeddingServiceClient(conn)

	// Optionally build the index (if user runs with --build)
	if *buildFlag {
		if err := buildIndex(grpcClient); err != nil {
			log.Fatal(err)
		}
	}

	// Set up HTTP routes
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/search", handleSearch)

	fmt.Println("Web server running on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// handleIndex just shows a basic HTML form
func handleIndex(w http.ResponseWriter, r *http.Request) {
	html := `
	<html>
		<head><title>Movie Recommender</title></head>
		<body>
			<h1>Movie Recommender</h1>
			<form action="/search" method="get">
				<input type="text" name="q" placeholder="Enter a movie title" />
				<input type="submit" value="Search" />
			</form>
		</body>
	</html>
	`
	w.Write([]byte(html))
}

// handleSearch handles the /search route.
//
// It takes the user's query from the URL, looks up the first matching TMDb movie ID,
// calls the gRPC server to get recommendations, and then builds a simple HTML page
// with the results.
func handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if strings.TrimSpace(query) == "" {
		http.Error(w, "Please provide a movie title.", http.StatusBadRequest)
		return
	}

	// Find the first matching TMDb movie ID
	movieID, err := getMovieID(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error looking up movie: %v", err), http.StatusInternalServerError)
		return
	}

	similarResp, err := grpcClient.GetSimilarMovies(context.Background(), &pb.SimilarMoviesRequest{
		MovieId: int32(movieID),
		Limit:   5, // we want 5 recommendations
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("gRPC call failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Build a simple results page
	var sb strings.Builder
	sb.WriteString("<html><head><title>Recommendations</title></head><body>")
	sb.WriteString(fmt.Sprintf("<h2>Recommendations for '%s'</h2>", query))
	sb.WriteString("<ul>")
	for _, rec := range similarResp.Recommendations {
		sb.WriteString(fmt.Sprintf("<li>%s (Movie ID: %d)</li>", rec.Title, rec.MovieId))
	}
	sb.WriteString("</ul>")
	sb.WriteString(`<a href="/">Search again</a>`)
	sb.WriteString("</body></html>")

	w.Write([]byte(sb.String()))
}
