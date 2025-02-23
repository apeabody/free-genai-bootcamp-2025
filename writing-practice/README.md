# Spanish Text Recognition App

This Streamlit application uses Google's Gemini 2.0 Flash Lite Preview model to extract and translate Spanish text from uploaded images. This model provides fast and efficient processing of visual content.

## Setup

1. Install the required dependencies:
```bash
pip install -r requirements.txt
```

2. Set up your Google API key in the `.env` file:
```bash
# Edit the .env file and replace 'your-api-key-here' with your actual Google API key
GOOGLE_API_KEY=your-api-key-here
```

3. Run the application:
```bash
streamlit run app.py
```

## Features

- Upload images containing Spanish text
- Extract text, translate, and receive context using Google's Gemini 2.0 Flash Lite Preview model

## Requirements

- Python 3.12+
- Google API key with access to Gemini 2.0 Flash Lite Preview model
- Internet connection for API access
