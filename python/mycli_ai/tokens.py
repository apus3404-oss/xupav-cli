"""Token counting utilities using tiktoken."""
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
    if not text:
        return 0

    try:
        encoding = tiktoken.encoding_for_model(model)
    except KeyError:
        # Fallback to cl100k_base for unknown models
        encoding = tiktoken.get_encoding("cl100k_base")

    tokens = encoding.encode(text)
    return len(tokens)
