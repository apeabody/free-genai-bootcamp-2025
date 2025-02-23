from typing import List, Dict
import json

from backend.chat import GoogleGenAIChat
from backend.vector_store import QuestionVectorStore
from backend.audio_generator import AudioGenerator


class QuestionGenerator:
    def __init__(self):
        self.chat = GoogleGenAIChat()
        self.vector_store = QuestionVectorStore()
        self.audio_generator = AudioGenerator()

    def generate_questions(
        self, conversation: str, num_questions: int = 1
    ) -> List[Dict]:
        """Generate new questions based on a conversation.

        Args:
            conversation: The Spanish conversation to use as an example style
            num_questions: Number of questions to generate (default: 1)

        Returns:
            List of dictionaries containing conversations, questions and answers
        """
        # Create the prompt
        prompt = """You are a Spanish language teacher creating multiple choice questions for listening comprehension practice.
        For each question, you will create:
        1. A new conversation between two people
        2. A specific question about that conversation
        3. One correct answer
        4. Three incorrect but plausible answers in Spanish

        IMPORTANT: For the conversation, use this exact format:
        [Speaker 1]: ¡Hola! ¿Cómo estás?
        [Speaker 2]: Muy bien, gracias. ¿Y tú?

        Make sure each line of dialogue is on its own line.
        Do not add any extra text or descriptions.

        Here are some example questions from our database:
        """

        # Add the target conversation as an example
        prompt += f"""
        Here's an example conversation in the style we want:
        {conversation}

        Now generate {num_questions} multiple choice questions. For each question:
        1. Create a NEW conversation between two people in a similar style and difficulty level
        2. Make the conversation natural and relevant to everyday Spanish usage
        3. Create a specific question about the conversation
        4. Create one correct answer that directly answers the question
        5. Create three incorrect but plausible answers in Spanish that could trick a student

        Use this exact format for each question:
        Conversation: [new conversation between two people]
        Question: [specific question about the conversation]
        Correct Answer: [correct answer]
        Incorrect Answer 1: [plausible wrong answer]
        Incorrect Answer 2: [plausible wrong answer]
        Incorrect Answer 3: [plausible wrong answer]

        Make sure each conversation is different and covers various everyday topics and situations.
        """

        # Get response from AI
        response = self.chat.generate_response(prompt)

        # Parse the response into a list of questions
        questions = []
        current_question = None
        conversation_lines = []

        for line in response.split("\n"):
            line = line.strip()
            if not line:
                continue

            if line.startswith("Conversation:"):
                # Start a new question
                if current_question and "Question" in current_question:
                    if conversation_lines:
                        current_question["Conversation"] = "\n".join(conversation_lines)
                    questions.append(current_question)
                current_question = {}
                conversation_lines = []
            elif line.startswith("Question:"):
                # End of conversation section
                if current_question and conversation_lines:
                    current_question["Conversation"] = "\n".join(conversation_lines)
                current_question["Question"] = line[9:].strip()
            elif (
                ":" in line
                and current_question is not None
                and not any(
                    line.startswith(prefix)
                    for prefix in ["Question:", "Correct Answer:", "Incorrect Answer"]
                )
            ):
                # This is a conversation line (e.g. "María: Hola")
                conversation_lines.append(line)
                current_question["Question"] = line[9:].strip()
            elif line.startswith("Correct Answer:"):
                current_question["Answer"] = line[15:].strip()
                current_question["Choices"] = [
                    line[15:].strip()
                ]  # Add correct answer as first choice
            elif line.startswith("Incorrect Answer 1:"):
                current_question["Choices"].append(line[19:].strip())
            elif line.startswith("Incorrect Answer 2:"):
                current_question["Choices"].append(line[19:].strip())
            elif line.startswith("Incorrect Answer 3:"):
                current_question["Choices"].append(line[19:].strip())

        # Add the last question
        if current_question:
            questions.append(current_question)

        # Generate audio for each question
        for q in questions:
            # Parse conversation into turns
            conversation_lines = []
            current_conversation = []

            # Split conversation into lines
            for line in q["Conversation"].strip().split("\n"):
                line = line.strip()
                if line and "]:" in line:
                    # Extract speaker info and text
                    parts = line.split("]:")
                    if len(parts) == 2:
                        speaker_info = parts[0] + "]"
                        text = parts[1].strip()
                        current_conversation.append(f"{speaker_info}: {text}")
                    else:
                        current_conversation.append(line)
                elif line:
                    current_conversation.append(line)

            # Join the conversation lines back together
            q["Conversation"] = "\n".join(current_conversation)

            # Generate audio for conversation
            audio_result = self.audio_generator.generate_conversation_audio(
                q["Conversation"]
            )
            q["conversation_audio"] = audio_result["conversation_audio"]

            # Generate audio for question and choices
            qa_result = self.audio_generator.generate_question_audio(
                q["Question"], q["Choices"]
            )
            q["question_audio"] = qa_result["question_audio"]
            q["choice_audio"] = qa_result["choice_audio"]

        return questions

    def save_generated_questions(self, questions: List[Dict], filename: str) -> None:
        """Save generated questions to a JSON file."""
        with open(filename, "w", encoding="utf-8") as f:
            json.dump(questions, f, ensure_ascii=False, indent=2)


if __name__ == "__main__":
    # Example usage of QuestionGenerator
    generator = QuestionGenerator()
    print("QuestionGenerator initialized with ChromaDB vector store.")
    print("Use generate_questions(conversation, num_questions) to generate questions.")
