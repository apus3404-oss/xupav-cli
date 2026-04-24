# Code Modifications Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement safe code modification system with diff preview, backup, undo, and three operation modes (safe, interactive, auto).

**Architecture:** File writer with backup system, diff generator, approval flow, and git integration. Respects user's chosen code_mode from config.

**Tech Stack:**
- **go-git** - Git operations
- **go-diff** - Diff generation
- Native file I/O with atomic writes

---

## File Structure

```
internal/
  writer/
    file.go          # File writer with backup
    diff.go          # Diff generation and parsing
    undo.go          # Undo mechanism
    git.go           # Git staging operations
    approval.go      # Interactive approval flow
```

---

## Task Summary

### Task 1: Backup System
- Create `~/.mycli/backups/<timestamp>/` on write
- Copy original file before modification
- Maintain backup history (last 10)
- **Test:** Write file, verify backup created

### Task 2: Diff Generation
- Parse AI response for file changes
- Generate unified diff format
- Color-coded display (red -, green +)
- **Test:** Generate diff, verify format

### Task 3: Interactive Approval Flow
- Show diff in TUI
- Prompt: [y]es [n]o [d]etailed [e]dit [a]ll
- Handle each option appropriately
- **Test:** Mock approval, verify each option

### Task 4: File Writer
- Atomic writes (write to temp, then rename)
- Preserve file permissions
- Handle write errors gracefully
- **Test:** Write file, verify atomicity

### Task 5: Undo Mechanism
- `mycli undo` command
- Restore from latest backup
- Unstage from git if staged
- **Test:** Write, undo, verify restoration

### Task 6: Git Integration
- Auto-stage modified files (if in repo)
- Respect .gitignore
- Handle git errors gracefully
- **Test:** Modify file, verify git status

### Task 7: Mode Implementation
- Safe mode: only show suggestions
- Interactive mode: diff + approval
- Auto mode: apply immediately
- **Test:** Each mode with same change

---

## Expected Results

**Safety:** No data loss (backups always created)
**UX:** Clear diff preview, easy approval
**Git:** Seamless integration with workflow

---

## All Plans Complete!

**Total:** 8 implementation plans
**Estimated Time:** 6-8 weeks for full MVP
**Next Step:** Choose execution approach
