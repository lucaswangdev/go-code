package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/corecoder/go-code/internal/agent"
	"github.com/corecoder/go-code/internal/config"
	"github.com/corecoder/go-code/internal/llm"
	"github.com/corecoder/go-code/internal/session"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("go-code " + version)
		return
	}

	cfg := &config.Config{}
	cfg.SetDefaults()

	model := os.Getenv("CORECODER_MODEL")
	if model == "" {
		model = cfg.Model
	}
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("CORECODER_API_KEY")
	}
	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL == "" {
		baseURL = os.Getenv("OPENAI_API_BASE")
	}
	if baseURL == "" {
		baseURL = cfg.BaseURL
	}

	if apiKey == "" {
		fmt.Println("Error: No API key found.")
		fmt.Println("Set one of: OPENAI_API_KEY, CORECODER_API_KEY")
		fmt.Println("\nExamples:")
		fmt.Println("  # OpenAI")
		fmt.Println("  export OPENAI_API_KEY=sk-...")
		fmt.Println("\n  # DeepSeek")
		fmt.Println("  export OPENAI_API_KEY=sk-... OPENAI_BASE_URL=https://api.deepseek.com")
		fmt.Println("\n  # Ollama (local)")
		fmt.Println("  export OPENAI_API_KEY=ollama OPENAI_BASE_URL=http://localhost:11434/v1 CORECODER_MODEL=qwen2.5-coder")
		os.Exit(1)
	}

	llmClient := llm.NewLLM(model, apiKey, baseURL, cfg.Temperature, cfg.MaxTokens)
	ag := agent.NewAgent(llmClient, cfg.MaxContextTokens)

	args := os.Args[1:]
	if len(args) > 0 && args[0] == "-p" && len(args) > 1 {
		prompt := args[1]
		ag.Chat(prompt, func(tok string) {
			fmt.Print(tok)
		}, nil)
		fmt.Println()
		return
	}

	if len(args) > 0 && args[0] == "-r" && len(args) > 1 {
		sessionID := args[1]
		messages, loadedModel, err := session.LoadSession(sessionID)
		if err != nil {
			fmt.Printf("Error loading session: %v\n", err)
			os.Exit(1)
		}
		ag.Messages = messages
		if model == "" {
			ag.LLM.Model = loadedModel
		}
		fmt.Printf("Resumed session: %s (model: %s)\n", sessionID, ag.LLM.Model)
	}

	if len(args) > 0 && args[0] == "--help" {
		printHelp()
		return
	}

	fmt.Printf("CoreCoder CLI (go-code)\nModel: %s\nType /help for commands, quit to exit.\n\n", ag.LLM.Model)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("You > ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if input == "" {
			continue
		}

		switch input {
		case "quit", "exit", "/quit", "/exit":
			fmt.Println("Bye!")
			return
		case "/help":
			printHelp()
			continue
		case "/reset":
			ag.Reset()
			fmt.Println("Conversation reset.")
			continue
		case "/tokens":
			fmt.Printf("Tokens: %d prompt + %d completion = %d total\n",
				ag.LLM.TotalPromptTokens, ag.LLM.TotalCompletionTokens,
				ag.LLM.TotalPromptTokens+ag.LLM.TotalCompletionTokens)
			if cost := ag.LLM.EstimatedCost(); cost != nil {
				fmt.Printf("  (~%.4f USD)\n", *cost)
			}
			continue
		case "/model":
			fmt.Printf("Current model: %s\n", ag.LLM.Model)
			continue
		case "/save":
			sid, err := session.SaveSession(ag.Messages, ag.LLM.Model, "")
			if err != nil {
				fmt.Printf("Error saving session: %v\n", err)
				continue
			}
			fmt.Printf("Session saved: %s\n", sid)
			continue
		case "/sessions":
			sessions := session.ListSessions()
			if len(sessions) == 0 {
				fmt.Println("No saved sessions.")
			}
			for _, s := range sessions {
				fmt.Printf("  %s (%s, %s) %s\n", s["id"], s["model"], s["saved_at"], s["preview"])
			}
			continue
		case "/diff":
			fmt.Println("Files modified this session.")
			continue
		}

		if len(input) > 6 && input[:7] == "/model " {
			newModel := input[7:]
			ag.LLM.Model = newModel
			fmt.Printf("Switched to %s\n", newModel)
			continue
		}

		response := ag.Chat(input, func(tok string) {
			fmt.Print(tok)
		}, func(name string, args map[string]interface{}) {
			fmt.Printf("\n> %s(%v)\n", name, args)
		})
		if response != "" {
			fmt.Println()
		}
	}
}

func printHelp() {
	fmt.Println(`Commands:
  /help          Show this help
  /reset         Clear conversation history
  /model         Show current model
  /model <name>  Switch model mid-conversation
  /tokens        Show token usage
  /save          Save session to disk
  /sessions      List saved sessions
  /diff          Show files modified this session
  quit           Exit CoreCoder

Input:
  Enter          Submit message
  Ctrl+C         Interrupt`)
}
