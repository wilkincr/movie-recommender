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

// MovieSearchResponse represents the TMDb search movie API response
type MovieSearchResponse struct {
	Results []MovieResult `json:"results"`
}

// MovieResult represents an individual movie result
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

// MovieID captures only the movie ID.
type MovieID struct {
	ID int `json:"id"`
}

// TopRatedResponse represents the JSON structure returned from the top-rated movies endpoint.
type TopRatedResponse struct {
	Page         int       `json:"page"`
	Results      []MovieID `json:"results"`
	TotalPages   int       `json:"total_pages"`
	TotalResults int       `json:"total_results"`
}

// EmbeddingResponse represents the response from the Python API.
type EmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

// getMovieID fetches the movie ID from TMDb given a title
func getMovieID(title string) (int, error) {
	baseURL := "https://api.themoviedb.org/3/search/movie"
	query := url.QueryEscape(title) // Encode spaces correctly

	// Construct the API request URL
	apiURL := fmt.Sprintf("%s?api_key=%s&query=%s", baseURL, apiKey, query)

	// Make the request
	resp, err := http.Get(apiURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	// Parse JSON response
	var searchResponse MovieSearchResponse
	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		return 0, err
	}

	// If no results found, return error
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

func main() {

	movie_id, err := getMovieID("Dune")

	if err != nil {
		log.Fatal(err)
	}

	movie, err := getMovieInfo(movie_id)

	if err != nil {
		log.Fatal(err)
	}

	// Print the parsed struct
	fmt.Printf("Title: %s\n", movie.Title)
	fmt.Printf("Overview: %s\n", movie.Overview)
	fmt.Printf("Keywords: %s\n", movie.Keywords)

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewEmbeddingServiceClient(conn)

	// Build the index
	buildFlag := flag.Bool("build", false, "Build the index")
	flag.Parse()

	if *buildFlag {
		err = buildIndex(client)
		if err != nil {
			log.Fatal(err)
		}
	}
	if err != nil {
		log.Fatal(err)
	}

	similarResp, err := client.GetSimilarMovie(context.Background(), &pb.SimilarMovieRequest{
		MovieId: int32(157336), // set queryMovieID to the desired movie ID
	})
	if err != nil {
		log.Fatalf("Error calling GetSimilarMovie: %v", err)
	}
	movie, err = getMovieInfo(int(similarResp.MovieId))

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Most similar movie to %d is %d: %s", movie.Title, similarResp.MovieId, movie.Title)

	// Optionally, if your service requires an explicit insertion into the vector database,
	// 	you might call an AddMovieEmbedding RPC here. For example:

	// conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalf("Failed to connect to gRPC server: %v", err)
	// }
	// defer conn.Close()

	// client := pb.NewEmbeddingServiceClient(conn)
	// result, errora := client.GetMovieEmbedding(context.Background(), &pb.MovieRequest{
	// 	Title:    "Dune",
	// 	Overview: "A mythic and emotionally charged hero’s journey...",
	// })

	// if err != nil {
	// 	log.Fatalf("Error calling gRPC: %v", errora)
	// }

	// log.Printf("Received embedding: %v", result.Embedding)

}
