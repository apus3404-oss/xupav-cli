"""Base provider interface for AI services."""
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
