package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mereska0/cliplink/api/gen/linkpb"
	"github.com/mereska0/cliplink/internal/grpcclient"
)

type screen int

const (
	screenMenu screen = iota
	screenCreateURL
	screenCreateAlias
	screenList
	screenDelete
)

type Model struct {
	client *grpcclient.LinkClient

	screen screen
	cursor int

	inputURL   string
	inputAlias string
	inputCode  string

	links []*linkpb.Link

	message string
	err     error
}

func NewModel(client *grpcclient.LinkClient) Model {
	return Model{
		client: client,
		screen: screenMenu,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
