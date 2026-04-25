"""OpenRouter API provider implementation."""
from typing import Optional, Iterator
import openai
from .base import BaseProvider, ProviderResponse

# Model pricing (per 1M tokens)
MODEL_PRICING = {
    "deepseek/deepseek-r1": {"input": 0.14, "output": 0.14},
    "stepfun/step-2-16k": {"input": 0.10, "output": 0.10},
    "anthropic/claude-sonnet-4": {"input": 3.00, "output": 15.00},
}


class OpenRouterProvider(BaseProvider):
    """OpenRouter API provider"""

    def __init__(self, api_key: str, base_url: str = "https://openrouter.ai/api/v1", **kwargs):
        super().__init__(api_key, **kwargs)
        self.client = openai.OpenAI(
            api_key=api_key,
            base_url=base_url
        )
        self.timeout = kwargs.get("timeout", 60)

    def chat(self, message: str, model: str, **kwargs) -> ProviderResponse:
        """Send chat message to OpenRouter"""
        max_tokens = kwargs.get("max_tokens", 4096)
        temperature = kwargs.get("temperature", 0.7)

        try:
            response = self.client.chat.completions.create(
                model=model,
                messages=[{"role": "user", "content": message}],
                max_tokens=max_tokens,
                temperature=temperature,
                timeout=self.timeout
            )

            content = response.choices[0].message.content
            tokens = response.usage.total_tokens
            cost = self._calculate_cost(
                model,
                response.usage.prompt_tokens,
                response.usage.completion_tokens
            )

            return ProviderResponse(
                content=content,
                tokens=tokens,
                cost=cost,
                metadata={"model": model}
            )

        except Exception as e:
            raise RuntimeError(f"OpenRouter API error: {str(e)}")

    def stream_chat(self, message: str, model: str, **kwargs) -> Iterator[str]:
        """Stream chat response from OpenRouter"""
        max_tokens = kwargs.get("max_tokens", 4096)
        temperature = kwargs.get("temperature", 0.7)

        try:
            stream = self.client.chat.completions.create(
                model=model,
                messages=[{"role": "user", "content": message}],
                max_tokens=max_tokens,
                temperature=temperature,
                timeout=self.timeout,
                stream=True
            )

            for chunk in stream:
                if chunk.choices[0].delta.content:
                    yield chunk.choices[0].delta.content

        except Exception as e:
            raise RuntimeError(f"OpenRouter streaming error: {str(e)}")

    def _calculate_cost(self, model: str, prompt_tokens: int, completion_tokens: int) -> float:
        """Calculate cost based on token usage"""
        if model not in MODEL_PRICING:
            # Default pricing if model not found
            return (prompt_tokens + completion_tokens) * 0.0001 / 1000

        pricing = MODEL_PRICING[model]
        input_cost = (prompt_tokens / 1_000_000) * pricing["input"]
        output_cost = (completion_tokens / 1_000_000) * pricing["output"]

        return input_cost + output_cost
