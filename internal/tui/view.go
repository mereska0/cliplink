package tui

import (
	"fmt"
	"strings"
)

func (m Model) View() string {
	switch m.screen {
	case screenMenu:
		return m.viewMenu()
	case screenCreateURL:
		return m.viewCreateURL()
	case screenCreateAlias:
		return m.viewCreateAlias()
	case screenList:
		return m.viewList()
	case screenDelete:
		return m.viewDelete()
	default:
		return "Unknown screen\n"
	}
}

func (m Model) viewMenu() string {
	var b strings.Builder

	b.WriteString("ClipLink\n\n")

	items := []string{
		"Create short link",
		"List links",
		"Delete link",
	}

	for i, item := range items {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		b.WriteString(fmt.Sprintf("%s %s\n", cursor, item))
	}

	b.WriteString("\n")
	b.WriteString("enter: select • q: quit\n")

	if m.message != "" {
		b.WriteString("\n")
		b.WriteString(m.message)
		b.WriteString("\n")
	}

	return b.String()
}

func (m Model) viewCreateURL() string {
	return fmt.Sprintf(
		"Create short link\n\nOriginal URL:\n%s\n\nenter: next • esc: back\n",
		m.inputURL,
	)
}

func (m Model) viewCreateAlias() string {
	return fmt.Sprintf(
		"Create short link\n\nOriginal URL:\n%s\n\nCustom alias optional:\n%s\n\nenter: create • esc: back\n",
		m.inputURL,
		m.inputAlias,
	)
}

func (m Model) viewList() string {
	var b strings.Builder

	b.WriteString("Links\n\n")

	if len(m.links) == 0 {
		b.WriteString("No links yet.\n")
	} else {
		for _, link := range m.links {
			b.WriteString(fmt.Sprintf(
				"[%s] http://localhost:8080/%s -> %s | clicks: %d\n",
				link.GetShortCode(),
				link.GetShortCode(),
				link.GetOriginalUrl(),
				link.GetClicks(),
			))
		}
	}

	b.WriteString("\n")
	b.WriteString("esc/enter/q: back\n")

	return b.String()
}

func (m Model) viewDelete() string {
	return fmt.Sprintf(
		"Delete link\n\nShort code:\n%s\n\nenter: delete • esc: back\n",
		m.inputCode,
	)
}
