from setuptools import setup, find_packages

setup(
    name="mycli-ai",
    version="0.1.0",
    packages=find_packages(),
    install_requires=[
        "openai==1.12.0",
        "ollama==0.1.7",
        "tiktoken==0.6.0",
    ],
    extras_require={
        "dev": [
            "pytest==8.0.0",
            "pytest-mock==3.12.0",
            "pytest-asyncio==0.23.5",
        ]
    },
    python_requires=">=3.9",
)
