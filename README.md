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

## Commands

### `chat`

Start an interactive chat session.

```bash
./scarlett chat
```

| Flag | Description |
|------|-------------|
| `--no-tts` | Disable text-to-speech audio |

Requires `OPENAI_API_KEY`. Also requires `FISH_API_KEY` unless `--no-tts` is set.

#### Controls

- **Enter** — Send message
- **Esc / Ctrl+C** — Quit
- **Scroll** — Browse chat history

### `say`

Speak text aloud using TTS.

```bash
./scarlett say "Hello, world!"
```

| Flag | Short | Description |
|------|-------|-------------|
| `--save` | `-s` | Save MP3 to current directory (`scarlett-YYYYMMDD-HHMMSS.mp3`) |

Requires `FISH_API_KEY`.
