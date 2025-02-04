# server.py
from concurrent import futures
import grpc
from proto import movie_pb2, movie_pb2_grpc
from sentence_transformers import SentenceTransformer
import numpy as np
import faiss

# Set the embedding dimension for 'all-MiniLM-L6-v2'.
EMBEDDING_DIM = 384

# Create a global Faiss index for L2 similarity.
index = faiss.IndexFlatL2(EMBEDDING_DIM)

# A global dictionary mapping movie IDs to their index positions in the Faiss index.
movie_id_to_index = {}

# A simple counter to keep track of the next index position.
next_index = 0

# Load your sentence transformer model.
model = SentenceTransformer('all-MiniLM-L6-v2')

class EmbeddingService(movie_pb2_grpc.EmbeddingServiceServicer):
    def GetMovieEmbedding(self, request, context):
        # try:
        #     # Immediately log the incoming raw request (using repr)
        #     print("Received raw request:", repr(request))
        #     print(request.movie_id)
        #     raise Exception("movie_id")
        # except Exception as e:
        #     print("Exception during request logging:", e)
        #     raise
        print(f"Received request: movie_id={request.movie_id}, title={request.title}")
        global next_index

        # Create a text string from movie details.
        text = f"{request.title}: {request.overview}"
        
        # Generate the embedding as a NumPy array.
        embedding = model.encode(text, convert_to_numpy=True)
        
        # Convert the embedding to a list so it can be returned in the response.
        embedding_list = embedding.tolist()
        
        # Add the embedding to the Faiss index.
        # Faiss expects a 2D array, so we wrap the embedding in another array.
        embedding_np = np.array([embedding], dtype=np.float32)
        index.add(embedding_np)
        
        # Use the provided movie_id (assumed to be in the request) as the key.
        movie_id = request.movie_id  
        movie_id_to_index[movie_id] = next_index
        
        # Increment our counter for the next index.
        next_index += 1

        # Optionally, print some debug info.
        print(f"Added movie {movie_id} at index position {movie_id_to_index[movie_id]}. Total movies in index: {index.ntotal}")
        
        return movie_pb2.EmbeddingResponse(embedding=embedding_list)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    movie_pb2_grpc.add_EmbeddingServiceServicer_to_server(EmbeddingService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("Server started. Listening on port 50051.")
    server.wait_for_termination()

if __name__ == '__main__':
    serve()
