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


def test_openrouter_chat(mocker):
    """Test OpenRouter provider chat method"""
    from mycli_ai.providers.openrouter import OpenRouterProvider

    # Mock OpenAI client response
    mock_message = mocker.Mock()
    mock_message.content = "test response"

    mock_choice = mocker.Mock()
    mock_choice.message = mock_message

    mock_usage = mocker.Mock()
    mock_usage.total_tokens = 150
    mock_usage.prompt_tokens = 50
    mock_usage.completion_tokens = 100

    mock_response = mocker.Mock()
    mock_response.choices = [mock_choice]
    mock_response.usage = mock_usage

    mock_client = mocker.Mock()
    mock_client.chat.completions.create.return_value = mock_response
    mocker.patch('openai.OpenAI', return_value=mock_client)

    # Test
    provider = OpenRouterProvider(api_key="test-key")
    response = provider.chat("test message", model="deepseek/deepseek-r1")

    assert response.content == "test response"
    assert response.tokens == 150
    assert response.cost > 0


def test_openrouter_stream(mocker):
    """Test OpenRouter provider streaming"""
    from mycli_ai.providers.openrouter import OpenRouterProvider

    # Mock streaming response chunks
    mock_delta1 = mocker.Mock()
    mock_delta1.content = "Hello"
    mock_choice1 = mocker.Mock()
    mock_choice1.delta = mock_delta1
    mock_chunk1 = mocker.Mock()
    mock_chunk1.choices = [mock_choice1]

    mock_delta2 = mocker.Mock()
    mock_delta2.content = " world"
    mock_choice2 = mocker.Mock()
    mock_choice2.delta = mock_delta2
    mock_chunk2 = mocker.Mock()
    mock_chunk2.choices = [mock_choice2]

    mock_delta3 = mocker.Mock()
    mock_delta3.content = "!"
    mock_choice3 = mocker.Mock()
    mock_choice3.delta = mock_delta3
    mock_chunk3 = mocker.Mock()
    mock_chunk3.choices = [mock_choice3]

    mock_client = mocker.Mock()
    mock_client.chat.completions.create.return_value = iter([mock_chunk1, mock_chunk2, mock_chunk3])
    mocker.patch('openai.OpenAI', return_value=mock_client)

    # Test
    provider = OpenRouterProvider(api_key="test-key")
    chunks = list(provider.stream_chat("test", model="deepseek/deepseek-r1"))

    assert chunks == ["Hello", " world", "!"]


def test_ollama_chat(mocker):
    """Test Ollama provider chat method"""
    from mycli_ai.providers.ollama import OllamaProvider

    # Mock ollama client response
    mock_response = {
        "message": {"content": "test response"},
        "eval_count": 100,
        "prompt_eval_count": 50
    }

    mock_client = mocker.Mock()
    mock_client.chat.return_value = mock_response
    mocker.patch('ollama.Client', return_value=mock_client)

    # Test
    provider = OllamaProvider(base_url="http://localhost:11434")
    response = provider.chat("test message", model="codellama:13b")

    assert response.content == "test response"
    assert response.tokens == 150
    assert response.cost == 0.0  # Ollama is free


def test_ollama_stream(mocker):
    """Test Ollama provider streaming"""
    from mycli_ai.providers.ollama import OllamaProvider

    # Mock streaming response chunks
    mock_chunks = [
        {"message": {"content": "Hello"}},
        {"message": {"content": " world"}},
        {"message": {"content": "!"}},
    ]

    mock_client = mocker.Mock()
    mock_client.chat.return_value = iter(mock_chunks)
    mocker.patch('ollama.Client', return_value=mock_client)

    # Test
    provider = OllamaProvider(base_url="http://localhost:11434")
    chunks = list(provider.stream_chat("test", model="codellama:13b"))

    assert chunks == ["Hello", " world", "!"]
