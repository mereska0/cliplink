package tui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mereska0/cliplink/api/gen/linkpb"
)

type createLinkMsg struct {
	link *linkpb.Link
	err  error
}

type listLinksMsg struct {
	links []*linkpb.Link
	err   error
}

type deleteLinkMsg struct {
	err error
}

func (m Model) createLinkCmd() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		link, err := m.client.CreateLink(ctx, m.inputURL, m.inputAlias)

		return createLinkMsg{
			link: link,
			err:  err,
		}
	}
}

func (m Model) listLinksCmd() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		links, err := m.client.ListLinks(ctx)

		return listLinksMsg{
			links: links,
			err:   err,
		}
	}
}

func (m Model) deleteLinkCmd() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := m.client.DeleteLink(ctx, m.inputCode)

		return deleteLinkMsg{
			err: err,
		}
	}
}
