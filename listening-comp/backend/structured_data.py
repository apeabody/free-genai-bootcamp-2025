import os
from .chat import GoogleGenAIChat


class StructuredData:
    def __init__(self):
        self.chat = GoogleGenAIChat()

    def read_transcript(self, transcript_file: str) -> str:
        with open(transcript_file, "r") as f:
            return f.read()

    def extract_questions(self, transcript: str):
        prompt = """You are an expert Spanish listening comprehension test
            question extractor. Given the following transcript of a Spanish
            listening comprehension test, extract the conversation, questions,
            and answer. Format the output in raw JSON as follows:

            Conversation: "conversation"
            Question: "question"
            Answer: "answer"

            Do not include the inital introduction.
            Do not number the conversations.
            Add punctuation to the questions.
            """

        prompt += f"\n\nTranscript: {transcript}"
        response = self.chat.generate_response(prompt)

        # Strip ```json from the response
        questions = response.replace("```json", "").replace("```", "").strip()
        return questions

    def save_questions(self, questions: str, filename: str) -> None:
        os.makedirs(os.path.dirname(filename), exist_ok=True)
        with open(filename, "w") as f:
            f.write(questions)


if __name__ == "__main__":
    video_id = "RYaTvO_ZcMA"
    structured_data = StructuredData()
    transcript_file = f"./transcripts/{video_id}.txt"
    transcript = structured_data.read_transcript(transcript_file)
    questions = structured_data.extract_questions(transcript)
    output_file = f"./questions/{video_id}.json"
    structured_data.save_questions(questions, output_file)
