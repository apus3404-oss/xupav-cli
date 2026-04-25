"""Tests for response parser."""
import pytest
from mycli_ai.parser import extract_code_blocks, detect_file_changes, CodeBlock, FileChange


def test_extract_code_blocks():
    """Test extracting multiple code blocks"""
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
    """Test extracting code blocks with filename comments"""
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
    """Test detecting file changes from response"""
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


def test_extract_code_blocks_no_language():
    """Test extracting code blocks without language specified"""
    response = """
    ```
    plain code block
    ```
    """

    blocks = extract_code_blocks(response)

    assert len(blocks) == 1
    assert blocks[0].language == ""
    assert "plain code block" in blocks[0].code


def test_extract_code_blocks_empty():
    """Test with no code blocks"""
    response = "Just plain text with no code blocks"

    blocks = extract_code_blocks(response)

    assert len(blocks) == 0


def test_detect_file_changes_no_files():
    """Test detecting file changes when no filenames present"""
    response = """
    ```python
    def function():
        pass
    ```
    """

    changes = detect_file_changes(response)

    assert len(changes) == 0
