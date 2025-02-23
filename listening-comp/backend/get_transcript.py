import os
from youtube_transcript_api import YouTubeTranscriptApi
from typing import Optional, List, Dict

LANGUAGES = ["en", "es"]


class YouTubeTranscriptDownloader:
    def extract_video_id(self, url: str) -> Optional[str]:
        return url.split("v=")[1][:11]

    def get_transcript(self, video_id: str) -> Optional[List[Dict]]:
        return YouTubeTranscriptApi.get_transcript(video_id, LANGUAGES)

    def save_transcript(self, transcript: List[Dict], filename: str) -> bool:
        os.makedirs(os.path.dirname(filename), exist_ok=True)
        with open(filename, "w") as f:
            for entry in transcript:
                f.write(f"{entry['text']}\n")
        return True


def main(video_url: str):
    downloader = YouTubeTranscriptDownloader()
    video_id = downloader.extract_video_id(video_url)
    transcript = downloader.get_transcript(video_id)
    downloader.save_transcript(transcript, f"./transcripts/{video_id}.txt")


if __name__ == "__main__":
    video_url = "https://www.youtube.com/watch?v=RYaTvO_ZcMA"
    transcript = main(video_url)
