package chat

import (
	"context"

	"github.com/openai/openai-go"
)

// StreamCompletion streams a chat completion and calls onToken for each delta.
// Returns the full accumulated response text.
func StreamCompletion(ctx context.Context, client openai.Client,
	messages []openai.ChatCompletionMessageParamUnion,
	onToken func(string)) (string, error) {

	stream := client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model:    "gpt-5",
		Messages: messages,
	})

	var full string
	for stream.Next() {
		chunk := stream.Current()
		if len(chunk.Choices) > 0 {
			delta := chunk.Choices[0].Delta.Content
			if delta != "" {
				full += delta
				onToken(delta)
			}
		}
	}
	if err := stream.Err(); err != nil {
		return full, err
	}
	return full, nil
}
