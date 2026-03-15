package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	fishaudio "github.com/fishaudio/fish-audio-go"
	"github.com/openai/openai-go"

	"github.com/scrap/scarlett/internal/chat"
	"github.com/scrap/scarlett/internal/tts"
)

var (
	userStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	assistantStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true)
	errorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	headerStyle    = lipgloss.NewStyle().Background(lipgloss.Color("5")).Foreground(lipgloss.Color("0")).Bold(true).Padding(0, 1)
)

type chatMessage struct {
	role    string
	content string
}

type Model struct {
	viewport        viewport.Model
	textarea        textarea.Model
	messages        []chatMessage
	history         []openai.ChatCompletionMessageParamUnion
	streaming       bool
	currentResponse string
	chatClient      openai.Client
	ttsClient       *fishaudio.Client
	noTTS           bool
	width           int
	height          int
	ready           bool
	tokenCh         chan tea.Msg
}

func NewModel(chatClient openai.Client, ttsClient *fishaudio.Client, noTTS bool) Model {
	ta := textarea.New()
	ta.Placeholder = "Type a message..."
	ta.CharLimit = 0
	ta.SetHeight(3)
	ta.ShowLineNumbers = false
	ta.Focus()
	ta.KeyMap.InsertNewline.SetEnabled(false)

	return Model{
		textarea:   ta,
		chatClient: chatClient,
		ttsClient:  ttsClient,
		noTTS:      noTTS,
		history: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are Scarlett, a helpful and friendly conversational assistant. Keep responses concise and natural."),
		},
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.streaming {
				break
			}
			val := strings.TrimSpace(m.textarea.Value())
			if val == "" {
				break
			}
			m.textarea.Reset()
			m.messages = append(m.messages, chatMessage{role: "user", content: val})
			m.history = append(m.history, openai.UserMessage(val))
			m.streaming = true
			m.currentResponse = ""
			m.updateViewport()
			return m, m.startStream()
		}

	case StreamTokenMsg:
		m.currentResponse += msg.Token
		m.updateViewport()
		return m, m.waitForToken()

	case StreamDoneMsg:
		m.streaming = false
		m.messages = append(m.messages, chatMessage{role: "assistant", content: msg.Full})
		m.history = append(m.history, openai.AssistantMessage(msg.Full))
		m.currentResponse = ""
		m.updateViewport()
		if !m.noTTS {
			return m, m.speakText(msg.Full)
		}
		return m, nil

	case StreamErrMsg:
		m.streaming = false
		m.messages = append(m.messages, chatMessage{role: "error", content: msg.Err.Error()})
		m.updateViewport()
		return m, nil

	case TTSDoneMsg:
		return m, nil

	case TTSErrMsg:
		m.messages = append(m.messages, chatMessage{role: "error", content: "TTS error: " + msg.Err.Error()})
		m.updateViewport()
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		headerHeight := 1
		inputHeight := 5
		vpHeight := m.height - headerHeight - inputHeight
		if vpHeight < 1 {
			vpHeight = 1
		}
		if !m.ready {
			m.viewport = viewport.New(m.width, vpHeight)
			m.ready = true
		} else {
			m.viewport.Width = m.width
			m.viewport.Height = vpHeight
		}
		m.textarea.SetWidth(m.width)
		m.updateViewport()
		return m, nil
	}

	var taCmd tea.Cmd
	m.textarea, taCmd = m.textarea.Update(msg)
	cmds = append(cmds, taCmd)

	var vpCmd tea.Cmd
	m.viewport, vpCmd = m.viewport.Update(msg)
	cmds = append(cmds, vpCmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) updateViewport() {
	var sb strings.Builder
	for _, msg := range m.messages {
		switch msg.role {
		case "user":
			sb.WriteString(userStyle.Render("You") + ": " + msg.content + "\n\n")
		case "assistant":
			sb.WriteString(assistantStyle.Render("Scarlett") + ": " + msg.content + "\n\n")
		case "error":
			sb.WriteString(errorStyle.Render("Error: "+msg.content) + "\n\n")
		}
	}
	if m.streaming && m.currentResponse != "" {
		sb.WriteString(assistantStyle.Render("Scarlett") + ": " + m.currentResponse + "▍\n\n")
	} else if m.streaming {
		sb.WriteString(assistantStyle.Render("Scarlett") + ": ▍\n\n")
	}
	m.viewport.SetContent(sb.String())
	m.viewport.GotoBottom()
}

func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}
	header := headerStyle.Render(fmt.Sprintf(" Scarlett Chat %s", strings.Repeat(" ", max(0, m.width-16))))
	return header + "\n" + m.viewport.View() + "\n" + m.textarea.View()
}

// startStream launches the GPT streaming goroutine and returns a cmd to wait for the first token.
func (m *Model) startStream() tea.Cmd {
	ch := make(chan tea.Msg, 64)
	m.tokenCh = ch

	history := make([]openai.ChatCompletionMessageParamUnion, len(m.history))
	copy(history, m.history)
	client := m.chatClient

	go func() {
		full, err := chat.StreamCompletion(context.Background(), client, history, func(token string) {
			ch <- StreamTokenMsg{Token: token}
		})
		if err != nil {
			ch <- StreamErrMsg{Err: err}
		} else {
			ch <- StreamDoneMsg{Full: full}
		}
		close(ch)
	}()

	return m.waitForToken()
}

// waitForToken returns a cmd that waits for the next message on the token channel.
func (m Model) waitForToken() tea.Cmd {
	ch := m.tokenCh
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return nil
		}
		return msg
	}
}

func (m Model) speakText(text string) tea.Cmd {
	client := m.ttsClient
	return func() tea.Msg {
		err := tts.Speak(context.Background(), client, text)
		if err != nil {
			return TTSErrMsg{Err: err}
		}
		return TTSDoneMsg{}
	}
}
