from concurrent import futures
import grpc
from proto import movie_pb2, movie_pb2_grpc
from sentence_transformers import SentenceTransformer
import numpy as np
import faiss

# Embedding model and Faiss setup
EMBEDDING_DIM = 384
index = faiss.IndexFlatL2(EMBEDDING_DIM)

# Data storage
movie_id_to_index = {}   # Maps movie_id to Faiss index position.
movie_id_to_title = {}   # Maps movie_id to title.
embeddings_list = []     # Stores embeddings in order.
next_index = 0           # Tracks the next index position.

# Load sentence transformer model
model = SentenceTransformer('all-MiniLM-L6-v2')

class EmbeddingService(movie_pb2_grpc.EmbeddingServiceServicer):
    def GetMovieEmbedding(self, request, context):
        global next_index

        # Generate the embedding from the title and overview
        text = f"{request.title}: {request.overview}"
        embedding = model.encode(text, convert_to_numpy=True)

        # Add the embedding to the Faiss index
        embedding_np = np.array([embedding], dtype=np.float32)
        index.add(embedding_np)

        # Store information
        movie_id = request.movie_id
        movie_id_to_index[movie_id] = next_index
        movie_id_to_title[movie_id] = request.title  # Store movie title
        embeddings_list.append(embedding)
        next_index += 1

        print(f"Added movie {movie_id} ({request.title}) at index {next_index - 1}. Total: {index.ntotal}")
        return movie_pb2.EmbeddingResponse(embedding=embedding.tolist())

    def GetSimilarMovie(self, request, context):
        query_movie_id = request.movie_id
        if query_movie_id not in movie_id_to_index:
            context.set_details("Movie ID not found in index")
            context.set_code(grpc.StatusCode.NOT_FOUND)
            return movie_pb2.SimilarMovieResponse()

        query_index = movie_id_to_index[query_movie_id]
        query_embedding = np.array([embeddings_list[query_index]], dtype=np.float32)

        # Search for 2 nearest neighbors (1st is itself)
        distances, indices = index.search(query_embedding, 2)
        if indices.shape[1] < 2:
            context.set_details("Not enough neighbors found")
            context.set_code(grpc.StatusCode.NOT_FOUND)
            return movie_pb2.SimilarMovieResponse()

        # Second result is the most similar movie
        similar_index = int(indices[0][1])
        similar_movie_id = next((mid for mid, idx in movie_id_to_index.items() if idx == similar_index), None)

        if similar_movie_id is None:
            context.set_details("Similar movie not found")
            context.set_code(grpc.StatusCode.NOT_FOUND)
            return movie_pb2.SimilarMovieResponse()

        similar_title = movie_id_to_title.get(similar_movie_id, "")
        print(f"Most similar movie to {query_movie_id}: {similar_movie_id} - {similar_title}")

        return movie_pb2.SimilarMovieResponse(movie_id=similar_movie_id, title=similar_title)

    def GetSimilarMovies(self, request, context):
        query_movie_id = request.movie_id
        if query_movie_id not in movie_id_to_index:
            context.set_details("Movie ID not found in index")
            context.set_code(grpc.StatusCode.NOT_FOUND)
            return movie_pb2.SimilarMoviesResponse()

        query_index = movie_id_to_index[query_movie_id]
        query_embedding = np.array([embeddings_list[query_index]], dtype=np.float32)

        # Search for `limit + 1` nearest neighbors (excluding itself)
        num_neighbors = request.limit + 1
        distances, indices = index.search(query_embedding, num_neighbors)

        similar_movies = []
        for i in range(1, min(num_neighbors, indices.shape[1])):  # Start from 1 to skip itself
            similar_index = int(indices[0][i])
            similar_movie_id = next((mid for mid, idx in movie_id_to_index.items() if idx == similar_index), None)

            if similar_movie_id is not None:
                similar_movies.append(movie_pb2.SimilarMovieResponse(
                    movie_id=similar_movie_id,
                    title=movie_id_to_title.get(similar_movie_id, "Unknown")
                ))

        print(f"Top {len(similar_movies)} similar movies to {query_movie_id}: {[m.movie_id for m in similar_movies]}")
        return movie_pb2.SimilarMoviesResponse(recommendations=similar_movies)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    movie_pb2_grpc.add_EmbeddingServiceServicer_to_server(EmbeddingService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("Server started on port 50051.")
    server.wait_for_termination()

if __name__ == '__main__':
    serve()
