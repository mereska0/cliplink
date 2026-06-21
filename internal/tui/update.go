package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)

	case createLinkMsg:
		if msg.err != nil {
			m.message = "Failed to create link: " + msg.err.Error()
			m.screen = screenMenu
			return m, nil
		}

		m.message = "Created: http://localhost:8080/" + msg.link.GetShortCode()
		m.inputURL = ""
		m.inputAlias = ""
		m.screen = screenMenu
		return m, nil

	case listLinksMsg:
		if msg.err != nil {
			m.message = "Failed to list links: " + msg.err.Error()
			m.screen = screenMenu
			return m, nil
		}

		m.links = msg.links
		m.screen = screenList
		return m, nil

	case deleteLinkMsg:
		if msg.err != nil {
			m.message = "Failed to delete link: " + msg.err.Error()
			m.screen = screenMenu
			return m, nil
		}

		m.message = "Link deleted"
		m.inputCode = ""
		m.screen = screenMenu
		return m, nil
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.screen {
	case screenMenu:
		return m.handleMenuKey(msg)
	case screenCreateURL:
		return m.handleCreateURLKey(msg)
	case screenCreateAlias:
		return m.handleCreateAliasKey(msg)
	case screenList:
		return m.handleListKey(msg)
	case screenDelete:
		return m.handleDeleteKey(msg)
	default:
		return m, nil
	}
}

func (m Model) handleMenuKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < 2 {
			m.cursor++
		}

	case "enter":
		m.message = ""

		switch m.cursor {
		case 0:
			m.screen = screenCreateURL
			m.inputURL = ""
			m.inputAlias = ""
			return m, nil
		case 1:
			return m, m.listLinksCmd()
		case 2:
			m.screen = screenDelete
			m.inputCode = ""
			return m, nil
		}
	}

	return m, nil
}

func (m Model) handleCreateURLKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.screen = screenMenu
		return m, nil
	case "enter":
		if strings.TrimSpace(m.inputURL) == "" {
			m.message = "URL is required"
			m.screen = screenMenu
			return m, nil
		}

		m.screen = screenCreateAlias
		return m, nil
	case "backspace":
		if len(m.inputURL) > 0 {
			m.inputURL = m.inputURL[:len(m.inputURL)-1]
		}
	default:
		if len(msg.Runes) > 0 {
			m.inputURL += string(msg.Runes)
		}
	}

	return m, nil
}

func (m Model) handleCreateAliasKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.screen = screenMenu
		return m, nil
	case "enter":
		return m, m.createLinkCmd()
	case "backspace":
		if len(m.inputAlias) > 0 {
			m.inputAlias = m.inputAlias[:len(m.inputAlias)-1]
		}
	default:
		m.inputAlias += msg.String()
	}

	return m, nil
}

func (m Model) handleListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", "enter", "q":
		m.screen = screenMenu
		return m, nil
	}

	return m, nil
}

func (m Model) handleDeleteKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.screen = screenMenu
		return m, nil
	case "enter":
		if strings.TrimSpace(m.inputCode) == "" {
			m.message = "Short code is required"
			m.screen = screenMenu
			return m, nil
		}

		return m, m.deleteLinkCmd()
	case "backspace":
		if len(m.inputCode) > 0 {
			m.inputCode = m.inputCode[:len(m.inputCode)-1]
		}
	default:
		m.inputCode += msg.String()
	}

	return m, nil
}
