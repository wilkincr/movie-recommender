# server.py
from concurrent import futures
import grpc
import movie_pb2
import movie_pb2_grpc
from sentence_transformers import SentenceTransformer

# Load your sentence transformer model
model = SentenceTransformer('all-MiniLM-L6-v2')

class EmbeddingService(movie_pb2_grpc.EmbeddingServiceServicer):
    def GetMovieEmbedding(self, request, context):
        # Create a text string from movie details
        text = f"{request.title}: {request.overview}"
        # Generate embedding (a list of floats)
        embedding = model.encode(text, convert_to_numpy=True).tolist()
        return movie_pb2.EmbeddingResponse(embedding=embedding)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    movie_pb2_grpc.add_EmbeddingServiceServicer_to_server(EmbeddingService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("Server started. Listening on port 50051.")
    server.wait_for_termination()

if __name__ == '__main__':
    serve()
