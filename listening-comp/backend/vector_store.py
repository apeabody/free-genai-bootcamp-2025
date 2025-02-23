import json
import os
from typing import List, Dict

import chromadb
from chromadb.utils import embedding_functions


class QuestionVectorStore:
    def __init__(self, persist_directory: str = "./chroma_db"):
        self.persist_directory = persist_directory
        # Ensure the directory exists
        os.makedirs(persist_directory, exist_ok=True)

        # Create a persistent client
        self.client = chromadb.PersistentClient(path=persist_directory)

        # Use the default embedding function
        self.embedding_function = embedding_functions.DefaultEmbeddingFunction()

        # Create or get the collection
        self.collection = self.client.get_or_create_collection(
            name="spanish_questions", embedding_function=self.embedding_function
        )

    def add_questions(self, questions_file: str) -> None:
        """Add questions from a JSON file to the vector store."""
        with open(questions_file, "r") as f:
            questions = json.load(f)

        # Prepare the documents and metadata
        documents = []
        metadatas = []
        ids = []

        for i, q in enumerate(questions):
            # Create a combined text for embedding
            text = (
                f"Conversation: {q['Conversation']} "
                f"Question: {q['Question']} "
                f"Answer: {q['Answer']}"
            )
            documents.append(text)

            # Store the original data as metadata
            metadatas.append(
                {
                    "conversation": q["Conversation"],
                    "question": q["Question"],
                    "answer": q["Answer"],
                }
            )

            # Create a unique ID
            ids.append(f"q_{i}")

        # Add to the collection
        self.collection.add(documents=documents, metadatas=metadatas, ids=ids)

    def generate_similar_questions(self, query: str, n_results: int = 5) -> List[Dict]:
        """Find similar questions based on the query."""
        results = self.collection.query(query_texts=[query], n_results=n_results)

        # Format the results
        similar_questions = []
        for i in range(len(results["ids"][0])):
            metadata = results["metadatas"][0][i]
            similar_questions.append(
                {
                    "Conversation": metadata["conversation"],
                    "Question": metadata["question"],
                    "Answer": metadata["answer"],
                }
            )

        return similar_questions

    def get_all_questions(self) -> List[Dict]:
        """Retrieve all questions from the collection."""
        results = self.collection.get()

        questions = []
        for i in range(len(results["ids"])):
            metadata = results["metadatas"][i]
            questions.append(
                {
                    "Conversation": metadata["conversation"],
                    "Question": metadata["question"],
                    "Answer": metadata["answer"],
                }
            )
        return questions

    def get_all_conversations(self) -> List[str]:
        """Get all unique conversations from the collection."""
        results = self.collection.get()
        conversations = set()
        for metadata in results["metadatas"]:
            conversations.add(metadata["conversation"])
        return sorted(list(conversations))


if __name__ == "__main__":
    # Example usage
    store = QuestionVectorStore()

    # Add questions from a file
    store.add_questions("./questions/RYaTvO_ZcMA.json")
