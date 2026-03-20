package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	fishaudio "github.com/fishaudio/fish-audio-go"
	"github.com/spf13/cobra"

	"github.com/desecho/scarlett/internal/tts"
)

// SayCmd returns the say cobra command.
func SayCmd() *cobra.Command {
	var save bool

	cmd := &cobra.Command{
		Use:   "say [text]",
		Short: "Speak text aloud using TTS",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if os.Getenv("FISH_API_KEY") == "" {
				return fmt.Errorf("FISH_API_KEY environment variable is required")
			}

			if strings.TrimSpace(args[0]) == "" {
				return fmt.Errorf("text cannot be empty")
			}

			client := fishaudio.NewClient()
			ctx := context.Background()

			audio, err := tts.Convert(ctx, client, args[0])
			if err != nil {
				return fmt.Errorf("TTS conversion failed: %w", err)
			}

			if save {
				filename := fmt.Sprintf("scarlett-%s.mp3", time.Now().Format("20060102-150405"))
				if err := os.WriteFile(filename, audio, 0644); err != nil {
					return fmt.Errorf("failed to save file: %w", err)
				}
				fmt.Println("Saved", filename)
			}

			tmp, err := os.CreateTemp("", "scarlett-say-*.mp3")
			if err != nil {
				return err
			}
			defer os.Remove(tmp.Name())

			if _, err := tmp.Write(audio); err != nil {
				tmp.Close()
				return err
			}
			tmp.Close()

			return exec.CommandContext(ctx, "afplay", tmp.Name()).Run()
		},
	}

	cmd.Flags().BoolVarP(&save, "save", "s", false, "Save MP3 to current directory")
	return cmd
}
