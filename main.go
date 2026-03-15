package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/desecho/scarlett/cmd"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "scarlett",
		Short: "Scarlett - CLI chatbot with TTS",
		Long:  "A conversational CLI chatbot powered by GPT-5 with Fish Audio text-to-speech.",
	}

	rootCmd.AddCommand(cmd.ChatCmd())
	rootCmd.AddCommand(cmd.SayCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
