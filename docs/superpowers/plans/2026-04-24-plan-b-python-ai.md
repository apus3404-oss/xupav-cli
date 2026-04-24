# Python AI Layer Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the Python AI layer with OpenRouter and Ollama clients, JSON-RPC server for Go communication, response parsing, and cost tracking.

**Architecture:** Python package that runs as a subprocess, communicates via JSON-RPC over stdin/stdout, handles API requests to OpenRouter/Ollama, parses responses, and tracks token usage/costs.

**Tech Stack:**
- **openai** - OpenRouter API client (compatible endpoint)
- **ollama-python** - Ollama client
- **tiktoken** - Token counting
- **pytest** - Testing framework

---

## File Structure

**New Files:**
```
python/
  mycli_ai/
    __init__.py
    server.py              # JSON-RPC server (stdin/stdout)
    providers/
      __init__.py
      base.py              # Provider interface
      openrouter.py        # OpenRouter implementation
      ollama.py            # Ollama implementation
    parser.py              # Response parser
    cost.py                # Cost calculation
    tokens.py              # Token counting
  requirements.txt         # Dependencies
  setup.py                 # Package metadata
  tests/
    test_providers.py
    test_parser.py
    test_server.py
    test_cost.py
```

---
## Task 1: Python Project Setup

**Files:**
- Create: `python/requirements.txt`
- Create: `python/setup.py`
- Create: `python/mycli_ai/__init__.py`

- [ ] **Step 1: Create requirements.txt**

```txt
openai==1.12.0
ollama==0.1.7
tiktoken==0.6.0
pytest==8.0.0
pytest-mock==3.12.0
```

- [ ] **Step 2: Create setup.py**

```python
# python/setup.py
from setuptools import setup, find_packages

setup(
    name="mycli-ai",
    version="0.1.0",
    packages=find_packages(),
    install_requires=[
        "openai>=1.12.0",
        "ollama>=0.1.7",
        "tiktoken>=0.6.0",
    ],
    python_requires=">=3.8",
    author="Your Name",
    description="AI layer for mycli",
)
```

- [ ] **Step 3: Create package __init__.py**

```python
# python/mycli_ai/__init__.py
"""
mycli AI layer - handles AI provider communication
"""

__version__ = "0.1.0"
```

- [ ] **Step 4: Create virtual environment and install**

```bash
cd python
python3 -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
pip install -e .
```

Expected: All dependencies installed successfully

- [ ] **Step 5: Verify installation**

```bash
python -c "import mycli_ai; print(mycli_ai.__version__)"
```

Expected output: `0.1.0`

- [ ] **Step 6: Commit**

```bash
git add python/
git commit -m "chore(python): initialize Python package structure"
```

---

## Task 2: Provider Base Interface

**Files:**
- Create: `python/mycli_ai/providers/__init__.py`
- Create: `python/mycli_ai/providers/base.py`
- Create: `python/tests/test_providers.py`

- [ ] **Step 1: Write failing test for provider interface**

```python
# python/tests/test_providers.py
import pytest
from mycli_ai.providers.base import BaseProvider, ProviderResponse

def test_provider_response():
    response = ProviderResponse(
        content="test response",
        tokens=100,
        cost=0.001
    )
    
    assert response.content == "test response"
    assert response.tokens == 100
    assert response.cost == 0.001

def test_base_provider_not_implemented():
    provider = BaseProvider()
    
    with pytest.raises(NotImplementedError):
        provider.chat("test message", model="test-model")
```

- [ ] **Step 2: Run test to verify it fails**

```bash
cd python
pytest tests/test_providers.py -v
```

Expected: FAIL with "ModuleNotFoundError: No module named 'mycli_ai.providers.base'"

- [ ] **Step 3: Create providers __init__.py**

```python
# python/mycli_ai/providers/__init__.py
from .base import BaseProvider, ProviderResponse

__all__ = ["BaseProvider", "ProviderResponse"]
```

- [ ] **Step 4: Implement base provider**

```python
# python/mycli_ai/providers/base.py
from abc import ABC, abstractmethod
from dataclasses import dataclass
from typing import Optional, Dict, Any

@dataclass
class ProviderResponse:
    """Response from an AI provider"""
    content: str
    tokens: int
    cost: float
    metadata: Optional[Dict[str, Any]] = None

class BaseProvider(ABC):
    """Base class for AI providers"""
    
    def __init__(self, api_key: Optional[str] = None, **kwargs):
        self.api_key = api_key
        self.config = kwargs
    
    @abstractmethod
    def chat(self, message: str, model: str, **kwargs) -> ProviderResponse:
        """
        Send a chat message and get response
        
        Args:
            message: User message
            model: Model identifier
            **kwargs: Additional provider-specific parameters
            
        Returns:
            ProviderResponse with content, tokens, and cost
        """
        raise NotImplementedError("Subclasses must implement chat()")
    
    @abstractmethod
    def stream_chat(self, message: str, model: str, **kwargs):
        """
        Stream chat response chunk by chunk
        
        Args:
            message: User message
            model: Model identifier
            **kwargs: Additional provider-specific parameters
            
        Yields:
            String chunks of the response
        """
        raise NotImplementedError("Subclasses must implement stream_chat()")
```

- [ ] **Step 5: Run test to verify it passes**

```bash
pytest tests/test_providers.py -v
```

Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add python/mycli_ai/providers/
git add python/tests/
git commit -m "feat(python): add provider base interface"
```

---
## Task 3: OpenRouter Provider

**Files:**
- Create: `python/mycli_ai/providers/openrouter.py`
- Modify: `python/tests/test_providers.py`

- [ ] **Step 1: Write failing test for OpenRouter**

```python
# python/tests/test_providers.py

def test_openrouter_chat(mocker):
    from mycli_ai.providers.openrouter import OpenRouterProvider
    
    # Mock OpenAI client
    mock_response = {
        "choices": [{
            "message": {"content": "test response"}
        }],
        "usage": {
            "total_tokens": 150,
            "prompt_tokens": 50,
            "completion_tokens": 100
        }
    }
    
    mock_client = mocker.Mock()
    mock_client.chat.completions.create.return_value = mocker.Mock(**mock_response)
    mocker.patch('openai.OpenAI', return_value=mock_client)
    
    # Test
    provider = OpenRouterProvider(api_key="test-key")
    response = provider.chat("test message", model="deepseek/deepseek-r1")
    
    assert response.content == "test response"
    assert response.tokens == 150
    assert response.cost > 0

def test_openrouter_stream(mocker):
    from mycli_ai.providers.openrouter import OpenRouterProvider
    
    # Mock streaming response
    mock_chunks = [
        mocker.Mock(choices=[mocker.Mock(delta=mocker.Mock(content="Hello"))]),
        mocker.Mock(choices=[mocker.Mock(delta=mocker.Mock(content=" world"))]),
        mocker.Mock(choices=[mocker.Mock(delta=mocker.Mock(content="!"))]),
    ]
    
    mock_client = mocker.Mock()
    mock_client.chat.completions.create.return_value = iter(mock_chunks)
    mocker.patch('openai.OpenAI', return_value=mock_client)
    
    # Test
    provider = OpenRouterProvider(api_key="test-key")
    chunks = list(provider.stream_chat("test", model="deepseek/deepseek-r1"))
    
    assert chunks == ["Hello", " world", "!"]
```

- [ ] **Step 2: Run test to verify it fails**

```bash
pytest tests/test_providers.py::test_openrouter_chat -v
```

Expected: FAIL with "ModuleNotFoundError: No module named 'mycli_ai.providers.openrouter'"

- [ ] **Step 3: Implement OpenRouter provider**

```python
# python/mycli_ai/providers/openrouter.py
from typing import Optional, Dict, Any, Iterator
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
                metadata={
                    "model": model,
                    "prompt_tokens": response.usage.prompt_tokens,
                    "completion_tokens": response.usage.completion_tokens
                }
            )
        except Exception as e:
            raise RuntimeError(f"OpenRouter API error: {e}")
    
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
                stream=True,
                timeout=self.timeout
            )
            
            for chunk in stream:
                if chunk.choices[0].delta.content:
                    yield chunk.choices[0].delta.content
        except Exception as e:
            raise RuntimeError(f"OpenRouter streaming error: {e}")
    
    def _calculate_cost(self, model: str, prompt_tokens: int, completion_tokens: int) -> float:
        """Calculate cost based on token usage"""
        if model not in MODEL_PRICING:
            # Default pricing if model not in list
            return (prompt_tokens + completion_tokens) / 1_000_000 * 0.50
        
        pricing = MODEL_PRICING[model]
        input_cost = (prompt_tokens / 1_000_000) * pricing["input"]
        output_cost = (completion_tokens / 1_000_000) * pricing["output"]
        
        return input_cost + output_cost
```

- [ ] **Step 4: Run test to verify it passes**

```bash
pytest tests/test_providers.py::test_openrouter -v
```

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add python/mycli_ai/providers/openrouter.py
git add python/tests/test_providers.py
git commit -m "feat(python): add OpenRouter provider implementation"
```

---

## Task 4: Ollama Provider

**Files:**
- Create: `python/mycli_ai/providers/ollama.py`
- Modify: `python/tests/test_providers.py`

- [ ] **Step 1: Write failing test for Ollama**

```python
# python/tests/test_providers.py

def test_ollama_chat(mocker):
    from mycli_ai.providers.ollama import OllamaProvider
    
    # Mock ollama client
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
    from mycli_ai.providers.ollama import OllamaProvider
    
    # Mock streaming response
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
```

- [ ] **Step 2: Run test to verify it fails**

```bash
pytest tests/test_providers.py::test_ollama -v
```

Expected: FAIL with "ModuleNotFoundError: No module named 'mycli_ai.providers.ollama'"

- [ ] **Step 3: Implement Ollama provider**

```python
# python/mycli_ai/providers/ollama.py
from typing import Optional, Dict, Any, Iterator
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
            raise RuntimeError(f"Ollama error: {e}")
    
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
            raise RuntimeError(f"Ollama streaming error: {e}")
```

- [ ] **Step 4: Run test to verify it passes**

```bash
pytest tests/test_providers.py::test_ollama -v
```

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add python/mycli_ai/providers/ollama.py
git add python/tests/test_providers.py
git commit -m "feat(python): add Ollama provider implementation"
```

---
## Task 5: Response Parser

**Files:**
- Create: `python/mycli_ai/parser.py`
- Create: `python/tests/test_parser.py`

- [ ] **Step 1: Write failing test for parser**

```python
# python/tests/test_parser.py
import pytest
from mycli_ai.parser import extract_code_blocks, detect_file_changes, CodeBlock, FileChange

def test_extract_code_blocks():
    response = """
    Here's the fix:
    
    ```python
    def fixed_function():
        return True
    ```
    
    And here's the test:
    
    ```python
    def test_fixed():
        assert fixed_function() == True
    ```
    """
    
    blocks = extract_code_blocks(response)
    
    assert len(blocks) == 2
    assert blocks[0].language == "python"
    assert "fixed_function" in blocks[0].code
    assert blocks[1].language == "python"
    assert "test_fixed" in blocks[1].code

def test_extract_code_blocks_with_filename():
    response = """
    Update db.py:
    
    ```python
    # db.py
    def connect():
        return connection
    ```
    """
    
    blocks = extract_code_blocks(response)
    
    assert len(blocks) == 1
    assert blocks[0].language == "python"
    assert blocks[0].filename == "db.py"

def test_detect_file_changes():
    response = """
    I'll update these files:
    
    ```python
    # db.py
    def new_function():
        pass
    ```
    
    ```go
    // main.go
    func main() {
        fmt.Println("updated")
    }
    ```
    """
    
    changes = detect_file_changes(response)
    
    assert len(changes) == 2
    assert changes[0].file_path == "db.py"
    assert "new_function" in changes[0].content
    assert changes[1].file_path == "main.go"
    assert "updated" in changes[1].content
```

- [ ] **Step 2: Run test to verify it fails**

```bash
pytest tests/test_parser.py -v
```

Expected: FAIL with "ModuleNotFoundError: No module named 'mycli_ai.parser'"

- [ ] **Step 3: Implement parser**

```python
# python/mycli_ai/parser.py
import re
from dataclasses import dataclass
from typing import List, Optional

@dataclass
class CodeBlock:
    """Represents a code block from AI response"""
    language: str
    code: str
    filename: Optional[str] = None

@dataclass
class FileChange:
    """Represents a file modification"""
    file_path: str
    content: str
    language: str

def extract_code_blocks(text: str) -> List[CodeBlock]:
    """
    Extract code blocks from markdown text
    
    Supports:
    - ```language
    - ```language
      # filename.ext
    """
    blocks = []
    
    # Pattern: ```language\n(# filename)?\ncode\n```
    pattern = r'```(\w+)\n(.*?)```'
    matches = re.finditer(pattern, text, re.DOTALL)
    
    for match in matches:
        language = match.group(1)
        code = match.group(2).strip()
        
        # Check for filename in first line
        filename = None
        lines = code.split('\n')
        if lines:
            first_line = lines[0].strip()
            # Match: # filename.py or // filename.go
            filename_match = re.match(r'^[#/]+\s+(\S+\.\w+)', first_line)
            if filename_match:
                filename = filename_match.group(1)
        
        blocks.append(CodeBlock(
            language=language,
            code=code,
            filename=filename
        ))
    
    return blocks

def detect_file_changes(text: str) -> List[FileChange]:
    """
    Detect file modifications from AI response
    
    Looks for code blocks with filenames
    """
    changes = []
    blocks = extract_code_blocks(text)
    
    for block in blocks:
        if block.filename:
            changes.append(FileChange(
                file_path=block.filename,
                content=block.code,
                language=block.language
            ))
    
    return changes

def strip_filename_comment(code: str, language: str) -> str:
    """Remove filename comment from code"""
    lines = code.split('\n')
    if not lines:
        return code
    
    first_line = lines[0].strip()
    
    # Check if first line is a filename comment
    if language in ['python', 'ruby', 'perl', 'r']:
        if first_line.startswith('#') and '.' in first_line:
            return '\n'.join(lines[1:]).strip()
    elif language in ['go', 'rust', 'javascript', 'typescript', 'java', 'c', 'cpp']:
        if first_line.startswith('//') and '.' in first_line:
            return '\n'.join(lines[1:]).strip()
    
    return code
```

- [ ] **Step 4: Run test to verify it passes**

```bash
pytest tests/test_parser.py -v
```

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add python/mycli_ai/parser.py
git add python/tests/test_parser.py
git commit -m "feat(python): add response parser for code blocks and file changes"
```

---

## Task 6: Token Counting and Cost Tracking

**Files:**
- Create: `python/mycli_ai/tokens.py`
- Create: `python/mycli_ai/cost.py`
- Create: `python/tests/test_cost.py`

- [ ] **Step 1: Write failing test for token counting**

```python
# python/tests/test_cost.py
import pytest
from mycli_ai.tokens import count_tokens
from mycli_ai.cost import CostTracker

def test_count_tokens():
    text = "Hello, world! This is a test message."
    tokens = count_tokens(text)
    
    assert tokens > 0
    assert tokens < 20  # Should be around 10 tokens

def test_cost_tracker():
    tracker = CostTracker()
    
    # Add some usage
    tracker.add_usage(
        model="deepseek/deepseek-r1",
        prompt_tokens=100,
        completion_tokens=50
    )
    
    tracker.add_usage(
        model="deepseek/deepseek-r1",
        prompt_tokens=200,
        completion_tokens=100
    )
    
    # Check totals
    assert tracker.total_tokens == 450
    assert tracker.total_cost > 0
    
    # Check summary
    summary = tracker.get_summary()
    assert summary["total_tokens"] == 450
    assert summary["total_cost"] > 0
    assert "deepseek/deepseek-r1" in summary["by_model"]
```

- [ ] **Step 2: Run test to verify it fails**

```bash
pytest tests/test_cost.py -v
```

Expected: FAIL with "ModuleNotFoundError"

- [ ] **Step 3: Implement token counting**

```python
# python/mycli_ai/tokens.py
import tiktoken

def count_tokens(text: str, model: str = "gpt-4") -> int:
    """
    Count tokens in text using tiktoken
    
    Args:
        text: Text to count tokens for
        model: Model name (for encoding selection)
        
    Returns:
        Number of tokens
    """
    try:
        encoding = tiktoken.encoding_for_model(model)
    except KeyError:
        # Fallback to cl100k_base (GPT-4 encoding)
        encoding = tiktoken.get_encoding("cl100k_base")
    
    return len(encoding.encode(text))

def estimate_tokens(text: str) -> int:
    """
    Quick token estimation without tiktoken
    
    Rough estimate: 1 token ≈ 4 characters
    """
    return len(text) // 4
```

- [ ] **Step 4: Implement cost tracking**

```python
# python/mycli_ai/cost.py
from typing import Dict, List
from dataclasses import dataclass, field

# Model pricing (per 1M tokens)
MODEL_PRICING = {
    "deepseek/deepseek-r1": {"input": 0.14, "output": 0.14},
    "stepfun/step-2-16k": {"input": 0.10, "output": 0.10},
    "anthropic/claude-sonnet-4": {"input": 3.00, "output": 15.00},
}

@dataclass
class UsageRecord:
    """Single usage record"""
    model: str
    prompt_tokens: int
    completion_tokens: int
    cost: float

class CostTracker:
    """Track token usage and costs across requests"""
    
    def __init__(self):
        self.records: List[UsageRecord] = []
        self.total_tokens = 0
        self.total_cost = 0.0
    
    def add_usage(self, model: str, prompt_tokens: int, completion_tokens: int):
        """Add a usage record"""
        cost = self._calculate_cost(model, prompt_tokens, completion_tokens)
        
        record = UsageRecord(
            model=model,
            prompt_tokens=prompt_tokens,
            completion_tokens=completion_tokens,
            cost=cost
        )
        
        self.records.append(record)
        self.total_tokens += prompt_tokens + completion_tokens
        self.total_cost += cost
    
    def get_summary(self) -> Dict:
        """Get usage summary"""
        by_model = {}
        
        for record in self.records:
            if record.model not in by_model:
                by_model[record.model] = {
                    "tokens": 0,
                    "cost": 0.0,
                    "requests": 0
                }
            
            by_model[record.model]["tokens"] += record.prompt_tokens + record.completion_tokens
            by_model[record.model]["cost"] += record.cost
            by_model[record.model]["requests"] += 1
        
        return {
            "total_tokens": self.total_tokens,
            "total_cost": self.total_cost,
            "total_requests": len(self.records),
            "by_model": by_model
        }
    
    def _calculate_cost(self, model: str, prompt_tokens: int, completion_tokens: int) -> float:
        """Calculate cost for usage"""
        if model not in MODEL_PRICING:
            # Default pricing
            return (prompt_tokens + completion_tokens) / 1_000_000 * 0.50
        
        pricing = MODEL_PRICING[model]
        input_cost = (prompt_tokens / 1_000_000) * pricing["input"]
        output_cost = (completion_tokens / 1_000_000) * pricing["output"]
        
        return input_cost + output_cost
```

- [ ] **Step 5: Run test to verify it passes**

```bash
pytest tests/test_cost.py -v
```

Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add python/mycli_ai/tokens.py
git add python/mycli_ai/cost.py
git add python/tests/test_cost.py
git commit -m "feat(python): add token counting and cost tracking"
```

---
## Task 7: JSON-RPC Server

**Files:**
- Create: `python/mycli_ai/server.py`
- Create: `python/tests/test_server.py`

- [ ] **Step 1: Write failing test for JSON-RPC server**

```python
# python/tests/test_server.py
import pytest
import json
from mycli_ai.server import JSONRPCServer, handle_request

def test_handle_chat_request(mocker):
    # Mock provider
    mock_provider = mocker.Mock()
    mock_provider.chat.return_value = mocker.Mock(
        content="test response",
        tokens=100,
        cost=0.001
    )
    
    request = {
        "jsonrpc": "2.0",
        "id": 1,
        "method": "chat",
        "params": {
            "message": "test message",
            "model": "deepseek/deepseek-r1"
        }
    }
    
    response = handle_request(request, mock_provider)
    
    assert response["jsonrpc"] == "2.0"
    assert response["id"] == 1
    assert "result" in response
    assert response["result"]["content"] == "test response"
    assert response["result"]["tokens"] == 100

def test_handle_invalid_method():
    request = {
        "jsonrpc": "2.0",
        "id": 1,
        "method": "invalid_method",
        "params": {}
    }
    
    response = handle_request(request, None)
    
    assert "error" in response
    assert response["error"]["code"] == -32601  # Method not found

def test_handle_malformed_request():
    request = {"invalid": "request"}
    
    response = handle_request(request, None)
    
    assert "error" in response
    assert response["error"]["code"] == -32600  # Invalid request
```

- [ ] **Step 2: Run test to verify it fails**

```bash
pytest tests/test_server.py -v
```

Expected: FAIL with "ModuleNotFoundError"

- [ ] **Step 3: Implement JSON-RPC server**

```python
# python/mycli_ai/server.py
import sys
import json
from typing import Dict, Any, Optional
from mycli_ai.providers.openrouter import OpenRouterProvider
from mycli_ai.providers.ollama import OllamaProvider
from mycli_ai.cost import CostTracker

# JSON-RPC error codes
PARSE_ERROR = -32700
INVALID_REQUEST = -32600
METHOD_NOT_FOUND = -32601
INVALID_PARAMS = -32602
INTERNAL_ERROR = -32603

class JSONRPCServer:
    """JSON-RPC server for AI provider communication"""
    
    def __init__(self, config: Dict[str, Any]):
        self.config = config
        self.cost_tracker = CostTracker()
        self.providers = self._init_providers()
    
    def _init_providers(self) -> Dict[str, Any]:
        """Initialize AI providers based on config"""
        providers = {}
        
        # OpenRouter
        if self.config.get("openrouter", {}).get("enabled"):
            api_key = self.config["openrouter"].get("api_key")
            if api_key:
                providers["openrouter"] = OpenRouterProvider(
                    api_key=api_key,
                    base_url=self.config["openrouter"].get("base_url", "https://openrouter.ai/api/v1")
                )
        
        # Ollama
        if self.config.get("ollama", {}).get("enabled"):
            providers["ollama"] = OllamaProvider(
                base_url=self.config["ollama"].get("base_url", "http://localhost:11434")
            )
        
        return providers
    
    def run(self):
        """Run the JSON-RPC server (stdin/stdout)"""
        for line in sys.stdin:
            try:
                request = json.loads(line.strip())
                response = self.handle_request(request)
                print(json.dumps(response), flush=True)
            except json.JSONDecodeError:
                error_response = create_error_response(None, PARSE_ERROR, "Parse error")
                print(json.dumps(error_response), flush=True)
            except Exception as e:
                error_response = create_error_response(None, INTERNAL_ERROR, str(e))
                print(json.dumps(error_response), flush=True)
    
    def handle_request(self, request: Dict[str, Any]) -> Dict[str, Any]:
        """Handle a single JSON-RPC request"""
        # Validate request
        if not isinstance(request, dict):
            return create_error_response(None, INVALID_REQUEST, "Invalid request")
        
        if "jsonrpc" not in request or request["jsonrpc"] != "2.0":
            return create_error_response(None, INVALID_REQUEST, "Invalid JSON-RPC version")
        
        if "method" not in request:
            return create_error_response(request.get("id"), INVALID_REQUEST, "Missing method")
        
        method = request["method"]
        params = request.get("params", {})
        request_id = request.get("id")
        
        # Route to method handler
        if method == "chat":
            return self._handle_chat(request_id, params)
        elif method == "stream_chat":
            return self._handle_stream_chat(request_id, params)
        elif method == "get_cost_summary":
            return self._handle_get_cost_summary(request_id)
        else:
            return create_error_response(request_id, METHOD_NOT_FOUND, f"Method not found: {method}")
    
    def _handle_chat(self, request_id: Any, params: Dict[str, Any]) -> Dict[str, Any]:
        """Handle chat request"""
        try:
            message = params.get("message")
            model = params.get("model")
            
            if not message or not model:
                return create_error_response(request_id, INVALID_PARAMS, "Missing message or model")
            
            # Get provider
            provider = self._get_provider_for_model(model)
            if not provider:
                return create_error_response(request_id, INTERNAL_ERROR, "No provider available")
            
            # Make request
            response = provider.chat(message, model, **params)
            
            # Track cost
            if response.metadata:
                self.cost_tracker.add_usage(
                    model=model,
                    prompt_tokens=response.metadata.get("prompt_tokens", 0),
                    completion_tokens=response.metadata.get("completion_tokens", 0)
                )
            
            return create_success_response(request_id, {
                "content": response.content,
                "tokens": response.tokens,
                "cost": response.cost,
                "metadata": response.metadata
            })
        except Exception as e:
            return create_error_response(request_id, INTERNAL_ERROR, str(e))
    
    def _handle_stream_chat(self, request_id: Any, params: Dict[str, Any]) -> Dict[str, Any]:
        """Handle streaming chat request"""
        # Streaming handled differently - send chunks as notifications
        try:
            message = params.get("message")
            model = params.get("model")
            
            if not message or not model:
                return create_error_response(request_id, INVALID_PARAMS, "Missing message or model")
            
            provider = self._get_provider_for_model(model)
            if not provider:
                return create_error_response(request_id, INTERNAL_ERROR, "No provider available")
            
            # Send chunks
            for chunk in provider.stream_chat(message, model, **params):
                notification = {
                    "jsonrpc": "2.0",
                    "method": "chunk",
                    "params": {"content": chunk}
                }
                print(json.dumps(notification), flush=True)
            
            # Send completion
            return create_success_response(request_id, {"status": "completed"})
        except Exception as e:
            return create_error_response(request_id, INTERNAL_ERROR, str(e))
    
    def _handle_get_cost_summary(self, request_id: Any) -> Dict[str, Any]:
        """Handle cost summary request"""
        summary = self.cost_tracker.get_summary()
        return create_success_response(request_id, summary)
    
    def _get_provider_for_model(self, model: str):
        """Get appropriate provider for model"""
        # Check if model is for OpenRouter
        if "/" in model:  # OpenRouter models have format "provider/model"
            return self.providers.get("openrouter")
        else:
            # Assume local Ollama model
            return self.providers.get("ollama")

def handle_request(request: Dict[str, Any], provider: Any) -> Dict[str, Any]:
    """Standalone request handler for testing"""
    if not isinstance(request, dict):
        return create_error_response(None, INVALID_REQUEST, "Invalid request")
    
    if "method" not in request:
        return create_error_response(request.get("id"), INVALID_REQUEST, "Missing method")
    
    method = request["method"]
    request_id = request.get("id")
    params = request.get("params", {})
    
    if method == "chat":
        try:
            message = params.get("message")
            model = params.get("model")
            response = provider.chat(message, model)
            
            return create_success_response(request_id, {
                "content": response.content,
                "tokens": response.tokens,
                "cost": response.cost
            })
        except Exception as e:
            return create_error_response(request_id, INTERNAL_ERROR, str(e))
    else:
        return create_error_response(request_id, METHOD_NOT_FOUND, f"Method not found: {method}")

def create_success_response(request_id: Any, result: Any) -> Dict[str, Any]:
    """Create JSON-RPC success response"""
    return {
        "jsonrpc": "2.0",
        "id": request_id,
        "result": result
    }

def create_error_response(request_id: Any, code: int, message: str) -> Dict[str, Any]:
    """Create JSON-RPC error response"""
    return {
        "jsonrpc": "2.0",
        "id": request_id,
        "error": {
            "code": code,
            "message": message
        }
    }

def main():
    """Entry point for JSON-RPC server"""
    # Read config from first line
    config_line = sys.stdin.readline()
    config = json.loads(config_line.strip())
    
    # Start server
    server = JSONRPCServer(config)
    server.run()

if __name__ == "__main__":
    main()
```

- [ ] **Step 4: Run test to verify it passes**

```bash
pytest tests/test_server.py -v
```

Expected: PASS

- [ ] **Step 5: Test server manually**

```bash
# Create test config
echo '{"openrouter": {"enabled": true, "api_key": "test-key"}}' > /tmp/test-config.json

# Start server (will wait for input)
python -m mycli_ai.server < /tmp/test-config.json
```

- [ ] **Step 6: Commit**

```bash
git add python/mycli_ai/server.py
git add python/tests/test_server.py
git commit -m "feat(python): add JSON-RPC server for Go communication"
```

---
## Task 8: Integration Test

**Files:**
- Create: `python/tests/test_integration.py`

- [ ] **Step 1: Write integration test**

```python
# python/tests/test_integration.py
import pytest
from mycli_ai.providers.openrouter import OpenRouterProvider
from mycli_ai.parser import extract_code_blocks, detect_file_changes
from mycli_ai.cost import CostTracker

def test_full_flow_with_mocks(mocker):
    """Test complete flow from request to parsed response"""
    
    # Mock OpenAI response
    mock_response = mocker.Mock()
    mock_response.choices = [mocker.Mock(message=mocker.Mock(content="""
Here's the fix:

```python
# db.py
def fixed_query():
    return "SELECT * FROM users"
```
"""))]
    mock_response.usage = mocker.Mock(
        total_tokens=150,
        prompt_tokens=50,
        completion_tokens=100
    )
    
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
    
    # Step 4: Parse response
    blocks = extract_code_blocks(response.content)
    assert len(blocks) == 1
    assert blocks[0].language == "python"
    
    # Step 5: Detect file changes
    changes = detect_file_changes(response.content)
    assert len(changes) == 1
    assert changes[0].file_path == "db.py"
    assert "fixed_query" in changes[0].content
    
    # Step 6: Track cost
    tracker = CostTracker()
    tracker.add_usage(
        model="deepseek/deepseek-r1",
        prompt_tokens=response.metadata["prompt_tokens"],
        completion_tokens=response.metadata["completion_tokens"]
    )
    
    summary = tracker.get_summary()
    assert summary["total_tokens"] == 150
    assert summary["total_cost"] > 0
```

- [ ] **Step 2: Run integration test**

```bash
pytest tests/test_integration.py -v
```

Expected: PASS

- [ ] **Step 3: Run all tests**

```bash
pytest tests/ -v --cov=mycli_ai
```

Expected: All tests PASS with good coverage

- [ ] **Step 4: Commit**

```bash
git add python/tests/test_integration.py
git commit -m "test(python): add integration test for full AI flow"
```

---

## Completion Checklist

- [x] Python project initialized with dependencies
- [x] Provider base interface defined
- [x] OpenRouter provider implemented
- [x] Ollama provider implemented
- [x] Response parser (code blocks, file changes)
- [x] Token counting with tiktoken
- [x] Cost tracking system
- [x] JSON-RPC server for Go communication
- [x] Integration tests
- [x] All tests passing

## Manual Verification

```bash
cd python

# Activate venv
source venv/bin/activate

# Run all tests
pytest tests/ -v --cov=mycli_ai

# Test OpenRouter provider (requires API key)
python -c "
from mycli_ai.providers.openrouter import OpenRouterProvider
provider = OpenRouterProvider(api_key='your-key')
response = provider.chat('Hello', model='deepseek/deepseek-r1')
print(response.content)
"

# Test Ollama provider (requires Ollama running)
python -c "
from mycli_ai.providers.ollama import OllamaProvider
provider = OllamaProvider()
response = provider.chat('Hello', model='codellama:13b')
print(response.content)
"

# Test JSON-RPC server
echo '{"openrouter": {"enabled": true, "api_key": "test"}}' | python -m mycli_ai.server
```

## Next Steps

This plan establishes the Python AI layer. Next plans will build on this:
- **Plan C**: Go-Python Bridge (subprocess management, JSON-RPC client)
- **Plan D**: Basic TUI (Bubble Tea interface)
- **Plan E**: TUI Polish (colors, animations, diff preview)

---

## Notes for Implementation

**Key Design Decisions:**
1. JSON-RPC over stdin/stdout for Go-Python communication
2. Provider abstraction allows easy addition of new AI services
3. Cost tracking built-in from the start
4. Response parsing handles code blocks and file modifications
5. Streaming support for real-time responses

**Testing Strategy:**
- Mock external APIs (OpenRouter, Ollama) in tests
- Integration test covers full flow
- Use pytest-mock for clean mocking
- Skip tests requiring external services in CI

**Common Issues:**
- OpenRouter API key must be valid for real tests
- Ollama must be running locally for Ollama tests
- tiktoken may need internet for first-time encoding download
- JSON-RPC requires exact format (jsonrpc: "2.0")

**Performance Considerations:**
- Token counting with tiktoken is fast (<1ms)
- Streaming reduces perceived latency
- Cost calculation is O(1) per request
- Parser uses regex (fast for typical responses)
