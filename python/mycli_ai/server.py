"""JSON-RPC server for AI provider communication."""
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
        self.provider = self._initialize_provider()

    def _initialize_provider(self):
        """Initialize the appropriate provider based on config"""
        provider_type = self.config.get("provider", "openrouter")

        if provider_type == "openrouter":
            return OpenRouterProvider(
                api_key=self.config.get("api_key"),
                base_url=self.config.get("base_url", "https://openrouter.ai/api/v1")
            )
        elif provider_type == "ollama":
            return OllamaProvider(
                base_url=self.config.get("base_url", "http://localhost:11434")
            )
        else:
            raise ValueError(f"Unknown provider: {provider_type}")

    def run(self):
        """Run the JSON-RPC server (stdin/stdout)"""
        for line in sys.stdin:
            try:
                request = json.loads(line)
                response = handle_request(request, self.provider)
                print(json.dumps(response), flush=True)
            except json.JSONDecodeError:
                error_response = create_error_response(None, PARSE_ERROR, "Parse error")
                print(json.dumps(error_response), flush=True)
            except Exception as e:
                error_response = create_error_response(None, INTERNAL_ERROR, str(e))
                print(json.dumps(error_response), flush=True)


def handle_request(request: Dict[str, Any], provider) -> Dict[str, Any]:
    """
    Handle a JSON-RPC request

    Args:
        request: JSON-RPC request object
        provider: AI provider instance

    Returns:
        JSON-RPC response object
    """
    # Validate request structure
    if not isinstance(request, dict):
        return create_error_response(None, INVALID_REQUEST, "Invalid request")

    if "jsonrpc" not in request or request["jsonrpc"] != "2.0":
        return create_error_response(None, INVALID_REQUEST, "Invalid JSON-RPC version")

    if "method" not in request:
        return create_error_response(
            request.get("id"),
            INVALID_REQUEST,
            "Missing method"
        )

    request_id = request.get("id")
    method = request["method"]
    params = request.get("params", {})

    # Handle methods
    try:
        if method == "chat":
            return handle_chat(request_id, params, provider)
        elif method == "stream_chat":
            return handle_stream_chat(request_id, params, provider)
        else:
            return create_error_response(
                request_id,
                METHOD_NOT_FOUND,
                f"Method not found: {method}"
            )
    except KeyError as e:
        return create_error_response(
            request_id,
            INVALID_PARAMS,
            f"Missing parameter: {str(e)}"
        )
    except Exception as e:
        return create_error_response(
            request_id,
            INTERNAL_ERROR,
            str(e)
        )


def handle_chat(request_id: Any, params: Dict[str, Any], provider) -> Dict[str, Any]:
    """Handle chat method"""
    if "message" not in params or "model" not in params:
        return create_error_response(
            request_id,
            INVALID_PARAMS,
            "Missing required parameters: message, model"
        )

    message = params["message"]
    model = params["model"]
    kwargs = {k: v for k, v in params.items() if k not in ["message", "model"]}

    response = provider.chat(message, model, **kwargs)

    return {
        "jsonrpc": "2.0",
        "id": request_id,
        "result": {
            "content": response.content,
            "tokens": response.tokens,
            "cost": response.cost,
            "metadata": response.metadata
        }
    }


def handle_stream_chat(request_id: Any, params: Dict[str, Any], provider) -> Dict[str, Any]:
    """Handle stream_chat method"""
    if "message" not in params or "model" not in params:
        return create_error_response(
            request_id,
            INVALID_PARAMS,
            "Missing required parameters: message, model"
        )

    message = params["message"]
    model = params["model"]
    kwargs = {k: v for k, v in params.items() if k not in ["message", "model"]}

    chunks = list(provider.stream_chat(message, model, **kwargs))

    return {
        "jsonrpc": "2.0",
        "id": request_id,
        "result": {
            "chunks": chunks
        }
    }


def create_error_response(request_id: Any, code: int, message: str) -> Dict[str, Any]:
    """Create a JSON-RPC error response"""
    return {
        "jsonrpc": "2.0",
        "id": request_id,
        "error": {
            "code": code,
            "message": message
        }
    }
