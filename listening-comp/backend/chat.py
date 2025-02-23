import google.generativeai as genai
from .config import GOOGLE_API_KEY

MODEL_NAME = "gemini-2.0-flash-lite-preview-02-05"


class GoogleGenAIChat:
    def __init__(self, model_name: str = MODEL_NAME, api_key: str = GOOGLE_API_KEY):
        genai.configure(api_key=api_key)
        self.model = genai.GenerativeModel(model_name)
        self.chat = self.model.start_chat()

    def generate_response(self, message: str):
        response = self.chat.send_message(message)
        return response.text


if __name__ == "__main__":
    chat = GoogleGenAIChat()
    while True:
        user_input = input("User: ")
        response = chat.generate_response(user_input)
        print(f"Bot: {response}")
