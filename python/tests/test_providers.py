"""Tests for AI provider implementations."""
import pytest
from mycli_ai.providers.base import BaseProvider, ProviderResponse


def test_provider_response_creation():
    """Test ProviderResponse dataclass creation"""
    response = ProviderResponse(
        content="Hello, world!",
        tokens=10,
        cost=0.001,
        metadata={"model": "test-model"}
    )

    assert response.content == "Hello, world!"
    assert response.tokens == 10
    assert response.cost == 0.001
    assert response.metadata["model"] == "test-model"


def test_base_provider_is_abstract():
    """Test that BaseProvider cannot be instantiated directly"""
    with pytest.raises(TypeError):
        BaseProvider(api_key="test-key")


class MockProvider(BaseProvider):
    """Mock provider for testing"""

    def chat(self, message: str, model: str, **kwargs) -> ProviderResponse:
        return ProviderResponse(
            content=f"Mock response to: {message}",
            tokens=len(message.split()),
            cost=0.001
        )

    def stream_chat(self, message: str, model: str, **kwargs):
        words = message.split()
        for word in words:
            yield word + " "


def test_base_provider_initialization():
    """Test BaseProvider initialization with config"""
    provider = MockProvider(api_key="test-key", timeout=30)

    assert provider.api_key == "test-key"
    assert provider.config["timeout"] == 30


def test_mock_provider_chat():
    """Test mock provider chat method"""
    provider = MockProvider(api_key="test-key")
    response = provider.chat("Hello AI", model="test-model")

    assert "Mock response to: Hello AI" in response.content
    assert response.tokens == 2
    assert response.cost == 0.001


def test_mock_provider_stream():
    """Test mock provider streaming"""
    provider = MockProvider(api_key="test-key")
    chunks = list(provider.stream_chat("Hello world", model="test-model"))

    assert len(chunks) == 2
    assert chunks[0] == "Hello "
    assert chunks[1] == "world "
