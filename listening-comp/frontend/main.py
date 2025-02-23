"""Spanish Listening Comprehension App

This module provides a Streamlit-based web interface for practicing Spanish
listening comprehension through interactive questions and conversations.
"""

# Configure Python path
import sys
from pathlib import Path
sys.path.append(str(Path(__file__).parent.parent))

# External imports
import streamlit as st

# Local application imports
from backend.question_generator import QuestionGenerator

st.set_page_config(page_title="Spanish Listening Comp", layout="wide")

# Initialize question generator
if "question_generator" not in st.session_state:
    st.session_state.question_generator = QuestionGenerator()


def render_question_generator():
    st.header("Question Generator")
    st.markdown("Generate practice questions from Spanish conversations.")

    # Get all conversations from ChromaDB
    conversations = (
        st.session_state.question_generator.vector_store.get_all_conversations()
    )

    if not conversations:
        st.warning(
            "No conversations found in the database. Please add some questions first."
        )
        return

    # Dropdown for conversation selection
    conversation = st.selectbox("Select a conversation:", conversations, index=0)

    num_questions = st.number_input(
        "Number of questions to generate:", min_value=1, max_value=20, value=1
    )

    # Initialize questions in session state if not exists
    if "questions" not in st.session_state:
        st.session_state.questions = None

    # Generate button
    if st.button("Generate Questions", key="generate_btn") and conversation:
        with st.spinner("Generating questions..."):
            try:
                st.session_state.questions = (
                    st.session_state.question_generator.generate_questions(
                        conversation=conversation, num_questions=num_questions
                    )
                )
            except Exception as e:
                st.error(f"Error generating questions: {str(e)}")
                return

    # Display questions if they exist
    if st.session_state.questions:
        st.subheader("Generated Questions:")

        # Track total and correct answers
        total_checked = 0
        total_correct = 0

        for i, q in enumerate(st.session_state.questions, 1):
            with st.expander(f"Question {i}", expanded=True):
                # Display conversation with audio controls
                st.markdown("**Conversation:**")

                # Display conversation
                if "Conversation" in q:
                    # Generate audio if needed
                    if "conversation_audio" not in q:
                        audio_result = st.session_state.question_generator.audio_generator.generate_conversation_audio(
                            q["Conversation"]
                        )
                        q["conversation_audio"] = audio_result["conversation_audio"]

                    # Display the full conversation audio at the top
                    if "conversation_audio" in q:
                        st.audio(q["conversation_audio"], format="audio/mp3")

                    # Split conversation into lines and parse speakers
                    conversation_lines = []
                    for line in q["Conversation"].strip().split("\n"):
                        if "]:" in line:
                            # Extract speaker info [Name] and text
                            parts = line.split("]:")
                            if len(parts) == 2:
                                name = parts[0].strip("[")  # Remove brackets
                                text = parts[1].strip()

                                conversation_lines.append(
                                    {
                                        "speaker": name,
                                        "text": text,
                                        "is_first": len(conversation_lines)
                                        == 0,  # Track if this is the first speaker
                                    }
                                )

                    # Create columns for the conversation text
                    if conversation_lines:
                        cols = st.columns([1, 3])

                        # Display headers
                        with cols[0]:
                            st.markdown("*Speaker*")
                        with cols[1]:
                            st.markdown("*Text*")

                        # Display each line with alternating styles
                        for line in conversation_lines:
                            with cols[0]:
                                # Add a subtle distinction between speakers using emojis
                                speaker_icon = (
                                    "💬" if line.get("is_first", False) else "💭"
                                )
                                st.markdown(f"{speaker_icon} *{line['speaker']}*")
                            with cols[1]:
                                st.markdown(line["text"])

                # Question with audio controls
                col3, col4 = st.columns([3, 1])
                with col3:
                    st.markdown(f"**Question:**\n{q['Question']}")
                with col4:
                    if "question_audio" in q:
                        st.audio(q["question_audio"], format="audio/mp3")

                # Display multiple choice options
                st.markdown("**Choose the correct answer:**")

                # Store question state in session state if not exists
                q_state_key = f"q_state_{i}"
                if q_state_key not in st.session_state:
                    import random

                    # Randomize choices and their corresponding audio
                    choices = q["Choices"].copy()
                    choice_audio = (
                        q.get("choice_audio", []).copy() if "choice_audio" in q else []
                    )

                    # Create list of tuples (choice, audio) and shuffle together
                    choice_pairs = (
                        list(zip(choices, choice_audio))
                        if choice_audio
                        else [(c, None) for c in choices]
                    )
                    random.shuffle(choice_pairs)

                    # Unzip the shuffled pairs
                    choices, choice_audio = (
                        zip(*choice_pairs) if choice_pairs else (choices, [])
                    )

                    # Store randomized choices, audio, and correct answer
                    st.session_state[q_state_key] = {
                        "choices": list(choices),  # Convert to list
                        "choice_audio": (
                            list(choice_audio) if choice_audio else []
                        ),  # Convert to list
                        "correct_answer": q[
                            "Answer"
                        ],  # Use the correct answer from question
                        "selected": None,
                        "checked": False,
                        "is_correct": False,
                    }

                # Get state from session
                state = st.session_state[q_state_key]

                # Radio group for all choices
                selected = st.radio(
                    "Select your answer:",
                    options=state["choices"],
                    key=f"q_{i}_choices",
                    label_visibility="collapsed",
                    horizontal=True,
                )

                # Display audio players next to choices
                st.write("")
                for j, choice in enumerate(state["choices"]):
                    cols = st.columns([4, 1])
                    with cols[0]:
                        if choice == selected:
                            st.markdown(f"**→ {choice}**")
                        else:
                            st.markdown(f"&nbsp;&nbsp;&nbsp;{choice}")
                    with cols[1]:
                        if "choice_audio" in state and j < len(state["choice_audio"]):
                            st.audio(state["choice_audio"][j], format="audio/mp3")

                # Update selected answer in state
                state["selected"] = selected

                # Check Answer button
                if st.button("Check Answer", key=f"check_{i}"):
                    state["checked"] = True
                    state["is_correct"] = selected == state["correct_answer"]

                # Show answer feedback if checked
                if state["checked"]:
                    if state["is_correct"]:
                        st.success("¡Correcto! (Correct!)")
                    else:
                        st.error(
                            f"Incorrect. The correct answer is: {state['correct_answer']}"
                        )

                    # Update totals
                    total_checked += 1
                    if state["is_correct"]:
                        total_correct += 1

        # Display final score if all questions are checked
        if total_checked == len(st.session_state.questions):
            score_percentage = (total_correct / total_checked) * 100
            st.markdown("---")
            st.subheader("Final Score")
            st.markdown(
                f"You got **{total_correct}** out of **{total_checked}** questions correct (**{score_percentage:.1f}%**)"
            )

            # Add encouraging message based on score
            if score_percentage == 100:
                st.success("¡Perfecto! Perfect score! 🌟")
                st.markdown(
                    "[Click here for your special reward! 🎵](https://www.youtube.com/watch?v=dQw4w9WgXcQ)"
                )
            elif score_percentage >= 80:
                st.success("¡Muy bien! Great job! 🎉")
            elif score_percentage >= 60:
                st.info("¡Bien! Keep practicing! 💪")
            else:
                st.info("Keep studying! You'll improve! 📚")


def main():
    st.title("Spanish Learning Assistant")

    # Create tabs
    tabs = st.tabs(["Question Generator"])

    with tabs[0]:
        render_question_generator()


if __name__ == "__main__":
    main()
