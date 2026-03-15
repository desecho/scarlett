package tui

// StreamTokenMsg carries a single token from the GPT stream.
type StreamTokenMsg struct{ Token string }

// StreamDoneMsg signals the stream finished with the full text.
type StreamDoneMsg struct{ Full string }

// StreamErrMsg signals a streaming error.
type StreamErrMsg struct{ Err error }

// TTSDoneMsg signals TTS playback finished.
type TTSDoneMsg struct{}

// TTSErrMsg signals a TTS error.
type TTSErrMsg struct{ Err error }
