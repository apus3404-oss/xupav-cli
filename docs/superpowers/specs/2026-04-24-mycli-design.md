# Budget-Friendly AI Agent CLI - Design Specification

**Date:** 2026-04-24  
**Version:** 1.0.0  
**Status:** Approved

## Executive Summary

A terminal-based AI coding assistant that prioritizes cost efficiency and performance. Built with Go (CLI/TUI) and Python (AI integration), it offers smart token management to reduce API costs by 60-70% compared to traditional tools. Targets developers seeking Claude Code/Cursor alternatives with lower operational costs.

## Goals

**Primary:**
- Provide chat-based AI coding assistance with modern, colorful TUI
- Reduce token consumption through intelligent context management
- Support budget-friendly models (DeepSeek-R1, StepFun) via OpenRouter
- Deliver fast response times (<2s startup, <3s first token)

**Secondary:**
- Local AI support via Ollama integration
- LSP-powered codebase intelligence
- Multiple code modification modes (safe/interactive/auto)
- Cross-platform support (Linux, macOS, Windows)

## Non-Goals (MVP)

- Full IDE integration (VS Code extension, etc.)
- Multi-user collaboration features
- Cloud sync of conversation history
- Custom model fine-tuning
- Git operations beyond basic status/staging
- Plugin marketplace (future: v2.0)

## Target Users

1. **Solo developers** - Individual programmers working on personal projects
2. **Budget-conscious teams** - Small teams avoiding expensive subscriptions
3. **Students and hobbyists** - Learners seeking affordable AI tools
4. **Open source contributors** - Community members needing occasional AI help

## Success Metrics

- **Cost:** Average session cost <$0.10 (vs $0.30+ for competitors)
- **Performance:** 90th percentile response time <5s
- **Adoption:** 1K+ GitHub stars in first 3 months
- **Retention:** 60%+ weekly active users return after first week
## Architecture

### High-Level Design

**Approach:** Monolithic Go binary with embedded Python runtime for AI operations.

**Rationale:**
- Go provides fast startup, low memory footprint, excellent TUI libraries
- Python offers rich AI ecosystem (OpenRouter SDK, Ollama client)
- Single binary distribution simplifies user experience
- Subprocess communication keeps components isolated

**Architecture Diagram:**

```
┌─────────────────────────────────────┐
│   mycli (Go Binary)                 │
│                                     │
│  ┌──────────────┐  ┌─────────────┐ │
│  │   CLI/TUI    │  │   Config    │ │
│  │   (Bubble    │  │   Manager   │ │
│  │    Tea)      │  │             │ │
│  └──────┬───────┘  └──────┬──────┘ │
│         │                 │         │
│  ┌──────▼─────────────▼──────────┐ │
│  │   Core Engine                 │ │
│  │   - File Manager              │ │
│  │   - LSP Client                │ │
│  │   - Token Optimizer           │ │
│  │   - Context Builder           │ │
│  └──────┬────────────────────────┘ │
│         │                           │
│  ┌──────▼────────────────────────┐ │
│  │   Python Bridge (subprocess)  │ │
│  └───────────────────────────────┘ │
└─────────────┬───────────────────────┘
              │
    ┌─────────▼──────────┐
    │  Python AI Layer   │
    │  - OpenRouter SDK  │
    │  - Ollama Client   │
    │  - Response Parser │
    └────────────────────┘
```

### Technology Stack

**Go Components:**
- **Bubble Tea** - TUI framework for modern terminal UI
- **Cobra** - CLI command structure and flag parsing
- **go-lsp** - Language Server Protocol client
- **go-git** - Git operations (status, staging)
- **keyring** - Secure credential storage (OS keychain)

**Python Components:**
- **openai SDK** - OpenRouter API client (compatible endpoint)
- **ollama-python** - Local Ollama integration
- **tiktoken** - Token counting for budget tracking
- **tree-sitter** - AST parsing for code chunking

**Communication:**
- JSON-RPC over stdin/stdout between Go and Python
- Streaming support for real-time AI responses
- Timeout handling (60s default)

### Component Responsibilities

**CLI/TUI Layer (Go):**
- User input handling and command parsing
- Terminal rendering with colors, animations
- Slash command processing
- Session state management

**Core Engine (Go):**
- File system operations with .gitignore awareness
- LSP integration for code intelligence
- Token optimization algorithms
- Context assembly from multiple sources

**Python Bridge (Go):**
- Subprocess lifecycle management
- JSON-RPC protocol implementation
- Error handling and retry logic
- Fallback orchestration (OpenRouter → Ollama)

**AI Layer (Python):**
- API requests to OpenRouter/Ollama
- Streaming response handling
- Code block extraction and parsing
- Cost calculation and tracking

### Data Flow

**User Message → AI Response:**

1. User types message in TUI
2. Go parses input, checks for slash commands
3. Core Engine builds context (files, LSP data, history)
4. Token Optimizer reduces context size
5. Python Bridge sends JSON-RPC request
6. Python makes API call (streaming)
7. Response chunks flow back through bridge
8. TUI renders markdown with syntax highlighting
9. If code changes detected, show diff preview
10. User approves → File Writer applies changes

**Latency Budget:**
- Input processing: <50ms
- Context building: <500ms
- LSP queries: <200ms per query
- API first token: <3s (network dependent)
- Rendering: <100ms per chunk
## Smart Token Management

### Overview

The core differentiator: reduce token consumption by 60-70% through intelligent context optimization without sacrificing code understanding quality.

### Three-Layer Optimization

#### Layer 1: File Chunking

**Mechanism:**
- Parse files into AST using tree-sitter
- Extract semantic units (functions, classes, methods)
- Score each unit based on user query relevance
- Send only high-scoring units, not entire files

**Example:**
```
User: "fix the login bug"
File: auth.py (500 lines)

Traditional: Send all 500 lines (2000 tokens)
Smart: Send only login() function + dependencies (300 tokens)
Savings: 85%
```

**Implementation:**
- tree-sitter grammars for Go, Python, JS, TS, Rust
- Keyword matching: extract terms from user query
- Dependency tracking: include called functions
- Threshold: score >0.3 to include

#### Layer 2: Context Prioritization

**Mechanism:**
- Build dependency graph via LSP
- Score files by relevance to query
- Apply cutoff to exclude low-priority files
- Preserve critical context (imports, types)

**Scoring Algorithm:**
```
score = (keyword_matches * 0.4) + 
        (import_relevance * 0.3) + 
        (recent_edits * 0.2) + 
        (user_mentioned * 0.1)
```

**Example:**
```
User: "optimize database queries"
Codebase: 50 files

Scores:
- db.py: 0.95 ✓ include
- models.py: 0.78 ✓ include
- queries.py: 0.92 ✓ include
- utils.py: 0.45 ✓ include
- config.py: 0.12 ✗ exclude
- tests.py: 0.08 ✗ exclude
...

Result: 4 files sent instead of 50
Savings: 92%
```

#### Layer 3: Diff-Based Transmission

**Mechanism:**
- Track file hashes in conversation state
- If file previously sent, compute diff
- Send only changed lines with context
- Use unified diff format

**Example:**
```
First mention of file.py:
  Send: entire file (1000 tokens)

Second mention after edit:
  Send: diff only (50 tokens)
  
Savings: 95% on subsequent mentions
```

**Implementation:**
- SHA-256 hash per file
- Store in conversation state
- git diff algorithm for line-level changes
- ±3 lines context around changes

### Token Budget System

**Per-Request Budget:**
```
Model max tokens: 64K (DeepSeek-R1)
Output reserve: 4K (20%)
Available for context: 60K (80%)

Priority allocation:
1. System prompt: 2K
2. User message: 1K
3. Recent history (last 5 turns): 10K
4. Current files: 30K
5. Related files: 15K
6. Git status: 2K
```

**Budget Enforcement:**
- Calculate token count before sending (tiktoken)
- If over budget, apply aggressive filtering
- Remove lowest-scoring files first
- Warn user if critical context dropped

**Cost Tracking:**
```
Per message:
- Input tokens × model rate
- Output tokens × model rate
- Display running total in TUI header

Session summary:
- Total tokens used
- Total cost
- Average cost per message
- Savings vs. naive approach
```

### Optimization Trade-offs

**Accuracy vs. Cost:**
- More context = better answers, higher cost
- Less context = cheaper, risk missing info
- User configurable: `token_budget` in config
- Modes: `minimal` (30K), `balanced` (60K), `generous` (100K)

**Latency vs. Optimization:**
- AST parsing adds ~100ms per file
- LSP queries add ~200ms per query
- Acceptable for 60-70% cost savings
- Cache parsed ASTs to amortize cost

### Why This Matters

**Cost Comparison (1000-line codebase, 10-turn conversation):**

Traditional approach:
- Send all files every turn: 50K tokens/turn
- 10 turns × 50K = 500K tokens
- DeepSeek-R1: $0.14/M tokens
- Cost: $0.07

Smart approach:
- Optimized context: 15K tokens/turn
- 10 turns × 15K = 150K tokens
- Cost: $0.021
- **Savings: 70%**

**At scale (100 sessions/month):**
- Traditional: $7/month
- Smart: $2.10/month
- **Savings: $4.90/month per user**
## TUI/UX Design

### Design Philosophy

**Modern and Colorful:** Break away from plain terminal aesthetics. Use gradients, animations, and rich colors while maintaining readability and performance.

**Conversational:** Chat-based interaction feels natural. Users think in dialogue, not commands.

**Transparent:** Show what's happening (analyzing files, calling API, applying changes). No black boxes.

### Main Interface Layout

```
╔══════════════════════════════════════════════════════════╗
║  🤖 mycli v1.0.0          Model: deepseek-r1  💰 $0.02  ║
╠══════════════════════════════════════════════════════════╣
║                                                          ║
║  You: How can I optimize this database query?           ║
║                                                          ║
║  🤖 Assistant: Let me analyze your code...               ║
║     [Analyzing 3 files: db.py, models.py, queries.py]   ║
║                                                          ║
║     I found the issue in queries.py:45                   ║
║     You're using N+1 queries. Here's the fix:            ║
║                                                          ║
║     ┌─────────────────────────────────────────┐         ║
║     │ # queries.py:45                         │         ║
║     │ - for user in users:                    │         ║
║     │ -     posts = get_posts(user.id)        │         ║
║     │ + posts = get_posts_bulk([u.id for...]) │         ║
║     └─────────────────────────────────────────┘         ║
║                                                          ║
║     Apply changes? [y/n/d=diff]                         ║
║                                                          ║
╠══════════════════════════════════════════════════════════╣
║  📝 Type your message... (Ctrl+S to send, Ctrl+C exit)  ║
╚══════════════════════════════════════════════════════════╝
```

**Sections:**
1. **Header** - Status bar with version, model, cost
2. **Chat Area** - Scrollable conversation history
3. **Input Area** - Multi-line text input with hints

### Visual Elements

#### Color Scheme

**Header Gradient:**
- Cyan (#00D9FF) → Blue (#0080FF)
- Terminal capability detection: fallback to solid cyan if no gradient support

**Message Colors:**
- User messages: White text, subtle gray background
- AI messages: Cyan accent for username, white text
- System messages: Yellow for warnings, red for errors

**Code Blocks:**
- Background: Dark gray (#1E1E1E)
- Syntax highlighting via Chroma library
- Diff colors: Red (-), Green (+), Gray (context)

**Status Indicators:**
- 🟢 Green: Online, API connected
- 🟡 Yellow: Slow response, rate limited
- 🔴 Red: Offline, error state
- 🔵 Blue: Using local Ollama

#### Animations

**Spinner (AI thinking):**
```
⠋ ⠙ ⠹ ⠸ ⠼ ⠴ ⠦ ⠧ ⠇ ⠏
```
- 10 frames, 80ms per frame
- Shown during API calls
- Accompanied by status text: "Analyzing files...", "Waiting for response..."

**Progress Bar (file analysis):**
```
[████████░░░░░░░░░░░░] 40% (4/10 files)
```
- Updates in real-time as files are processed
- Shows percentage and count

**Fade-in Effect:**
- New messages fade in over 200ms
- Smooth transition, not jarring

**Shake Animation (errors):**
- Input box shakes horizontally on error
- 3 shakes, 50ms each
- Visual feedback for invalid input

#### Typography

**Fonts:**
- Monospace for code blocks (respects terminal font)
- Regular for prose
- Bold for emphasis (file names, commands)

**Sizing:**
- Normal text: default terminal size
- Headers: same size but bold + colored
- Code: same size, different background

### Input Experience

**Multi-line Support:**
- Enter: send message (default)
- Shift+Enter: new line
- Configurable: swap behavior if user prefers

**Auto-complete:**
- File paths: Tab to complete from project files
- Slash commands: Tab to complete from available commands
- History: Up/Down arrows to navigate previous messages

**Character Counter:**
- Bottom-right of input area
- Shows estimated tokens: "~450 tokens"
- Turns yellow when approaching budget limit
- Turns red when over budget

**Placeholder Text:**
```
When empty: "Ask me anything about your code..."
When typing: (disappears)
After error: "Try again or type /help for assistance"
```

### Slash Commands

**Available Commands:**
```
/help          Show command list and keyboard shortcuts
/clear         Clear conversation history
/files         List files in current context
/add <path>    Add file to context
/remove <path> Remove file from context
/model <name>  Switch AI model
/cost          Show session cost breakdown
/export        Export conversation to markdown
/reset         Start fresh conversation
/config        Open config in editor
/debug         Toggle debug mode
```

**Command Autocomplete:**
- Type `/` to see command list
- Type `/he` + Tab → `/help`
- Fuzzy matching: `/fi` matches `/files`

### Notifications

**Toast Notifications:**
- Position: Top-right corner
- Size: 40 chars wide, auto-height
- Duration: 3 seconds, then fade out
- Types:
  - Success: ✓ Green background
  - Error: ✗ Red background
  - Warning: ⚠ Yellow background
  - Info: ℹ Blue background

**Examples:**
```
✓ Changes applied to db.py
✗ API key invalid
⚠ Token budget exceeded
ℹ Switched to deepseek-r1
```

### Diff Preview

**Interactive Diff View:**
```
┌─────────────────────────────────────────────────────────┐
│ File: queries.py                                        │
├─────────────────────────────────────────────────────────┤
│  42 │ def get_user_posts(user_ids):                     │
│  43 │ -   results = []                                  │
│  44 │ -   for uid in user_ids:                          │
│  45 │ -       posts = db.query(Post).filter_by(         │
│  46 │ -           user_id=uid                           │
│  47 │ -       ).all()                                   │
│  48 │ -       results.extend(posts)                     │
│  49 │ +   results = db.query(Post).filter(              │
│  50 │ +       Post.user_id.in_(user_ids)                │
│  51 │ +   ).all()                                       │
│  52 │     return results                                │
└─────────────────────────────────────────────────────────┘

Apply this change? [y]es [n]o [d]etailed [e]dit [a]ll
```

**Options:**
- `y` - Apply this change
- `n` - Skip this change
- `d` - Show full file context (±10 lines)
- `e` - Open in $EDITOR for manual edit
- `a` - Apply all remaining changes (this session)

### Responsive Design

**Terminal Size Adaptation:**
- Minimum: 80×24 (standard)
- Optimal: 120×40
- Maximum: No limit, scales gracefully

**Small Terminal (<80 cols):**
- Disable animations
- Simplify layout (no borders)
- Truncate long lines with "..."

**Large Terminal (>150 cols):**
- Side-by-side diff view
- Wider code blocks
- More context lines in diffs

### Accessibility

**Color Blindness:**
- Don't rely solely on color for information
- Use symbols: ✓✗⚠ℹ
- High contrast mode available

**Screen Readers:**
- Semantic markup where possible
- Status updates announced
- Alternative text for symbols

**Keyboard-Only:**
- All features accessible via keyboard
- No mouse required
- Tab navigation through interactive elements
## Configuration and Setup

### Interactive Setup Flow

**First Run Experience (`mycli chat` without config):**

```
$ mycli chat

🤖 Welcome to mycli! Let's set you up.

? Select your primary AI provider:
  ❯ OpenRouter (recommended for budget)
    Ollama (local, free)
    I'll configure later

? Enter your OpenRouter API key:
  (paste here, will be hidden) ••••••••••••

? Select default model:
  ❯ deepseek/deepseek-r1 ($0.14/M tokens) ⭐ Recommended
    stepfun/step-2-16k ($0.10/M tokens)
    anthropic/claude-sonnet-4 ($3.00/M tokens)
    Custom model...

? Code modification mode:
  ❯ Interactive (show diff, ask before applying)
    Safe (only suggest, never auto-apply)
    Auto (apply immediately, trust me)

? Enable LSP integration? (smarter context)
  ❯ Yes (recommended)
    No

✓ Configuration saved to ~/.mycli/config.yaml
✓ API key stored in system keychain
✓ Ready to go! Type your first message...
```

**Design Principles:**
- Sensible defaults (OpenRouter + DeepSeek-R1 + Interactive mode)
- Progressive disclosure (advanced options via `mycli config`)
- Non-blocking (can skip and configure later)
- Secure by default (keychain storage)

### Configuration File Structure

**Location:** `~/.mycli/config.yaml`

```yaml
version: "1.0"

providers:
  openrouter:
    enabled: true
    default_model: "deepseek/deepseek-r1"
    api_key_source: "keychain"  # never plain text
    base_url: "https://openrouter.ai/api/v1"
    max_tokens: 4096
    temperature: 0.7
    timeout: 60000  # ms
    
  ollama:
    enabled: true
    base_url: "http://localhost:11434"
    default_model: "codellama:13b"
    fallback: true  # use if OpenRouter fails
    timeout: 120000  # local models can be slower

  # Future: anthropic, openai
  # anthropic:
  #   enabled: false
  #   api_key_source: "keychain"

behavior:
  code_mode: "interactive"  # safe | interactive | auto
  auto_add_files: true      # auto-detect relevant files via LSP
  max_context_files: 10     # limit to prevent token explosion
  token_budget: 60000       # max tokens for context (80% of model limit)
  conversation_memory: 10   # number of turns to remember

ui:
  theme: "gradient"         # gradient | minimal | custom
  animations: true          # disable for slow terminals
  syntax_highlight: true
  show_token_count: true
  show_cost: true
  color_scheme: "auto"      # auto | dark | light

lsp:
  enabled: true
  languages: ["go", "python", "javascript", "typescript", "rust"]
  timeout: 5000             # ms per LSP query
  max_queries: 5            # limit concurrent queries
  
cache:
  enabled: true
  ttl: 3600                 # seconds (1 hour)
  max_size: "100MB"
  location: "~/.mycli/cache"

logging:
  level: "info"             # debug | info | warn | error
  file: "~/.mycli/logs/mycli.log"
  max_size: "10MB"
  max_backups: 3
```

**Validation:**
- Schema validation on load
- Type checking for all fields
- Range validation (e.g., token_budget must be >1000)
- Warn on deprecated fields

### Environment Variable Overrides

**Priority:** ENV VAR > config.yaml > defaults

**Supported Variables:**
```bash
# API Keys (most important)
MYCLI_OPENROUTER_KEY=sk-or-...
MYCLI_ANTHROPIC_KEY=sk-ant-...
MYCLI_OPENAI_KEY=sk-...

# Provider Selection
MYCLI_PROVIDER=openrouter  # openrouter | ollama | anthropic
MYCLI_MODEL=deepseek/deepseek-r1

# Behavior
MYCLI_CODE_MODE=auto       # safe | interactive | auto
MYCLI_TOKEN_BUDGET=80000

# UI
MYCLI_THEME=minimal
MYCLI_NO_ANIMATIONS=true   # disable animations
MYCLI_NO_COLOR=true        # disable colors (CI/CD)

# Debug
MYCLI_DEBUG=true
MYCLI_LOG_LEVEL=debug
MYCLI_TRACE_API=true       # log all API requests/responses
```

**Use Cases:**
- CI/CD: `MYCLI_NO_COLOR=true MYCLI_CODE_MODE=safe`
- Testing: `MYCLI_PROVIDER=ollama MYCLI_DEBUG=true`
- Temporary override: `MYCLI_MODEL=claude-sonnet-4 mycli chat`

### Code Modification Modes

#### Safe Mode
**Behavior:**
- AI generates code suggestions
- Displays in code blocks
- User manually copies and applies
- No file system writes

**When to use:**
- Untrusted codebases
- Learning/exploration
- Reviewing AI suggestions before applying

**Trade-off:** Safest but slowest workflow

#### Interactive Mode (Default)
**Behavior:**
- AI generates code changes
- Shows unified diff preview
- Prompts for approval: `[y/n/d/e/a]`
- Applies only on explicit approval
- Creates backup before applying

**When to use:**
- Normal development workflow
- Want to review changes before applying
- Balance between safety and speed

**Trade-off:** Balanced approach

#### Auto Mode
**Behavior:**
- AI directly modifies files
- Shows notification of changes
- Automatically stages changes in git
- Provides undo command: `mycli undo`

**When to use:**
- Trusted AI models
- Rapid prototyping
- Experienced users who trust the tool

**Trade-off:** Fastest but requires trust

**Undo Mechanism:**
```bash
$ mycli undo
? Undo last change to db.py? (y/n) y
✓ Restored db.py from backup
✓ Unstaged from git
```

### API Key Management

**Security Requirements:**
- Never store keys in plain text
- Use OS-native secure storage
- Support key rotation
- Warn on key exposure (e.g., in logs)

**Storage Strategy:**

**macOS:** Keychain Access
```go
import "github.com/zalando/go-keyring"

keyring.Set("mycli", "openrouter", apiKey)
key, _ := keyring.Get("mycli", "openrouter")
```

**Windows:** Credential Manager
```go
// Same API via go-keyring
keyring.Set("mycli", "openrouter", apiKey)
```

**Linux:** Secret Service API
```go
// Supports gnome-keyring, kwallet
keyring.Set("mycli", "openrouter", apiKey)
```

**Fallback (no keychain):**
- Encrypt key with AES-256
- Derive encryption key from machine ID + user ID
- Store in `~/.mycli/secrets.enc`
- Warn user: "Keychain not available, using encrypted storage"

**Key Rotation:**
```bash
$ mycli config set-key openrouter
? Enter new OpenRouter API key: ••••••••••••
✓ Key updated and stored securely
✓ Old key removed from keychain
```

### Configuration Commands

**View current config:**
```bash
$ mycli config show
Provider: openrouter
Model: deepseek/deepseek-r1
Code Mode: interactive
LSP: enabled
Token Budget: 60000
```

**Edit config file:**
```bash
$ mycli config edit
# Opens ~/.mycli/config.yaml in $EDITOR
```

**Set individual values:**
```bash
$ mycli config set model stepfun/step-2-16k
✓ Model updated to stepfun/step-2-16k

$ mycli config set code_mode auto
✓ Code mode updated to auto
```

**Reset to defaults:**
```bash
$ mycli config reset
? This will reset all settings to defaults. Continue? (y/n) y
✓ Configuration reset
✓ API keys preserved
```

### Project-Specific Configuration

**Override global config per project:**

**Location:** `.mycli/config.yaml` (in project root)

**Use Case:**
- Different model for different projects
- Stricter token budget for large codebases
- Disable LSP for specific languages

**Example:**
```yaml
# .mycli/config.yaml (project-specific)
providers:
  openrouter:
    default_model: "anthropic/claude-sonnet-4"  # use better model for this project
    
behavior:
  token_budget: 40000  # stricter budget for large codebase
  max_context_files: 5
  
lsp:
  languages: ["go"]  # only Go LSP for this project
```

**Merge Strategy:**
- Project config overrides global config
- Unspecified fields inherit from global
- Explicit `null` disables global setting
## Data Flow and Error Handling

### End-to-End Data Flow

**Complete Request Lifecycle:**

```
1. User Input (TUI)
   ↓
2. Message Parser (Go)
   - Slash command? → Command Handler
   - Normal message? → Context Builder
   ↓
3. Context Builder (Go)
   - Load conversation history (last 10 turns)
   - Query LSP for relevant files
   - Read file contents
   - Extract git status
   ↓
4. Token Optimizer (Go)
   - Parse files to AST
   - Score and filter chunks
   - Apply diff-based optimization
   - Enforce token budget
   ↓
5. JSON Payload Assembly (Go)
   {
     "method": "chat",
     "params": {
       "message": "user message",
       "context": {...},
       "model": "deepseek/deepseek-r1",
       "max_tokens": 4096
     }
   }
   ↓
6. Python Bridge (Go → Python)
   - Start Python subprocess if not running
   - Send JSON-RPC request via stdin
   - Set 60s timeout
   ↓
7. AI Provider (Python)
   - Route to OpenRouter or Ollama
   - Make streaming API request
   - Handle rate limits, retries
   ↓
8. Response Streaming (Python → Go)
   - Chunk-by-chunk transmission
   - Each chunk: JSON-RPC notification
   - Final chunk: completion marker
   ↓
9. Response Parser (Go)
   - Accumulate chunks
   - Extract code blocks
   - Detect file modifications
   - Calculate token usage and cost
   ↓
10. TUI Renderer (Go)
    - Markdown rendering
    - Syntax highlighting
    - Animate text appearance
    ↓
11. Code Change Detection (Go)
    - If modifications detected → Diff Preview
    - Show unified diff
    - Prompt for approval
    ↓
12. User Approval (if interactive mode)
    - y → Apply changes
    - n → Skip
    - d → Show detailed diff
    - e → Open in editor
    ↓
13. File Writer (Go)
    - Create backup: ~/.mycli/backups/<timestamp>/
    - Apply changes
    - Stage in git (if in repo)
    - Show success notification
```

**Timing Breakdown (Target):**
- Steps 1-5: <500ms (local processing)
- Step 6: <100ms (subprocess communication)
- Step 7: 2-5s (API dependent)
- Steps 8-13: <200ms (local processing)
- **Total: 3-6s** (dominated by API latency)

### Error Handling Strategy

#### Layer 1: Python Bridge Errors

**Python Runtime Not Found:**
```
Error: Python 3.8+ not found
Solution: Install Python from python.org

Would you like to:
1. Open installation guide
2. Specify custom Python path
3. Exit
```

**Python Subprocess Crash:**
```
Action: Automatic restart (up to 3 times)
Backoff: 1s, 2s, 4s
After 3 failures: Fallback to Ollama or error message
```

**JSON-RPC Protocol Error:**
```
Log: Full request/response to debug log
User: "Internal communication error. Check logs: ~/.mycli/logs/mycli.log"
Recovery: Restart Python subprocess
```

#### Layer 2: API Errors

**Rate Limit (HTTP 429):**
```
Detection: "rate_limit_exceeded" in response
Action: Exponential backoff
Backoff: 1s, 2s, 4s, 8s, 16s (max 5 retries)
UI: "⏳ Rate limited. Retrying in 4s..."
```

**Invalid API Key (HTTP 401):**
```
Detection: "invalid_api_key" or 401 status
Action: Stop immediately (no retry)
UI: "❌ API key invalid. Run: mycli config set-key openrouter"
Recovery: User must update key
```

**Model Not Found (HTTP 404):**
```
Detection: "model_not_found" in response
Action: Fetch available models from API
UI: "❌ Model 'xyz' not found. Available models:
     - deepseek/deepseek-r1
     - stepfun/step-2-16k
     - ..."
Recovery: User selects valid model
```

**Request Timeout (>60s):**
```
Detection: No response within timeout
Action: Cancel request
UI: "⏱️ Request timed out. Model may be overloaded.
     Try again or switch to faster model?"
Options: [Retry] [Change Model] [Cancel]
```

**Network Error:**
```
Detection: Connection refused, DNS failure, etc.
Action: Retry 3 times with 2s delay
Fallback: Switch to Ollama if available
UI: "🔌 Network error. Trying Ollama..."
```

**Token Limit Exceeded:**
```
Detection: "context_length_exceeded" in response
Action: Aggressive context pruning
- Remove lowest-priority files
- Reduce conversation history to 5 turns
- Retry request
UI: "⚠️ Context too large. Reducing and retrying..."
```

#### Layer 3: File System Errors

**File Not Found:**
```
Detection: os.IsNotExist(err)
Action: Fuzzy search for similar files
UI: "❌ File not found: db.py
     Did you mean:
     - database.py
     - models/db.py"
Recovery: User corrects path
```

**Permission Denied:**
```
Detection: os.IsPermission(err)
Action: Check file permissions
UI: "❌ Cannot write to db.py (permission denied)
     Current permissions: -r--r--r--
     Run: chmod u+w db.py"
Recovery: User fixes permissions
```

**Disk Full:**
```
Detection: "no space left on device"
Action: Stop immediately
UI: "❌ Disk full. Free up space and try again.
     Current usage: 98% (45GB / 46GB)"
Recovery: User frees disk space
```

**Git Conflict:**
```
Detection: Uncommitted changes in target file
Action: Warn before overwriting
UI: "⚠️ db.py has uncommitted changes.
     Applying AI changes will overwrite them.
     Continue? [y/n]"
Recovery: User commits or stashes first
```

#### Layer 4: LSP Errors

**LSP Server Crash:**
```
Detection: LSP process exit
Action: Disable LSP for this session
UI: "⚠️ LSP server crashed. Continuing without code intelligence."
Fallback: Simple file scanning
Log: Full LSP logs to debug file
```

**LSP Timeout:**
```
Detection: No response within 5s
Action: Cancel query, continue without result
UI: (No user-facing message, silent degradation)
Fallback: Skip this LSP query
```

**LSP Not Available for Language:**
```
Detection: No LSP server configured
Action: Skip LSP, use file scanning
UI: (No message, expected behavior)
```

### Graceful Degradation

**Feature Priority Hierarchy:**

1. **Critical (must work):**
   - Chat functionality
   - Message display
   - Basic file reading

2. **Important (degrade gracefully):**
   - File writing (fallback to suggestions)
   - Syntax highlighting (fallback to plain text)
   - Cost tracking (fallback to estimates)

3. **Nice-to-have (disable if broken):**
   - LSP integration
   - Animations
   - Token optimization (fallback to simple truncation)

**Degradation Scenarios:**

**No LSP:**
```
Impact: Less accurate file selection
Fallback: User manually specifies files with /add
Performance: Slightly higher token usage
```

**No OpenRouter:**
```
Impact: Cannot use cloud models
Fallback: Switch to Ollama automatically
UI: "🔵 Using local Ollama (OpenRouter unavailable)"
```

**No Ollama:**
```
Impact: No local fallback
Fallback: Suggestion-only mode (no API calls)
UI: "⚠️ No AI provider available. Running in suggestion mode."
```

**No Color Support:**
```
Impact: Plain terminal output
Fallback: Disable colors, use ASCII symbols
UI: Plain text, no gradients
```

**Slow Terminal:**
```
Detection: Render time >100ms per frame
Action: Disable animations automatically
UI: Static output, no spinners
```

### Logging and Debugging

**Log Levels:**

**ERROR:** Critical failures requiring user action
```
[ERROR] API key invalid (401)
[ERROR] Cannot write to file: permission denied
[ERROR] Python subprocess crashed after 3 retries
```

**WARN:** Degraded functionality, but continuing
```
[WARN] LSP timeout after 5s, skipping query
[WARN] Rate limited, retrying in 4s
[WARN] Token budget exceeded, pruning context
```

**INFO:** Normal operations
```
[INFO] Model changed to deepseek-r1
[INFO] File added to context: db.py
[INFO] Applied changes to 3 files
```

**DEBUG:** Detailed information for troubleshooting
```
[DEBUG] Token count: 15234 (budget: 60000)
[DEBUG] LSP query: textDocument/definition
[DEBUG] API request: POST /chat/completions (2.3s)
```

**Log File Format (JSON):**
```json
{
  "timestamp": "2026-04-24T15:30:45Z",
  "level": "ERROR",
  "component": "python_bridge",
  "message": "Subprocess crashed",
  "context": {
    "exit_code": 1,
    "stderr": "ModuleNotFoundError: openai",
    "attempt": 3
  }
}
```

**Log Rotation:**
- Max size: 10MB per file
- Keep last 3 files
- Compress old logs (gzip)
- Location: `~/.mycli/logs/`

**Debug Mode:**
```bash
$ MYCLI_DEBUG=true mycli chat

[DEBUG] Starting Python subprocess: python3 -m mycli_ai.server
[DEBUG] Subprocess PID: 12345
[DEBUG] Building context: 5 files, 12000 tokens
[DEBUG] API request: {model: "deepseek-r1", max_tokens: 4096}
[DEBUG] API response: 2.3s, 450 tokens, $0.0063
```

**Trace Mode (API requests):**
```bash
$ MYCLI_TRACE_API=true mycli chat

→ POST https://openrouter.ai/api/v1/chat/completions
  Headers: {Authorization: "Bearer sk-or-...", ...}
  Body: {model: "deepseek/deepseek-r1", messages: [...]}

← 200 OK (2.3s)
  Body: {choices: [{message: {content: "..."}}], usage: {...}}
```
## Installation and Distribution

### Installation Strategy

**Primary Method:** One-line installer script (curl | bash)

**Rationale:**
- Simplest user experience
- Handles platform detection automatically
- Sets up Python dependencies
- Configures PATH
- No package manager required

### Installer Script (`install.sh`)

**Usage:**
```bash
curl -fsSL https://raw.githubusercontent.com/yourusername/mycli/main/install.sh | bash
```

**Script Responsibilities:**

#### 1. Platform Detection
```bash
OS=$(uname -s)      # Linux, Darwin, MINGW64_NT (Windows Git Bash)
ARCH=$(uname -m)    # x86_64, arm64, aarch64

# Normalize
case "$OS" in
  Linux*)   OS="linux" ;;
  Darwin*)  OS="darwin" ;;
  MINGW*)   OS="windows" ;;
esac

case "$ARCH" in
  x86_64)   ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
esac
```

#### 2. Python Version Check
```bash
# Require Python 3.8+
PYTHON_CMD=$(command -v python3 || command -v python)
PYTHON_VERSION=$($PYTHON_CMD --version 2>&1 | awk '{print $2}')

if ! version_gte "$PYTHON_VERSION" "3.8.0"; then
  echo "❌ Python 3.8+ required (found: $PYTHON_VERSION)"
  echo "Install from: https://python.org"
  exit 1
fi
```

#### 3. Binary Download
```bash
VERSION="v1.0.0"  # or fetch latest from GitHub API
BINARY_URL="https://github.com/yourusername/mycli/releases/download/$VERSION/mycli-$OS-$ARCH"
CHECKSUM_URL="$BINARY_URL.sha256"

# Download binary
curl -fsSL "$BINARY_URL" -o /tmp/mycli

# Verify checksum
EXPECTED=$(curl -fsSL "$CHECKSUM_URL")
ACTUAL=$(sha256sum /tmp/mycli | awk '{print $1}')

if [ "$EXPECTED" != "$ACTUAL" ]; then
  echo "❌ Checksum mismatch. Download corrupted."
  exit 1
fi

echo "✓ Binary verified"
```

#### 4. Python Virtual Environment
```bash
# Create venv in ~/.mycli/venv
VENV_DIR="$HOME/.mycli/venv"
$PYTHON_CMD -m venv "$VENV_DIR"

# Activate and install dependencies
source "$VENV_DIR/bin/activate"
pip install --quiet --upgrade pip
pip install --quiet -r https://raw.githubusercontent.com/yourusername/mycli/main/python/requirements.txt

echo "✓ Python dependencies installed"
```

**requirements.txt (pinned versions):**
```
openai==1.12.0
ollama==0.1.7
tiktoken==0.6.0
tree-sitter==0.20.4
tree-sitter-python==0.20.4
tree-sitter-go==0.20.0
tree-sitter-javascript==0.20.3
tree-sitter-typescript==0.20.5
tree-sitter-rust==0.20.4
```

#### 5. Binary Installation
```bash
# Install location
if [ "$OS" = "windows" ]; then
  INSTALL_DIR="$HOME/bin"
else
  INSTALL_DIR="/usr/local/bin"
fi

# Create directory if needed (no sudo required for ~/bin)
mkdir -p "$INSTALL_DIR"

# Move binary
mv /tmp/mycli "$INSTALL_DIR/mycli"
chmod +x "$INSTALL_DIR/mycli"

echo "✓ Installed to $INSTALL_DIR/mycli"
```

#### 6. PATH Configuration
```bash
# Add to PATH if not already present
case "$SHELL" in
  */bash)  RC_FILE="$HOME/.bashrc" ;;
  */zsh)   RC_FILE="$HOME/.zshrc" ;;
  */fish)  RC_FILE="$HOME/.config/fish/config.fish" ;;
  *)       RC_FILE="$HOME/.profile" ;;
esac

if ! grep -q "$INSTALL_DIR" "$RC_FILE" 2>/dev/null; then
  echo "export PATH=\"$INSTALL_DIR:\$PATH\"" >> "$RC_FILE"
  echo "✓ Added to PATH in $RC_FILE"
  echo "⚠️  Run: source $RC_FILE (or restart terminal)"
fi
```

#### 7. Initial Configuration
```bash
# Run interactive setup
"$INSTALL_DIR/mycli" init

echo ""
echo "🎉 Installation complete!"
echo ""
echo "Get started:"
echo "  mycli chat          Start chatting"
echo "  mycli --help        Show all commands"
echo "  mycli config show   View configuration"
```

### Security Measures

**Checksum Verification:**
- SHA-256 hash for every binary
- Published alongside releases
- Verified before installation

**HTTPS Only:**
- All downloads over HTTPS
- GitHub's SSL certificates

**No Sudo Required:**
- Install to user directory (`~/bin` or `~/.local/bin`)
- No system-wide changes
- Safer for users

**Optional GPG Signature:**
```bash
# For paranoid users
curl -fsSL "$BINARY_URL.sig" -o /tmp/mycli.sig
gpg --verify /tmp/mycli.sig /tmp/mycli
```

### Alternative Installation Methods

#### Manual Installation (for developers)

```bash
# Clone repository
git clone https://github.com/yourusername/mycli.git
cd mycli

# Build Go binary
make build
# Output: ./bin/mycli

# Install Python dependencies
cd python
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt

# Run
./bin/mycli chat
```

#### From Source (Go developers)

```bash
# Requires Go 1.21+
go install github.com/yourusername/mycli/cmd/mycli@latest

# Still need Python dependencies
cd $(go env GOPATH)/src/github.com/yourusername/mycli/python
pip install -r requirements.txt
```

### Update Mechanism

**Automatic Update Check:**
- Check GitHub API every 24 hours
- Compare local version with latest release
- Show notification if update available
- Non-intrusive (doesn't block usage)

**Update Flow:**

```bash
$ mycli chat

ℹ️  New version available: v1.2.0 (current: v1.0.0)
   Run: mycli update

$ mycli update

🔍 Checking for updates...
✓ New version available: v1.2.0 (current: v1.0.0)

📝 Changelog:
  • Added Anthropic Claude support
  • Fixed LSP timeout issue
  • Improved token optimization by 15%
  • Security: Updated dependencies

? Update now? (y/n) y

⬇️  Downloading v1.2.0...
[████████████████████] 100% (5.2 MB / 5.2 MB)

✓ Download complete
✓ Checksum verified
✓ Backup created: ~/.mycli/backups/v1.0.0
✓ Binary updated
✓ Python dependencies updated

🎉 Successfully updated to v1.2.0!

Restart mycli to use the new version.
```

**Update Command Implementation:**
```bash
$ mycli update --check      # Check without installing
$ mycli update --force      # Skip confirmation
$ mycli update --version v1.1.0  # Install specific version
```

**Rollback Support:**
```bash
$ mycli update --rollback

Available backups:
  v1.1.0 (2026-04-20)
  v1.0.0 (2026-04-15)

? Select version to restore: v1.1.0

✓ Restored to v1.1.0
```

### Uninstallation

**Clean Removal:**

```bash
$ mycli uninstall

⚠️  This will remove:
  • Binary: /usr/local/bin/mycli
  • Configuration: ~/.mycli/config.yaml
  • Cache: ~/.mycli/cache/
  • Logs: ~/.mycli/logs/
  • Python venv: ~/.mycli/venv/
  • Backups: ~/.mycli/backups/

? Keep conversation history? (y/n) y

⏳ Removing mycli...
✓ Binary removed
✓ Cache cleared (45 MB freed)
✓ Logs removed
✓ Python venv removed
✓ Conversation history saved to: ~/mycli-backup-2026-04-24/

🎉 mycli uninstalled successfully

To reinstall: curl -fsSL https://mycli.dev/install.sh | bash
```

**Manual Uninstall:**
```bash
# Remove binary
rm /usr/local/bin/mycli

# Remove data directory
rm -rf ~/.mycli

# Remove PATH entry from shell RC file
# (manual edit of ~/.bashrc, ~/.zshrc, etc.)
```

### Platform-Specific Notes

#### macOS

**Gatekeeper Warning:**
```
"mycli" cannot be opened because it is from an unidentified developer.
```

**Solutions:**

1. **User workaround:**
```bash
xattr -d com.apple.quarantine /usr/local/bin/mycli
```

2. **Developer solution (future):**
- Sign binary with Apple Developer certificate
- Notarize with Apple
- Requires paid Apple Developer account ($99/year)

**Homebrew (future):**
```bash
brew tap yourusername/mycli
brew install mycli
```

#### Windows

**Windows Defender SmartScreen:**
```
Windows protected your PC
Microsoft Defender SmartScreen prevented an unrecognized app from starting.
```

**Solutions:**

1. **User workaround:**
- Click "More info"
- Click "Run anyway"

2. **Developer solution (future):**
- Code-sign binary with certificate
- Requires code signing certificate (~$100-300/year)

**Recommended Environment:**
- Git Bash (comes with Git for Windows)
- WSL2 (Windows Subsystem for Linux)
- Windows Terminal (better color support)

**Native CMD/PowerShell:**
- Limited color support
- Automatic fallback to minimal theme
- Full functionality preserved

#### Linux

**Distribution Support:**
- Ubuntu 20.04+ ✓
- Debian 11+ ✓
- Fedora 35+ ✓
- Arch Linux ✓
- Alpine Linux ✓ (musl libc)

**Python Availability:**
- Most distros include Python 3.8+
- Minimal distros may need: `apt install python3 python3-venv`

**Package Managers (future):**
```bash
# Debian/Ubuntu
apt install mycli

# Fedora
dnf install mycli

# Arch
yay -S mycli
```

### Distribution Checklist

**Pre-Release:**
- [ ] Build binaries for all platforms (linux-amd64, linux-arm64, darwin-amd64, darwin-arm64, windows-amd64)
- [ ] Generate SHA-256 checksums
- [ ] Test installer script on all platforms
- [ ] Verify Python dependencies install correctly
- [ ] Test update mechanism
- [ ] Write release notes

**Release:**
- [ ] Tag version in git: `git tag v1.0.0`
- [ ] Push to GitHub: `git push --tags`
- [ ] GitHub Actions builds and uploads binaries
- [ ] Create GitHub Release with changelog
- [ ] Update install.sh with new version
- [ ] Announce on social media, forums

**Post-Release:**
- [ ] Monitor error reports
- [ ] Track installation metrics
- [ ] Gather user feedback
- [ ] Plan next release
## Testing Strategy and Project Structure

### Project Directory Structure

```
mycli/
├── cmd/
│   └── mycli/
│       └── main.go                    # Entry point, CLI initialization
│
├── internal/
│   ├── cli/
│   │   ├── commands.go                # Cobra command definitions
│   │   ├── chat.go                    # Chat command implementation
│   │   ├── config.go                  # Config command implementation
│   │   ├── update.go                  # Update command implementation
│   │   └── flags.go                   # Global flags
│   │
│   ├── tui/
│   │   ├── app.go                     # Bubble Tea main app
│   │   ├── model.go                   # App state model
│   │   ├── update.go                  # Update logic
│   │   ├── view.go                    # Render logic
│   │   ├── components/
│   │   │   ├── chat.go                # Chat message component
│   │   │   ├── input.go               # Input box component
│   │   │   ├── diff.go                # Diff preview component
│   │   │   ├── header.go              # Status header component
│   │   │   └── notification.go        # Toast notification component
│   │   └── styles.go                  # Colors, themes, styling
│   │
│   ├── core/
│   │   ├── context.go                 # Context builder
│   │   ├── optimizer.go               # Token optimizer (3 layers)
│   │   ├── files.go                   # File manager (.gitignore aware)
│   │   ├── lsp/
│   │   │   ├── client.go              # LSP client interface
│   │   │   ├── go.go                  # Go LSP implementation
│   │   │   ├── python.go              # Python LSP implementation
│   │   │   └── javascript.go          # JS/TS LSP implementation
│   │   ├── git.go                     # Git operations
│   │   └── cache.go                   # File hash cache
│   │
│   ├── bridge/
│   │   ├── python.go                  # Python subprocess manager
│   │   ├── protocol.go                # JSON-RPC protocol
│   │   └── stream.go                  # Streaming response handler
│   │
│   ├── config/
│   │   ├── config.go                  # Config loading/saving
│   │   ├── keychain.go                # Secure API key storage
│   │   ├── defaults.go                # Default configuration
│   │   └── validation.go              # Config validation
│   │
│   └── writer/
│       ├── file.go                    # File writer with backup
│       ├── diff.go                    # Diff generation
│       └── undo.go                    # Undo mechanism
│
├── python/
│   ├── mycli_ai/
│   │   ├── __init__.py
│   │   ├── server.py                  # JSON-RPC server (stdin/stdout)
│   │   ├── providers/
│   │   │   ├── __init__.py
│   │   │   ├── base.py                # Provider interface
│   │   │   ├── openrouter.py          # OpenRouter implementation
│   │   │   └── ollama.py              # Ollama implementation
│   │   ├── parser.py                  # Response parser (code blocks, etc.)
│   │   ├── cost.py                    # Cost calculation
│   │   └── tokens.py                  # Token counting (tiktoken)
│   ├── requirements.txt               # Pinned dependencies
│   └── setup.py                       # Package metadata
│
├── scripts/
│   ├── install.sh                     # Installation script
│   ├── build.sh                       # Build script (all platforms)
│   └── release.sh                     # Release automation
│
├── tests/
│   ├── go/
│   │   ├── optimizer_test.go          # Token optimizer tests
│   │   ├── bridge_test.go             # Python bridge tests
│   │   ├── config_test.go             # Config tests
│   │   └── integration_test.go        # End-to-end tests
│   └── python/
│       ├── test_providers.py          # Provider tests
│       ├── test_parser.py             # Parser tests
│       └── test_server.py             # JSON-RPC server tests
│
├── docs/
│   ├── superpowers/
│   │   └── specs/
│   │       └── 2026-04-24-mycli-design.md  # This document
│   ├── architecture.md                # Architecture deep-dive
│   ├── contributing.md                # Contribution guide
│   └── api.md                         # JSON-RPC API spec
│
├── .github/
│   └── workflows/
│       ├── build.yml                  # Build on push
│       ├── test.yml                   # Run tests
│       └── release.yml                # Release automation
│
├── Makefile                           # Build commands
├── go.mod                             # Go dependencies
├── go.sum                             # Go dependency checksums
├── .gitignore
├── LICENSE                            # MIT or Apache 2.0
└── README.md                          # Project overview
```

### Testing Strategy (MVP Focus)

**Philosophy:** Minimal but effective testing for MVP. Focus on critical paths and error handling. Expand test coverage post-MVP based on real-world usage.

#### Unit Tests (Go)

**Token Optimizer:**
```go
// internal/core/optimizer_test.go

func TestOptimizer_FileChunking(t *testing.T) {
    code := `
    def login(username, password):
        return authenticate(username, password)
    
    def logout(session):
        return destroy_session(session)
    `
    
    chunks := optimizer.ChunkFile(code, "python")
    assert.Equal(t, 2, len(chunks))  // 2 functions
    assert.Contains(t, chunks[0].Code, "login")
    assert.Contains(t, chunks[1].Code, "logout")
}

func TestOptimizer_ContextPrioritization(t *testing.T) {
    files := []File{
        {Path: "db.py", Content: "database code"},
        {Path: "models.py", Content: "model definitions"},
        {Path: "utils.py", Content: "utility functions"},
    }
    
    query := "optimize database query"
    scored := optimizer.ScoreFiles(files, query)
    
    // db.py should score highest
    assert.Greater(t, scored[0].Score, scored[1].Score)
    assert.Equal(t, "db.py", scored[0].Path)
}

func TestOptimizer_DiffBased(t *testing.T) {
    original := "line1\nline2\nline3"
    modified := "line1\nline2_modified\nline3"
    
    diff := optimizer.GenerateDiff(original, modified)
    
    assert.Contains(t, diff, "-line2")
    assert.Contains(t, diff, "+line2_modified")
    assert.NotContains(t, diff, "line1")  // unchanged lines excluded
}
```

**Python Bridge:**
```go
// internal/bridge/python_test.go

func TestBridge_StartStop(t *testing.T) {
    bridge := NewPythonBridge()
    
    err := bridge.Start()
    assert.NoError(t, err)
    assert.True(t, bridge.IsRunning())
    
    err = bridge.Stop()
    assert.NoError(t, err)
    assert.False(t, bridge.IsRunning())
}

func TestBridge_JSONRPCRequest(t *testing.T) {
    bridge := NewPythonBridge()
    bridge.Start()
    defer bridge.Stop()
    
    request := JSONRPCRequest{
        Method: "chat",
        Params: map[string]interface{}{
            "message": "test",
            "model": "test-model",
        },
    }
    
    response, err := bridge.Send(request, 5*time.Second)
    assert.NoError(t, err)
    assert.NotNil(t, response)
}

func TestBridge_Timeout(t *testing.T) {
    bridge := NewPythonBridge()
    bridge.Start()
    defer bridge.Stop()
    
    // Send request with very short timeout
    _, err := bridge.Send(request, 1*time.Millisecond)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "timeout")
}
```

**Config:**
```go
// internal/config/config_test.go

func TestConfig_LoadSave(t *testing.T) {
    cfg := DefaultConfig()
    cfg.Providers.OpenRouter.DefaultModel = "test-model"
    
    err := cfg.Save("/tmp/test-config.yaml")
    assert.NoError(t, err)
    
    loaded, err := LoadConfig("/tmp/test-config.yaml")
    assert.NoError(t, err)
    assert.Equal(t, "test-model", loaded.Providers.OpenRouter.DefaultModel)
}

func TestConfig_EnvOverride(t *testing.T) {
    os.Setenv("MYCLI_MODEL", "env-model")
    defer os.Unsetenv("MYCLI_MODEL")
    
    cfg := LoadConfigWithEnv()
    assert.Equal(t, "env-model", cfg.Providers.OpenRouter.DefaultModel)
}

func TestConfig_Validation(t *testing.T) {
    cfg := Config{
        Behavior: BehaviorConfig{
            TokenBudget: 500,  // too low
        },
    }
    
    err := cfg.Validate()
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "token_budget must be >= 1000")
}
```

#### Unit Tests (Python)

**Providers:**
```python
# tests/python/test_providers.py

def test_openrouter_request(mocker):
    mock_response = {
        "choices": [{"message": {"content": "test response"}}],
        "usage": {"total_tokens": 100}
    }
    mocker.patch('openai.ChatCompletion.create', return_value=mock_response)
    
    provider = OpenRouterProvider(api_key="test-key")
    response = provider.chat("test message", model="test-model")
    
    assert response.content == "test response"
    assert response.tokens == 100

def test_ollama_fallback(mocker):
    # OpenRouter fails
    mocker.patch('openai.ChatCompletion.create', side_effect=Exception("API error"))
    
    # Ollama succeeds
    mock_ollama = mocker.patch('ollama.chat')
    mock_ollama.return_value = {"message": {"content": "fallback response"}}
    
    provider = ProviderManager(openrouter_key="test", ollama_enabled=True)
    response = provider.chat("test message")
    
    assert response.content == "fallback response"
    assert mock_ollama.called
```

**Parser:**
```python
# tests/python/test_parser.py

def test_extract_code_blocks():
    response = """
    Here's the fix:
    
    ```python
    def fixed_function():
        return True
    ```
    """
    
    blocks = extract_code_blocks(response)
    assert len(blocks) == 1
    assert blocks[0].language == "python"
    assert "fixed_function" in blocks[0].code

def test_detect_file_changes():
    response = """
    I'll update db.py:
    
    ```python
    # db.py
    def new_function():
        pass
    ```
    """
    
    changes = detect_file_changes(response)
    assert len(changes) == 1
    assert changes[0].file_path == "db.py"
    assert "new_function" in changes[0].content
```

#### Integration Tests

**End-to-End Flow:**
```go
// tests/go/integration_test.go

func TestE2E_ChatFlow(t *testing.T) {
    // Setup: mock API server
    mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "choices": []map[string]interface{}{
                {"message": map[string]string{"content": "test response"}},
            },
        })
    }))
    defer mockServer.Close()
    
    // Configure to use mock server
    cfg := DefaultConfig()
    cfg.Providers.OpenRouter.BaseURL = mockServer.URL
    
    // Create app
    app := NewApp(cfg)
    
    // Send message
    response, err := app.SendMessage("test message")
    assert.NoError(t, err)
    assert.Equal(t, "test response", response.Content)
}

func TestE2E_FileModification(t *testing.T) {
    // Create temp file
    tmpFile := "/tmp/test.py"
    ioutil.WriteFile(tmpFile, []byte("original content"), 0644)
    defer os.Remove(tmpFile)
    
    // Mock API response with code change
    mockResponse := `
    Here's the fix:
    
    ` + "```python\n# test.py\nnew content\n```"
    
    // Apply changes
    app := NewApp(DefaultConfig())
    err := app.ApplyChanges(mockResponse, InteractiveMode)
    assert.NoError(t, err)
    
    // Verify file changed
    content, _ := ioutil.ReadFile(tmpFile)
    assert.Contains(t, string(content), "new content")
    
    // Verify backup exists
    backupExists := fileExists("~/.mycli/backups/*/test.py")
    assert.True(t, backupExists)
}
```

#### Manual Testing Checklist (Pre-MVP Release)

**Installation:**
- [ ] Fresh install on Ubuntu 22.04
- [ ] Fresh install on macOS 13 (Intel)
- [ ] Fresh install on macOS 14 (Apple Silicon)
- [ ] Fresh install on Windows 11 (Git Bash)
- [ ] Fresh install on Windows 11 (WSL2)
- [ ] Install without Python (should show error)
- [ ] Install with Python 3.7 (should show version error)

**First Run:**
- [ ] `mycli init` interactive setup
- [ ] API key stored in keychain
- [ ] Config file created
- [ ] Can skip setup and configure later

**Chat Flow:**
- [ ] Send simple message, receive response
- [ ] Send message about code, AI analyzes files
- [ ] Multi-turn conversation maintains context
- [ ] Slash commands work (/help, /files, /model)
- [ ] Token count updates in header
- [ ] Cost tracking accurate

**Code Modifications:**
- [ ] Safe mode: shows suggestions only
- [ ] Interactive mode: shows diff, prompts for approval
- [ ] Auto mode: applies changes immediately
- [ ] Backup created before changes
- [ ] `mycli undo` restores from backup
- [ ] Git staging works (if in repo)

**Error Scenarios:**
- [ ] Invalid API key: shows helpful error
- [ ] Network offline: falls back to Ollama
- [ ] Ollama not running: shows error
- [ ] Rate limit: retries with backoff
- [ ] File permission denied: shows error
- [ ] LSP crash: continues without LSP

**Terminal Compatibility:**
- [ ] iTerm2 (macOS): full colors, animations
- [ ] Terminal.app (macOS): basic colors
- [ ] Windows Terminal: full colors
- [ ] Git Bash: basic colors
- [ ] gnome-terminal (Linux): full colors
- [ ] tmux: colors work correctly
- [ ] screen: colors work correctly

**Performance:**
- [ ] Startup time <2s (cold)
- [ ] Startup time <500ms (warm)
- [ ] First token <3s (with good network)
- [ ] Memory usage <200MB
- [ ] No memory leaks (run for 1 hour)

### Performance Targets

**Startup Performance:**
- Cold start (first run): <2 seconds
- Warm start (subsequent): <500ms
- Python subprocess start: <300ms

**Response Performance:**
- Context building: <500ms (for 10 files)
- LSP queries: <200ms per query
- Token optimization: <100ms per file
- First token from API: <3s (network dependent)
- Streaming: 30-50 tokens/second

**Resource Usage:**
- Go binary size: <20MB
- Memory (Go): ~50MB
- Memory (Python): ~100MB
- Total memory: <200MB (vs Claude Code ~1GB)
- Disk cache: <100MB

**Token Efficiency:**
- Baseline (no optimization): 50K tokens/conversation
- With optimization: 15K tokens/conversation
- **Target savings: 60-70%**

**Cost Efficiency:**
- Average session (10 turns): <$0.10
- Heavy session (50 turns): <$0.50
- Monthly (100 sessions): <$10

### Build and CI/CD

**Makefile:**
```makefile
.PHONY: build test install clean

build:
	go build -o bin/mycli cmd/mycli/main.go

test:
	go test ./...
	cd python && pytest

install:
	go install cmd/mycli

clean:
	rm -rf bin/
	rm -rf python/__pycache__

release:
	./scripts/build.sh
	./scripts/release.sh
```

**GitHub Actions (build.yml):**
```yaml
name: Build
on: [push, pull_request]

jobs:
  build:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: make build
      - run: make test
```
## Future Roadmap (Post-MVP)

### Phase 2: Enhanced AI Support

**Additional Providers:**
- Anthropic Claude (direct API)
- OpenAI GPT-4/GPT-4o
- Google Gemini
- Cohere Command
- Local models via llama.cpp

**Provider Features:**
- Multi-provider fallback chain
- Cost-based auto-selection
- A/B testing between models
- Provider health monitoring

### Phase 3: Advanced Features

**Codebase Intelligence:**
- Full project indexing (vector embeddings)
- Semantic code search
- Cross-file refactoring suggestions
- Architecture visualization

**Collaboration:**
- Shared conversation sessions
- Team configuration templates
- Usage analytics dashboard
- Cost allocation by team member

**IDE Integration:**
- VS Code extension
- JetBrains plugin
- Neovim plugin
- Emacs mode

**Git Integration:**
- Automatic commit message generation
- PR description generation
- Code review assistance
- Merge conflict resolution

### Phase 4: Plugin System

**Plugin Architecture:**
- Plugin discovery and installation
- Custom AI providers
- Custom commands
- Custom TUI themes

**Community Plugins:**
- Language-specific assistants
- Framework-specific helpers
- Custom token optimizers
- Integration with external tools

### Phase 5: Enterprise Features

**Security:**
- SSO/SAML authentication
- Audit logging
- Data residency controls
- On-premise deployment

**Management:**
- Centralized configuration
- Usage quotas and limits
- Cost budgets and alerts
- Team analytics

## Open Questions and Decisions

### Technical Decisions

**Q: Should we support streaming in the TUI?**
- **Decision:** Yes, for better UX. Show tokens as they arrive.
- **Trade-off:** More complex rendering logic, but worth it.

**Q: How to handle very large files (>1MB)?**
- **Decision:** Warn user, suggest chunking or summarization.
- **Implementation:** Add `max_file_size` config option (default: 1MB).

**Q: Should we cache LSP results?**
- **Decision:** Yes, with 1-hour TTL.
- **Rationale:** LSP queries are expensive, caching improves performance.

**Q: How to handle multiple Python versions?**
- **Decision:** Detect and use highest available (3.8+).
- **Fallback:** Let user specify via `MYCLI_PYTHON` env var.

### Product Decisions

**Q: Should we have a web UI?**
- **Decision:** Not for MVP. Terminal-first philosophy.
- **Future:** Maybe a web dashboard for analytics/settings.

**Q: Should we support non-code files (docs, configs)?**
- **Decision:** Yes, treat as plain text.
- **Limitation:** No special handling for MVP.

**Q: Should we have a free tier?**
- **Decision:** Tool is free, users pay for API usage.
- **Rationale:** Aligns with "budget-friendly" positioning.

**Q: Should we collect telemetry?**
- **Decision:** Opt-in only, privacy-first.
- **Data:** Anonymous usage stats (commands used, errors encountered).
- **No PII:** Never collect code, API keys, or personal data.

## Success Criteria

### MVP Launch (Week 0)

**Must Have:**
- ✓ Installation works on Linux, macOS, Windows
- ✓ Chat functionality with OpenRouter
- ✓ Basic token optimization (60%+ savings)
- ✓ Interactive code modification mode
- ✓ LSP integration for Go, Python, JS/TS
- ✓ Modern TUI with colors and animations
- ✓ Secure API key storage

**Nice to Have:**
- Ollama integration
- Syntax highlighting
- Slash commands
- Cost tracking

### Week 4 Metrics

**Adoption:**
- 500+ GitHub stars
- 100+ active users
- 50+ Discord/community members

**Quality:**
- <5% error rate (API calls)
- <10 critical bugs reported
- >4.0 average rating (if on marketplace)

**Performance:**
- 90th percentile response time <5s
- Average token savings >60%
- Average session cost <$0.10

### Month 3 Metrics

**Growth:**
- 1000+ GitHub stars
- 500+ weekly active users
- 10+ community contributions (PRs, issues)

**Retention:**
- 60%+ weekly active users return
- 40%+ monthly active users return
- <20% churn rate

**Revenue (if applicable):**
- Not applicable for MVP (free tool)
- Future: Premium features, enterprise licenses

## Risk Assessment

### Technical Risks

**Risk: Python dependency issues**
- **Likelihood:** Medium
- **Impact:** High (blocks installation)
- **Mitigation:** Clear error messages, installation guide, fallback to system Python

**Risk: LSP server crashes**
- **Likelihood:** Medium
- **Impact:** Low (graceful degradation)
- **Mitigation:** Catch crashes, disable LSP, continue without it

**Risk: API rate limits**
- **Likelihood:** High (for heavy users)
- **Impact:** Medium (delays responses)
- **Mitigation:** Exponential backoff, Ollama fallback, user warnings

**Risk: Token optimization too aggressive**
- **Likelihood:** Low
- **Impact:** Medium (poor AI responses)
- **Mitigation:** Configurable budget, user feedback loop, A/B testing

### Product Risks

**Risk: Users prefer GUI tools**
- **Likelihood:** Medium
- **Impact:** High (low adoption)
- **Mitigation:** Target terminal-native developers, emphasize speed/cost benefits

**Risk: Competitors copy features**
- **Likelihood:** High
- **Impact:** Low (open source, community-driven)
- **Mitigation:** Focus on execution, community, and continuous innovation

**Risk: AI providers change pricing**
- **Likelihood:** Medium
- **Impact:** Medium (affects value proposition)
- **Mitigation:** Multi-provider support, local model fallback

## Next Steps

### Immediate (This Week)

1. **Finalize spec** - Review and approve this document
2. **Create implementation plan** - Break down into tasks
3. **Set up repository** - Initialize Go project, Python package
4. **Design JSON-RPC protocol** - Define exact message format
5. **Prototype TUI** - Basic Bubble Tea app with mock data

### Short Term (Weeks 1-2)

1. **Core Engine** - File manager, context builder, token optimizer
2. **Python Bridge** - Subprocess management, JSON-RPC
3. **OpenRouter Integration** - API client, streaming support
4. **Basic TUI** - Chat interface, message display
5. **Configuration** - Config file, keychain integration

### Medium Term (Weeks 3-4)

1. **LSP Integration** - Go, Python, JS/TS support
2. **Code Modification** - Diff preview, file writer, backup
3. **Error Handling** - Retry logic, fallbacks, user messages
4. **Testing** - Unit tests, integration tests
5. **Documentation** - README, architecture docs, API docs

### Long Term (Weeks 5-8)

1. **Ollama Integration** - Local model support
2. **Advanced TUI** - Animations, syntax highlighting, themes
3. **Installer Script** - Cross-platform installation
4. **Polish** - Bug fixes, performance optimization, UX improvements
5. **Launch** - GitHub release, announcement, community building

## Conclusion

This design specification outlines a budget-friendly AI coding assistant that prioritizes cost efficiency through smart token management while maintaining a modern, user-friendly terminal interface. The hybrid Go/Python architecture leverages the strengths of both languages, and the three-layer token optimization strategy provides significant cost savings without sacrificing code understanding quality.

The MVP focuses on core functionality: chat-based interaction, intelligent context building, and safe code modifications. Post-MVP phases will expand AI provider support, add advanced features, and potentially introduce enterprise capabilities.

Success will be measured by adoption (GitHub stars, active users), quality (low error rates, high ratings), and efficiency (token savings, low costs). The project's open-source nature and terminal-first philosophy position it as a compelling alternative to existing tools for developers who value speed, cost-efficiency, and control.

**Key Differentiators:**
1. **60-70% token savings** through intelligent optimization
2. **Modern TUI** with colors, animations, and great UX
3. **Budget-friendly models** (DeepSeek-R1, StepFun) as defaults
4. **Fast performance** (<2s startup, <200MB memory)
5. **Open source** and community-driven

**Next Action:** Proceed to implementation planning phase using the `writing-plans` skill.
