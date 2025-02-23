# Test Documentation

This directory contains the test suite for the Spanish Text Recognition application. The tests are organized to ensure comprehensive coverage of both core functionality and end-to-end workflows.

## Test Structure

### Unit Tests (`test_app.py`)
Tests individual components with mocked dependencies.

- `test_process_image_success`: Verifies successful image processing
  - Tests the core image processing function with a mocked Gemini model
  - Ensures proper handling of valid image input
  - Verifies expected response format

- `test_process_image_error`: Tests error handling
  - Verifies proper error handling when model fails
  - Ensures user-friendly error messages
  - Tests system resilience

### End-to-End Tests (`test_app_e2e.py`)
Tests complete workflows with real dependencies.

- `test_process_image_e2e`: Tests core functionality with real API
  - Uses actual Google Gemini API
  - Processes a real test image
  - Verifies complete text extraction and translation workflow

- `test_core_functionality_e2e`: Comprehensive E2E test
  - Tests the entire image processing pipeline
  - Verifies text extraction, translation, and formatting
  - Ensures multi-line response with proper sections

## Running Tests

1. **Environment Setup**
   ```bash
   # Create and activate virtual environment
   python -m venv venv
   source venv/bin/activate

   # Install dependencies
   pip install -r requirements.txt
   ```

2. **Configure API Key**
   ```bash
   # Set your Google API key
   export GOOGLE_API_KEY='your-api-key'
   ```

3. **Run Tests**
   ```bash
   # Run all tests
   pytest tests/ -v

   # Run specific test file
   pytest tests/test_app.py -v

   # Run tests with specific marker
   pytest -m e2e
   ```

## Test Data
- `hola!.jpg`: Test image containing Spanish text
  - Used in E2E tests
  - Contains simple Spanish greeting
  - Ideal for testing basic functionality

## Notes
- Tests will skip if `GOOGLE_API_KEY` is not set
- E2E tests require internet connection for API access
- Test image should remain in the repository for consistent testing
