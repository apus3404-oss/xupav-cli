"""Ollama local AI provider implementation."""
from typing import Iterator
import ollama
from .base import BaseProvider, ProviderResponse


class OllamaProvider(BaseProvider):
    """Ollama local AI provider"""

    def __init__(self, base_url: str = "http://localhost:11434", **kwargs):
        super().__init__(api_key=None, **kwargs)
        self.client = ollama.Client(host=base_url)
        self.timeout = kwargs.get("timeout", 120)

    def chat(self, message: str, model: str, **kwargs) -> ProviderResponse:
        """Send chat message to Ollama"""
        try:
            response = self.client.chat(
                model=model,
                messages=[{"role": "user", "content": message}]
            )

            content = response["message"]["content"]

            # Calculate tokens (Ollama provides eval counts)
            prompt_tokens = response.get("prompt_eval_count", 0)
            completion_tokens = response.get("eval_count", 0)
            total_tokens = prompt_tokens + completion_tokens

            return ProviderResponse(
                content=content,
                tokens=total_tokens,
                cost=0.0,  # Ollama is free (local)
                metadata={
                    "model": model,
                    "prompt_tokens": prompt_tokens,
                    "completion_tokens": completion_tokens
                }
            )

        except Exception as e:
            raise RuntimeError(f"Ollama API error: {str(e)}")

    def stream_chat(self, message: str, model: str, **kwargs) -> Iterator[str]:
        """Stream chat response from Ollama"""
        try:
            stream = self.client.chat(
                model=model,
                messages=[{"role": "user", "content": message}],
                stream=True
            )

            for chunk in stream:
                if "message" in chunk and "content" in chunk["message"]:
                    yield chunk["message"]["content"]

        except Exception as e:
            raise RuntimeError(f"Ollama streaming error: {str(e)}")
