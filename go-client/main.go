package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	pb "movie-recommender/go-client/pb"

	"google.golang.org/grpc"
)

const apiKey = "c5baf44cf7c6de22b5d2e5c09cfbb255"

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

func getMovieInfo(movie_id int) (Movie, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d?api_key=%s", movie_id, apiKey)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var movie Movie

	err = json.Unmarshal(body, &movie)
	if err != nil {
		log.Fatal(err)
	}

	return movie, nil
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

	// Encode payload to JSON
	// reqBody, err := json.Marshal(movie)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewEmbeddingServiceClient(conn)
	res, err := client.GetMovieEmbedding(context.Background(), &pb.MovieRequest{
		Title:    "Dune",
		Overview: "A mythic and emotionally charged hero’s journey...",
	})
	if err != nil {
		log.Fatalf("Error calling gRPC: %v", err)
	}

	log.Printf("Received embedding: %v", res.Embedding)

}
