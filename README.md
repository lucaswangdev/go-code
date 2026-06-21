# go-code

AI coding agent CLI inspired by Claude Code's architecture, rewritten in Go.

## Installation & Uninstallation

### Uninstall (Before Installing)

**Linux / macOS:**
```bash
sudo rm /usr/local/bin/go-code
```

**Windows:**
```powershell
Remove-Item C:\Windows\System32\go-code.exe
```

---

### Install

**Windows (PowerShell as Admin):**
```powershell
Invoke-WebRequest -Uri "https://github.com/lucaswangdev/go-code/releases/download/v0.1.0/go-code-windows-amd64.exe" -OutFile "C:\Windows\System32\go-code.exe"
```

**macOS (Intel):**
```bash
curl -L https://github.com/lucaswangdev/go-code/releases/download/v0.1.0/go-code-darwin-amd64 -o /tmp/go-code && chmod +x /tmp/go-code && sudo mv /tmp/go-code /usr/local/bin/go-code
```

**macOS (Apple Silicon/M1/M2/M3):**
```bash
curl -L https://github.com/lucaswangdev/go-code/releases/download/v0.1.0/go-code-darwin-arm64 -o /tmp/go-code && chmod +x /tmp/go-code && sudo mv /tmp/go-code /usr/local/bin/go-code
```

**Linux (amd64):**
```bash
curl -L https://github.com/lucaswangdev/go-code/releases/download/v0.1.0/go-code-linux-amd64 -o /tmp/go-code && chmod +x /tmp/go-code && sudo mv /tmp/go-code /usr/local/bin/go-code
```

**Linux (arm64):**
```bash
curl -L https://github.com/lucaswangdev/go-code/releases/download/v0.1.0/go-code-linux-arm64 -o /tmp/go-code && chmod +x /tmp/go-code && sudo mv /tmp/go-code /usr/local/bin/go-code
```

---

## Usage

```bash
# Set API key and base URL
export OPENAI_API_KEY="your-api-key"
export OPENAI_BASE_URL="https://api.minimaxi.com/v1"
export CORECODER_MODEL="MiniMax-M2.7"

# Interactive mode
go-code

# One-shot mode
go-code -p "What is 2+2?"

# Show help
go-code --help
```

### Commands (in interactive mode)

- `/help` - Show help
- `/model` - Show current model
- `/model <name>` - Switch model
- `/tokens` - Show token usage
- `/save` - Save session
- `/sessions` - List saved sessions
- `/reset` - Clear history
- `quit` - Exit
