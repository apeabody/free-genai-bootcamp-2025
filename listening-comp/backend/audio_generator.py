from pathlib import Path
import hashlib
import json
import io
import tempfile
from typing import Dict, List
from gtts import gTTS
from pydub import AudioSegment


class AudioGenerator:
    def __init__(self, cache_dir: str = ".audio_cache"):
        # Create cache directory if it doesn't exist
        self.cache_dir = Path(cache_dir)
        self.cache_dir.mkdir(parents=True, exist_ok=True)

        # Define voice configurations
        self.voices = {
            "spain": {"lang": "es", "tld": "es", "slow": False},  # Spain Spanish
            "mexico": {"lang": "es", "tld": "com.mx", "slow": False},  # Mexican Spanish
            "narrator": {
                "lang": "es",
                "tld": "us",  # United States Spanish
                "slow": True,
            },
        }

    def _get_cache_path(self, text: str, speaker: str) -> Path:
        """Generate a unique cache path for the given text and speaker."""
        content_hash = hashlib.md5(f"{text}:{speaker}".encode()).hexdigest()
        return self.cache_dir / f"{content_hash}.mp3"

    def _generate_audio(self, text: str, speaker: str) -> bytes:
        """Generate audio for the given text using gTTS."""
        voice_config = self.voices[speaker]

        # Try different TLDs if the first one fails
        tlds = [voice_config["tld"], "com", "es"]
        last_error = None

        for tld in tlds:
            try:
                # Create gTTS object
                tts = gTTS(
                    text=text,
                    lang=voice_config["lang"],
                    tld=tld,
                    slow=voice_config["slow"],
                )

                # Save to a temporary file and read the bytes
                with tempfile.NamedTemporaryFile(suffix=".mp3") as temp_file:
                    tts.save(temp_file.name)
                    temp_file.seek(0)
                    audio_data = temp_file.read()
                return audio_data
            except Exception as e:
                last_error = e
                continue

        # If all TLDs failed, try one more time with minimal settings
        try:
            tts = gTTS(text=text, lang="es")
            with tempfile.NamedTemporaryFile(suffix=".mp3") as temp_file:
                tts.save(temp_file.name)
                temp_file.seek(0)
                audio_data = temp_file.read()
            return audio_data
        except Exception as e:
            msg = (
                f"Failed to generate audio: {str(last_error)} "
                f"and fallback also failed: {str(e)}"
            )
            raise Exception(msg)

    def _combine_audio_files(self, audio_files: List[str]) -> str:
        """Combine multiple audio files into a single file."""
        # Convert audio files to AudioSegment format

        # Convert each audio file to AudioSegment
        segments = []
        for audio_file in audio_files:
            with open(audio_file, "rb") as f:
                audio_bytes = io.BytesIO(f.read())
                segment = AudioSegment.from_mp3(audio_bytes)
                segments.append(segment)
                # Add a short pause between segments (500ms)
                segments.append(AudioSegment.silent(duration=500))

        # Combine all segments
        combined = sum(segments, AudioSegment.empty())

        # Export to a temporary file
        temp = tempfile.NamedTemporaryFile(suffix=".mp3", delete=False)
        combined.export(temp.name, format="mp3")
        return temp.name

    def generate_conversation_audio(self, conversation: str) -> Dict[str, str]:
        """Generate audio for a conversation with multiple speakers."""
        audio_files = []
        # Generate hash for conversation
        conv_bytes = conversation.encode()
        conversation_hash = hashlib.md5(conv_bytes).hexdigest()
        # Create path for combined audio file
        filename = f"combined_{conversation_hash}.mp3"
        combined_cache_path = self.cache_dir / filename

        # Check if we already have the combined audio
        if combined_cache_path.exists():
            return {"conversation_audio": str(combined_cache_path.absolute())}

        # Generate individual audio files for each line
        for line in conversation.strip().split("\n"):
            if "]:" in line:
                # Parse text
                text = line.split("]:")[1].strip()

                # Alternate between Spanish and Mexican voices
                voice = "spain" if len(audio_files) % 2 == 0 else "mexico"

                # Generate audio for this line
                cache_path = self._get_cache_path(text, voice)
                if not cache_path.exists():
                    audio_content = self._generate_audio(text, voice)
                    with open(cache_path, "wb") as f:
                        f.write(audio_content)

                audio_files.append(str(cache_path))

        # Combine all audio files
        if audio_files:
            combined_path = self._combine_audio_files(audio_files)
            # Move to cache location
            from shutil import move

            move(combined_path, str(combined_cache_path))

        return {"conversation_audio": str(combined_cache_path.absolute())}

    def generate_question_audio(
        self, question: str, choices: List[str]
    ) -> Dict[str, str | List[str]]:
        """Generate audio for a question and its choices."""
        # Generate audio for question
        question_cache = self._get_cache_path(question, "narrator")
        if not question_cache.exists():
            audio_content = self._generate_audio(question, "narrator")
            with open(question_cache, "wb") as f:
                f.write(audio_content)

        # Generate audio for choices
        choice_files = []
        for choice in choices:
            choice_cache = self._get_cache_path(choice, "narrator")
            if not choice_cache.exists():
                audio_content = self._generate_audio(choice, "narrator")
                with open(choice_cache, "wb") as f:
                    f.write(audio_content)
            choice_files.append(str(choice_cache.absolute()))

        return {
            "question_audio": str(question_cache.absolute()),
            "choice_audio": choice_files,
        }


if __name__ == "__main__":
    # Test the audio generator
    import asyncio

    async def test():
        generator = AudioGenerator()

        # Test conversation
        conversation = """
        María: ¡Hola! ¿Cómo estás?
        Juan: Muy bien, gracias. ¿Y tú?
        María: Excelente, acabo de terminar mi tarea.
        """

        conv_result = await generator.generate_conversation_audio(conversation)
        print("Conversation audio files:", json.dumps(conv_result, indent=2))

        # Test question
        question = "¿Qué acaba de hacer María?"
        choices = [
            "Terminar su tarea",
            "Ir al cine",
            "Hablar por teléfono",
            "Cocinar la cena",
        ]

        q_result = await generator.generate_question_audio(question, choices)
        print("\nQuestion audio files:", json.dumps(q_result, indent=2))

    asyncio.run(test())
