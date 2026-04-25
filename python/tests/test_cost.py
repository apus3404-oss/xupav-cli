"""Tests for token counting and cost tracking."""
import pytest
from mycli_ai.tokens import count_tokens
from mycli_ai.cost import CostTracker


def test_count_tokens():
    """Test token counting"""
    text = "Hello, world! This is a test message."
    tokens = count_tokens(text)

    assert tokens > 0
    assert tokens < 20  # Should be around 10 tokens


def test_count_tokens_empty():
    """Test token counting with empty string"""
    tokens = count_tokens("")
    assert tokens == 0


def test_count_tokens_long_text():
    """Test token counting with longer text"""
    text = "This is a longer text. " * 100
    tokens = count_tokens(text)
    assert tokens > 100


def test_cost_tracker():
    """Test cost tracker basic functionality"""
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


def test_cost_tracker_multiple_models():
    """Test cost tracker with multiple models"""
    tracker = CostTracker()

    tracker.add_usage(
        model="deepseek/deepseek-r1",
        prompt_tokens=100,
        completion_tokens=50
    )

    tracker.add_usage(
        model="anthropic/claude-sonnet-4",
        prompt_tokens=100,
        completion_tokens=50
    )

    summary = tracker.get_summary()
    assert len(summary["by_model"]) == 2
    assert "deepseek/deepseek-r1" in summary["by_model"]
    assert "anthropic/claude-sonnet-4" in summary["by_model"]


def test_cost_tracker_reset():
    """Test cost tracker reset functionality"""
    tracker = CostTracker()

    tracker.add_usage(
        model="deepseek/deepseek-r1",
        prompt_tokens=100,
        completion_tokens=50
    )

    assert tracker.total_tokens == 150

    tracker.reset()

    assert tracker.total_tokens == 0
    assert tracker.total_cost == 0
    assert len(tracker.get_summary()["by_model"]) == 0
