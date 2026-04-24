
## Task 1: Enhanced Color Themes

**Summary:** Add gradient colors, theme system, and terminal capability detection.

**Files:** `internal/tui/theme.go`, `internal/tui/styles.go`

**Key Steps:**
- Detect terminal color support (256-color, true-color)
- Implement gradient header (cyan → blue)
- Add theme variants (gradient, minimal, custom)
- Update all styles to use theme system

**Test:** Visual verification in different terminals

---

## Task 2: Syntax Highlighting

**Summary:** Add Chroma-based syntax highlighting for code blocks in AI responses.

**Files:** `internal/tui/highlight.go`, modify `internal/tui/view.go`

**Dependencies:** `go get github.com/alecthomas/chroma/v2@latest`

**Key Steps:**
- Parse code blocks from AI responses
- Apply Chroma highlighting
- Render with Lipgloss styles
- Support multiple languages (Go, Python, JS, etc.)

**Test:** Send message asking for code, verify highlighting

---

## Task 3: Diff Preview Component

**Summary:** Create interactive diff viewer for code changes.

**Files:** `internal/tui/components/diff.go`

**Key Steps:**
- Unified diff format parser
- Color-coded diff display (red -, green +, gray context)
- Interactive approval (y/n/d/e/a)
- Keyboard navigation

**Test:** Mock diff data, verify rendering and interaction

---

## Task 4: Toast Notifications

**Summary:** Add non-intrusive notifications for events.

**Files:** `internal/tui/components/notification.go`

**Key Steps:**
- Toast message queue
- Auto-dismiss after 3 seconds
- Position: top-right corner
- Types: success, error, warning, info

**Test:** Trigger various notification types

---

## Task 5: Loading Animations

**Summary:** Add spinner and progress indicators.

**Files:** `internal/tui/components/spinner.go`, `internal/tui/animations.go`

**Key Steps:**
- Custom spinner frames (⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏)
- Progress bar for file analysis
- Fade-in effect for new messages
- Smooth transitions

**Test:** Verify animations in different terminal speeds

---

## Completion Notes

**Total Implementation Time:** 4-6 hours

**Priority:** Medium (enhances UX but not critical for MVP)

**Testing:** Primarily visual/manual testing

**Next:** Plan F (Token Optimizer) - THE core differentiator
