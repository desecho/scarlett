# Scarlett

A CLI chatbot powered by GPT-5 with Fish Audio text-to-speech. Uses Bubble Tea for an interactive TUI and Cobra for CLI structure.

## Setup

### Prerequisites

- Go 1.21+
- macOS (uses `afplay` for audio playback)

### Environment Variables

```bash
export OPENAI_API_KEY=your-openai-key
export FISH_API_KEY=your-fish-audio-key
```

### Build

```bash
go build -o scarlett .
```

## Usage

```bash
# Start a chat session with TTS
./scarlett chat

# Chat without audio
./scarlett chat --no-tts
```

### Controls

- **Enter** — Send message
- **Esc / Ctrl+C** — Quit
- **Scroll** — Browse chat history
