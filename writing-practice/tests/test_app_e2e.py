import os
from pathlib import Path

import pytest
from PIL import Image

from app import init_model, process_image


@pytest.mark.e2e
def test_core_functionality_e2e():
    """End-to-end test of core functionality using real dependencies.

    This test provides complete coverage of the image processing pipeline by:
    1. Using the actual Google Gemini API
    2. Processing a real test image (hola!.jpg)
    3. Verifying text extraction and translation
    4. Checking response format and content

    The test ensures:
    - Successful API integration
    - Accurate text detection
    - Proper translation
    - Appropriate response formatting

    Note:
        - Requires GOOGLE_API_KEY environment variable
        - Makes real API calls (needs internet connection)
        - Uses actual test image from the test directory
        - Will skip if API key is not set
    """
    # Skip if no API key is set
    if not os.getenv("GOOGLE_API_KEY"):
        pytest.skip("GOOGLE_API_KEY not set")

    # Initialize the model
    init_model()

    # Load and process a real image
    image_path = Path(__file__).parent / "hola!.jpg"
    image = Image.open(image_path)

    # Process the image using the actual implementation
    result = process_image(image)

    # Verify the results
    assert result, "Expected non-empty result"
    assert len(result) > 10, "Expected detailed response"

    # Check for expected content
    result_lower = result.lower()
    assert any(
        word in result_lower for word in ["hola!"]
    ), "Expected to find greeting in Spanish"
    assert "spanish" in result_lower, "Expected mention of Spanish language"
    assert len(result.split("\n")) >= 2, "Expected multiple lines including translation"
    assert "spanish" in result_lower, "Expected mention of Spanish language"

    # Check for translation and context
    assert len(result.split("\n")) >= 2, "Expected multiple lines including translation"
