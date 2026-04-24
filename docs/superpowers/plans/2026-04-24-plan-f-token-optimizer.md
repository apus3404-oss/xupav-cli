# Token Optimizer Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement the 3-layer token optimization system that achieves 60-70% token savings - THE core differentiator of mycli.

**Architecture:** Three independent optimization layers that work together: file chunking (AST-based), context prioritization (LSP + scoring), and diff-based transmission (hash tracking).

**Tech Stack:**
- **tree-sitter** (via Python) - AST parsing
- **tiktoken** (Python) - Token counting
- Go cache system for file hashes

---

## File Structure

```
internal/
  core/
    optimizer.go       # Main optimizer orchestrator
    chunker.go         # Layer 1: File chunking
    prioritizer.go     # Layer 2: Context prioritization
    differ.go          # Layer 3: Diff-based transmission
    cache.go           # File hash cache
    budget.go          # Token budget management
```

---

## Task Summary

### Task 1: File Hash Cache System
- SHA-256 hashing for files
- In-memory cache with LRU eviction
- Detect file changes efficiently
- **Test:** Hash same file twice, verify cache hit

### Task 2: Layer 1 - File Chunking (AST-based)
- Call Python tree-sitter via bridge
- Extract functions/classes/methods
- Score chunks by keyword relevance
- **Test:** Parse Python file, verify chunks extracted

### Task 3: Layer 2 - Context Prioritization
- Keyword extraction from user query
- File scoring algorithm (4 factors)
- Threshold-based filtering
- **Test:** Score files for "database query", verify ranking

### Task 4: Layer 3 - Diff-Based Transmission
- Track sent files by hash
- Generate unified diffs for changed files
- Send only diffs on subsequent mentions
- **Test:** Send file twice, verify diff on second send

### Task 5: Token Budget System
- Per-model token limits
- Budget allocation (system, user, history, files)
- Overflow handling (aggressive filtering)
- **Test:** Exceed budget, verify filtering kicks in

### Task 6: Integration and Metrics
- Wire optimizer into chat flow
- Track savings metrics
- Display savings in TUI
- **Test:** Full conversation, measure actual savings

---

## Expected Results

**Token Savings:** 60-70% compared to naive approach
**Performance:** <500ms for 10 files
**Accuracy:** No loss of relevant context

---

## Next: Plan G (LSP Integration)
