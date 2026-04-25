"""Tests for JSON-RPC server."""
import pytest
import json
from mycli_ai.server import JSONRPCServer, handle_request


def test_handle_chat_request(mocker):
    """Test handling a chat request"""
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
    """Test handling invalid method"""
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
    """Test handling malformed request"""
    request = {"invalid": "request"}

    response = handle_request(request, None)

    assert "error" in response
    assert response["error"]["code"] == -32600  # Invalid request


def test_handle_stream_chat_request(mocker):
    """Test handling a streaming chat request"""
    # Mock provider
    mock_provider = mocker.Mock()
    mock_provider.stream_chat.return_value = iter(["Hello", " world", "!"])

    request = {
        "jsonrpc": "2.0",
        "id": 2,
        "method": "stream_chat",
        "params": {
            "message": "test message",
            "model": "deepseek/deepseek-r1"
        }
    }

    response = handle_request(request, mock_provider)

    assert response["jsonrpc"] == "2.0"
    assert response["id"] == 2
    assert "result" in response
    assert "chunks" in response["result"]


def test_handle_missing_params():
    """Test handling request with missing params"""
    request = {
        "jsonrpc": "2.0",
        "id": 3,
        "method": "chat",
        "params": {}
    }

    response = handle_request(request, None)

    assert "error" in response
    assert response["error"]["code"] == -32602  # Invalid params


def test_jsonrpc_server_initialization():
    """Test JSON-RPC server initialization"""
    config = {
        "provider": "openrouter",
        "api_key": "test-key",
        "model": "deepseek/deepseek-r1"
    }

    server = JSONRPCServer(config)

    assert server.config == config
    assert server.cost_tracker is not None
