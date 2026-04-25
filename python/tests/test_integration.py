"""Integration tests for complete AI flow."""
import pytest
from mycli_ai.providers.openrouter import OpenRouterProvider
from mycli_ai.parser import extract_code_blocks, detect_file_changes
from mycli_ai.cost import CostTracker


def test_full_flow_with_mocks(mocker):
    """Test complete flow from request to parsed response"""

    # Mock OpenAI response with code block
    mock_message = mocker.Mock()
    mock_message.content = """
Here's the fix:

```python
# db.py
def fixed_query():
    return "SELECT * FROM users"
```
"""

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

    # Step 1: Create provider
    provider = OpenRouterProvider(api_key="test-key")

    # Step 2: Make request
    response = provider.chat("fix the query", model="deepseek/deepseek-r1")

    # Step 3: Verify response
    assert response.content is not None
    assert response.tokens == 150
    assert response.cost > 0

    # Step 4: Parse response for code blocks
    blocks = extract_code_blocks(response.content)
    assert len(blocks) == 1
    assert blocks[0].language == "python"
    assert "fixed_query" in blocks[0].code

    # Step 5: Detect file changes
    changes = detect_file_changes(response.content)
    assert len(changes) == 1
    assert changes[0].file_path == "db.py"
    assert "fixed_query" in changes[0].content

    # Step 6: Track cost
    tracker = CostTracker()
    tracker.add_usage(
        model="deepseek/deepseek-r1",
        prompt_tokens=50,
        completion_tokens=100
    )

    summary = tracker.get_summary()
    assert summary["total_tokens"] == 150
    assert summary["total_cost"] > 0


def test_ollama_integration(mocker):
    """Test integration with Ollama provider"""
    from mycli_ai.providers.ollama import OllamaProvider

    # Mock ollama response
    mock_response = {
        "message": {"content": "Here's a simple function:\n\n```go\n// main.go\nfunc Hello() string {\n    return \"Hello\"\n}\n```"},
        "eval_count": 80,
        "prompt_eval_count": 40
    }

    mock_client = mocker.Mock()
    mock_client.chat.return_value = mock_response
    mocker.patch('ollama.Client', return_value=mock_client)

    # Create provider and make request
    provider = OllamaProvider(base_url="http://localhost:11434")
    response = provider.chat("write a hello function", model="codellama:13b")

    # Verify response
    assert response.content is not None
    assert response.tokens == 120
    assert response.cost == 0.0  # Ollama is free

    # Parse code blocks
    blocks = extract_code_blocks(response.content)
    assert len(blocks) == 1
    assert blocks[0].language == "go"
    assert blocks[0].filename == "main.go"

    # Detect file changes
    changes = detect_file_changes(response.content)
    assert len(changes) == 1
    assert changes[0].file_path == "main.go"


def test_multiple_providers_cost_tracking(mocker):
    """Test cost tracking across multiple providers"""
    # Mock OpenRouter
    mock_or_response = mocker.Mock()
    mock_or_response.choices = [mocker.Mock(message=mocker.Mock(content="OpenRouter response"))]
    mock_or_response.usage = mocker.Mock(total_tokens=100, prompt_tokens=40, completion_tokens=60)

    mock_or_client = mocker.Mock()
    mock_or_client.chat.completions.create.return_value = mock_or_response
    mocker.patch('openai.OpenAI', return_value=mock_or_client)

    # Mock Ollama
    mock_ollama_response = {
        "message": {"content": "Ollama response"},
        "eval_count": 50,
        "prompt_eval_count": 30
    }

    mock_ollama_client = mocker.Mock()
    mock_ollama_client.chat.return_value = mock_ollama_response
    mocker.patch('ollama.Client', return_value=mock_ollama_client)

    # Create tracker
    tracker = CostTracker()

    # Use OpenRouter
    from mycli_ai.providers.openrouter import OpenRouterProvider
    or_provider = OpenRouterProvider(api_key="test-key")
    or_response = or_provider.chat("test", model="deepseek/deepseek-r1")
    tracker.add_usage("deepseek/deepseek-r1", 40, 60)

    # Use Ollama
    from mycli_ai.providers.ollama import OllamaProvider
    ollama_provider = OllamaProvider()
    ollama_response = ollama_provider.chat("test", model="codellama:13b")
    tracker.add_usage("codellama:13b", 30, 50)

    # Verify tracking
    summary = tracker.get_summary()
    assert summary["total_tokens"] == 180
    assert len(summary["by_model"]) == 2
    assert "deepseek/deepseek-r1" in summary["by_model"]
    assert "codellama:13b" in summary["by_model"]
