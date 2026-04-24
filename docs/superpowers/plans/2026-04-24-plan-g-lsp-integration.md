# LSP Integration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Integrate Language Server Protocol clients for Go, Python, JavaScript/TypeScript to enable intelligent file discovery and context building.

**Architecture:** LSP client manager that spawns language servers, queries for definitions/references, and uses results to build smart context for AI requests.

**Tech Stack:**
- **go-lsp** - LSP client library
- Language servers: gopls (Go), pyright (Python), typescript-language-server (JS/TS)

---

## File Structure

```
internal/
  core/
    lsp/
      client.go        # LSP client interface
      manager.go       # Multi-language LSP manager
      go.go            # Go LSP (gopls)
      python.go        # Python LSP (pyright)
      javascript.go    # JS/TS LSP
      cache.go         # LSP result cache
```

---

## Task Summary

### Task 1: LSP Client Interface
- Define common LSP operations (definition, references, hover)
- Abstract protocol details
- Timeout and error handling
- **Test:** Mock LSP server, verify protocol

### Task 2: Go LSP Client (gopls)
- Start gopls subprocess
- Initialize with workspace root
- Query definitions and references
- **Test:** Find definition of Go function

### Task 3: Python LSP Client (pyright)
- Start pyright subprocess
- Handle Python-specific features
- Parse pyright responses
- **Test:** Find references to Python class

### Task 4: JavaScript/TypeScript LSP
- Start typescript-language-server
- Handle both JS and TS files
- Support import resolution
- **Test:** Find definition across files

### Task 5: LSP Manager
- Detect language from file extension
- Route requests to appropriate LSP
- Cache results (1-hour TTL)
- **Test:** Query multiple languages in one session

### Task 6: Integration with Context Builder
- Use LSP to find related files
- "Go to definition" for imports
- "Find references" for symbols
- **Test:** User mentions function, verify related files found

---

## Expected Results

**Accuracy:** 90%+ relevant file detection
**Performance:** <200ms per LSP query
**Coverage:** Go, Python, JS/TS (Rust future)

---

## Next: Plan H (Code Modifications)
