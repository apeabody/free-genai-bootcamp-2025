# Spanish Listening Comprehension App

## Getting Started

1. Install dependencies:
   ```bash
   pip install -r requirements.txt
   ```

2. Create a `.env` file in the `backend` directory with your API key:
   ```
   GOOGLE_API_KEY=your_api_key_here
   ```

3. Setup the backend:
   ```bash
   python -m backend.get_transcript
   python -m backend.structured_data
   python -m backend.vector_store
   ```

4. Run the app:
   ```bash
   streamlit run frontend/main.py
   ```

## Note
Make sure to keep your API key secure and never commit it to version control.
