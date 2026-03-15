package tts

import (
	"context"
	"os"
	"os/exec"

	fishaudio "github.com/fishaudio/fish-audio-go"
)

const voiceModelID = "7e9a17104fd644bb86b91a240b4f2055"

// Convert synthesizes text to speech and returns the raw MP3 bytes.
func Convert(ctx context.Context, client *fishaudio.Client, text string) ([]byte, error) {
	return client.TTS.Convert(ctx, &fishaudio.ConvertParams{
		Text:        text,
		ReferenceID: voiceModelID,
		Format:      fishaudio.AudioFormatMP3,
	})
}

// Speak synthesizes text to speech and plays it via afplay.
func Speak(ctx context.Context, client *fishaudio.Client, text string) error {
	audio, err := Convert(ctx, client, text)
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp("", "scarlett-tts-*.mp3")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.Write(audio); err != nil {
		tmp.Close()
		return err
	}
	tmp.Close()

	cmd := exec.CommandContext(ctx, "afplay", tmp.Name())
	return cmd.Run()
}
