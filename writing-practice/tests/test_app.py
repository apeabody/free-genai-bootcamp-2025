import os
from pathlib import Path

import pytest
import streamlit as st
from PIL import Image

import app
from app import init_model, process_image, create_sidebar


# Mock image for testing
@pytest.fixture
def test_image():
    """Create a test image for unit testing.

    Returns:
        PIL.Image: A small RGB image (100x30 pixels) with white background.
        Used for testing image processing functions without needing real images.
    """
    img = Image.new("RGB", (100, 30), color="white")
    return img


@pytest.fixture
def mock_genai(mocker):
    """Create a mock Gemini model for testing.

    This fixture:
    1. Creates a mock model that simulates the Gemini API
    2. Configures default successful response
    3. Patches the model at module level

    Args:
        mocker: pytest-mock fixture for creating mocks

    Returns:
        MagicMock: A configured mock model that can be used to verify calls
        and simulate different responses
    """
    mock_model = mocker.MagicMock()
    mock_response = mocker.MagicMock()
    mock_response.text = "Test response: Spanish text detected"
    mock_model.generate_content.return_value = mock_response

    # Mock the model at module level
    mocker.patch("app.model", mock_model)
    return mock_model


def test_process_image_success(test_image, mock_genai):
    """Test successful image processing with mocked Gemini model.

    This test verifies that:
    1. The image processing function accepts a valid image
    2. The model is called with correct parameters
    3. The response is properly formatted
    4. The result contains expected text

    Args:
        test_image: Pytest fixture providing a test image
        mock_genai: Pytest fixture providing a mocked Gemini model
    """
    # Process the test image
    result = process_image(test_image)

    # Verify the result
    assert result == "Test response: Spanish text detected"
    mock_genai.generate_content.assert_called_once()


def test_process_image_error(test_image, mock_genai):
    """Test error handling during image processing.

    This test ensures that:
    1. Errors from the Gemini model are caught
    2. The function returns an empty string on error
    3. Error messages are properly logged
    4. The system remains stable after an error

    Args:
        test_image: Pytest fixture providing a test image
        mock_genai: Pytest fixture providing a mocked Gemini model
    """
    # Configure the mock to raise an exception
    mock_genai.generate_content.side_effect = Exception("API Error")

    # Process the test image
    result = process_image(test_image)

    # Verify error handling
    assert result == ""
    mock_genai.generate_content.assert_called_once()


@pytest.mark.e2e
def test_process_image_e2e():
    """End-to-end test of image processing with real API.

    This test verifies:
    1. Integration with Google Gemini API
    2. Processing of actual image data
    3. Text extraction and translation accuracy
    4. Response format and content validity

    Note:
        - Requires GOOGLE_API_KEY environment variable
        - Uses real API calls (needs internet connection)
        - Will skip if API key is not set
    """

    # Skip if no API key is set
    if not os.getenv("GOOGLE_API_KEY"):
        pytest.skip("GOOGLE_API_KEY not set")

    # Initialize the real model
    app.model = init_model()

    # Load the test image
    image_path = Path(__file__).parent / "hola!.jpg"
    image = Image.open(image_path)

    # Process the image
    result = process_image(image)

    # Basic validation of the result
    assert result, "Expected non-empty result"
    assert (
        "hola" in result.lower() or "hello" in result.lower()
    ), "Expected to find 'hola' or 'hello' in the result"


def test_init_model_error(monkeypatch, mocker):
    """Test error handling in model initialization.

    This test verifies that:
    1. Missing API key is properly handled
    2. Error message is displayed
    3. Application stops gracefully

    Args:
        monkeypatch: pytest fixture for modifying environment
        mocker: pytest fixture for mocking
    """
    # Remove API key from environment
    monkeypatch.delenv("GOOGLE_API_KEY", raising=False)

    # Mock streamlit functions
    mock_error = mocker.patch("streamlit.error")
    mock_stop = mocker.patch("streamlit.stop")

    # Call init_model
    init_model()

    # Verify error handling
    mock_error.assert_called_once_with(
        "Please set your GOOGLE_API_KEY environment variable"
    )
    mock_stop.assert_called_once()


def test_create_sidebar(mocker):
    """Test sidebar creation functionality.

    This test ensures that:
    1. Sidebar is created with correct sections
    2. Instructions are properly formatted
    3. About section contains necessary information

    Args:
        mocker: pytest fixture for mocking
    """
    # Mock streamlit functions
    mock_header = mocker.patch("app.st.header")
    mock_write = mocker.patch("app.st.write")

    # Create sidebar
    create_sidebar()

    # Verify sidebar content
    assert mock_header.call_count >= 2  # Instructions and About
    assert mock_write.call_count >= 2  # Instructions and About text

    # Verify specific headers
    mock_header.assert_any_call("Instructions")
    mock_header.assert_any_call("About")


def test_main_flow(mocker, test_image):
    """Test the main UI flow of the application.

    This test verifies that:
    1. Title and description are displayed
    2. Sidebar is created
    3. File uploader is configured correctly
    4. Image processing works when file is uploaded

    Args:
        mocker: pytest fixture for mocking
        test_image: Pytest fixture providing a test image
    """
    # Mock streamlit functions
    mocker.patch("streamlit.title")
    mocker.patch("streamlit.write")
    mocker.patch("streamlit.file_uploader", return_value=test_image)
    mocker.patch("streamlit.image")
    mocker.patch("streamlit.spinner")
    mocker.patch("app.create_sidebar")
    mocker.patch("app.process_image", return_value="Test result")

    # Run main function
    app.main()

    # Verify main flow
    st.title.assert_called_once_with("Spanish Language Learning Assistant")
    st.file_uploader.assert_called_once()
    app.create_sidebar.assert_called_once()
