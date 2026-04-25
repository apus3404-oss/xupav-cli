"""Response parser for extracting code blocks and file changes."""
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
    - ```
      code without language
    """
    blocks = []
    pattern = r'```(\w*)\n(.*?)```'
    matches = re.findall(pattern, text, re.DOTALL)

    for language, code in matches:
        code = code.strip()
        filename = None

        # Try to extract filename from first line comment
        first_line = code.split('\n')[0].strip()

        # Python/Shell style: # filename.py
        if first_line.startswith('#') and not first_line.startswith('##'):
            potential_filename = first_line[1:].strip()
            if '.' in potential_filename and ' ' not in potential_filename:
                filename = potential_filename

        # C/Go/Java style: // filename.go
        elif first_line.startswith('//'):
            potential_filename = first_line[2:].strip()
            if '.' in potential_filename and ' ' not in potential_filename:
                filename = potential_filename

        blocks.append(CodeBlock(
            language=language,
            code=code,
            filename=filename
        ))

    return blocks


def detect_file_changes(text: str) -> List[FileChange]:
    """
    Detect file modifications from AI response

    Returns FileChange objects for code blocks that have filenames
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
