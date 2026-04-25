"""Cost tracking for AI API usage."""
from typing import Dict, Any
from collections import defaultdict

# Model pricing (per 1M tokens)
MODEL_PRICING = {
    "deepseek/deepseek-r1": {"input": 0.14, "output": 0.14},
    "stepfun/step-2-16k": {"input": 0.10, "output": 0.10},
    "anthropic/claude-sonnet-4": {"input": 3.00, "output": 15.00},
}


class CostTracker:
    """Track token usage and costs across API calls"""

    def __init__(self):
        self.total_tokens = 0
        self.total_cost = 0.0
        self.usage_by_model = defaultdict(lambda: {
            "prompt_tokens": 0,
            "completion_tokens": 0,
            "total_tokens": 0,
            "cost": 0.0
        })

    def add_usage(self, model: str, prompt_tokens: int, completion_tokens: int):
        """
        Add usage for a model

        Args:
            model: Model identifier
            prompt_tokens: Number of prompt tokens
            completion_tokens: Number of completion tokens
        """
        total = prompt_tokens + completion_tokens
        cost = self._calculate_cost(model, prompt_tokens, completion_tokens)

        self.total_tokens += total
        self.total_cost += cost

        self.usage_by_model[model]["prompt_tokens"] += prompt_tokens
        self.usage_by_model[model]["completion_tokens"] += completion_tokens
        self.usage_by_model[model]["total_tokens"] += total
        self.usage_by_model[model]["cost"] += cost

    def get_summary(self) -> Dict[str, Any]:
        """
        Get usage summary

        Returns:
            Dictionary with total_tokens, total_cost, and by_model breakdown
        """
        return {
            "total_tokens": self.total_tokens,
            "total_cost": round(self.total_cost, 6),
            "by_model": dict(self.usage_by_model)
        }

    def reset(self):
        """Reset all tracking data"""
        self.total_tokens = 0
        self.total_cost = 0.0
        self.usage_by_model.clear()

    def _calculate_cost(self, model: str, prompt_tokens: int, completion_tokens: int) -> float:
        """Calculate cost based on token usage"""
        if model not in MODEL_PRICING:
            # Default pricing if model not found
            return (prompt_tokens + completion_tokens) * 0.0001 / 1000

        pricing = MODEL_PRICING[model]
        input_cost = (prompt_tokens / 1_000_000) * pricing["input"]
        output_cost = (completion_tokens / 1_000_000) * pricing["output"]

        return input_cost + output_cost
