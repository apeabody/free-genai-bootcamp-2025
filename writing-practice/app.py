import os

import google.generativeai as genai
import streamlit as st
from dotenv import load_dotenv
from PIL import Image

# Load environment variables from .env file
load_dotenv()


def init_model():
    """Initialize and return the Gemini model.

    Returns:
        GenerativeModel: Initialized Gemini model

    Raises:
        SystemExit: If GOOGLE_API_KEY is not set
    """
    GOOGLE_API_KEY = os.getenv("GOOGLE_API_KEY")
    if not GOOGLE_API_KEY:
        st.error("Please set your GOOGLE_API_KEY environment variable")
        st.stop()

    genai.configure(api_key=GOOGLE_API_KEY)
    return genai.GenerativeModel("gemini-2.0-flash-lite-preview-02-05")


# Configure the page
st.set_page_config(page_title="Spanish Text Recognition", layout="wide")

# Initialize Google Gemini API
model = init_model()


def process_image(image) -> str:
    """Process an image to extract and translate Spanish text.

    Args:
        image: PIL Image object containing Spanish text

    Returns:
        str: Extracted text, translation, and context
    """
    prompt = """
    Please analyze this image and extract any Spanish text you see.
    If there is text in the image, provide:
    1. The original Spanish text
    2. An English translation
    3. Any relevant context about the text
    If no text is found, please indicate that.
    """

    try:
        response = model.generate_content([prompt, image])
        return response.text
    except Exception as e:
        st.error(f"Error processing image: {str(e)}")
        return ""


def create_sidebar() -> None:
    """Create the sidebar with instructions and about information."""
    with st.sidebar:
        st.header("Instructions")
        st.write(
            """
            1. Upload an image containing Spanish text
            2. Click 'Process Image'
            3. Wait for the results

            The app will:
            - Extract Spanish text from the image
            - Provide an English translation
            - Give context about the text
            """
        )

        st.header("About")
        st.write(
            """
            This app uses Google's Gemini 2.0 Flash Lite Preview model to analyze images
            containing Spanish text. It provides fast and efficient processing of
            Spanish content for writing practice.
            """
        )


def main() -> None:
    """Main function to run the Streamlit application."""
    # Initialize session state
    if "process_clicked" not in st.session_state:
        st.session_state.process_clicked = False

    # UI Components
    st.title("Spanish Language Learning Assistant")

    # Create tabs
    tab1, tab2 = st.tabs(["Writing Practice", "Coming Soon"])

    with tab1:
        st.header("Spanish Text Recognition")
        st.write("Upload an image containing Spanish text to extract and translate it.")

        # Create sidebar
        create_sidebar()

        # File uploader
        uploaded_file = st.file_uploader(
            "Choose an image file", type=["png", "jpg", "jpeg"]
        )

        if uploaded_file is not None:
            try:
                # Display the uploaded image
                image = Image.open(uploaded_file)
                st.image(image, caption="Uploaded Image", use_column_width=True)

                # Process button
                if st.button("Process Image") or st.session_state.process_clicked:
                    st.session_state.process_clicked = True
                    with st.spinner("Processing image..."):
                        result = process_image(image)
                        if result:
                            st.write("### Results:")
                            st.write(result)

            except Exception as e:
                st.error(f"Error loading image: {str(e)}")

    with tab2:
        st.info(
            "More features coming soon! Stay tuned for:"
            "\n- Speaking Practice"
            "\n- Vocabulary Builder"
            "\n- Grammar Exercises"
        )


if __name__ == "__main__":
    main()
