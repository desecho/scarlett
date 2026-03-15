package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	fishaudio "github.com/fishaudio/fish-audio-go"
	"github.com/openai/openai-go"
	"github.com/spf13/cobra"

	"github.com/desecho/scarlett/internal/tui"
)

var noTTS bool

// ChatCmd returns the chat cobra command.
func ChatCmd() *cobra.Command {
	return chatCmd
}

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start an interactive chat session",
	RunE: func(cmd *cobra.Command, args []string) error {
		if os.Getenv("OPENAI_API_KEY") == "" {
			return fmt.Errorf("OPENAI_API_KEY environment variable is required")
		}
		if !noTTS && os.Getenv("FISH_API_KEY") == "" {
			return fmt.Errorf("FISH_API_KEY environment variable is required (or use --no-tts)")
		}

		chatClient := openai.NewClient()

		var ttsClient *fishaudio.Client
		if !noTTS {
			ttsClient = fishaudio.NewClient()
		}

		model := tui.NewModel(chatClient, ttsClient, noTTS)
		p := tea.NewProgram(model, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	chatCmd.Flags().BoolVar(&noTTS, "no-tts", false, "Disable text-to-speech audio")
}
