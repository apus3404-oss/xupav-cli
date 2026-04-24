# TUI Polish Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Enhance the TUI with modern visual elements including gradient colors, animations, syntax highlighting, diff preview, and toast notifications.

**Architecture:** Build on Plan D's basic TUI by adding visual components, animation system, and enhanced styling using Lipgloss and Chroma.

**Tech Stack:**
- **charmbracelet/lipgloss** - Advanced styling
- **alecthomas/chroma** - Syntax highlighting
- **charmbracelet/bubbles/spinner** - Loading animations

---

## File Structure

**New Files:**
```
internal/
  tui/
    components/
      diff.go          # Diff preview component
      notification.go  # Toast notifications
      spinner.go       # Custom spinner
    theme.go           # Color themes
    animations.go      # Animation helpers
```

**Modified Files:**
```
internal/tui/styles.go
internal/tui/view.go
internal/tui/update.go
```

---

## Summary

This plan adds 5 tasks:
1. Enhanced color themes with gradients
2. Syntax highlighting for code blocks
3. Diff preview component
4. Toast notifications
5. Loading animations and spinners

Total estimated time: 4-6 hours
