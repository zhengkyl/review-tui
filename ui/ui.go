package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/components/textfield"
	"github.com/zhengkyl/review-ssh/ui/keymap"
	"github.com/zhengkyl/review-ssh/ui/pages/account"
	"github.com/zhengkyl/review-ssh/ui/pages/search"
	"github.com/zhengkyl/review-ssh/ui/styles"
	"github.com/zhengkyl/review-ssh/ui/util"
)

var (
	// 	testStyle = lipgloss.NewStyle().
	// 			Bold(true).
	// 			Foreground(lipgloss.Color("#FAFAFA")).
	// 			Background(lipgloss.Color("#7D56F4")).
	// 			PaddingTop(2).
	// 			PaddingLeft(4).
	// 			Width(22)
	// highlightColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	docStyle = lipgloss.NewStyle().Background(lipgloss.Color("#7D56F4")).Padding(1, 2)

// windowStyle    = lipgloss.NewStyle().BorderForeground(highlightColor)
)

var (
	// activeTabBorder = lipgloss.Border{
	// 	Top:         "─",
	// 	Bottom:      " ",
	// 	Left:        "│",
	// 	Right:       "│",
	// 	TopLeft:     "╭",
	// 	TopRight:    "╮",
	// 	BottomLeft:  "┘",
	// 	BottomRight: "└",
	// }

	// tabBorder = lipgloss.Border{
	// 	Top:         "─",
	// 	Bottom:      "─",
	// 	Left:        "│",
	// 	Right:       "│",
	// 	TopLeft:     "╭",
	// 	TopRight:    "╮",
	// 	BottomLeft:  "┴",
	// 	BottomRight: "┴",
	// }
	titleStyle     = lipgloss.NewStyle().Background(lipgloss.Color("#fb7185"))
	tabBorder      = lipgloss.NormalBorder()
	tabStyle       = lipgloss.NewStyle().Padding(0, 1).BorderForeground(lipgloss.Color("#7D56F4")).Border(tabBorder, true)
	activeTabStyle = lipgloss.NewStyle().Padding(0, 1).BorderForeground(lipgloss.Color("#7D56F4")).Border(tabBorder, true)
)

const (
	searchTab int = iota
	accountTab
	NUM_TABS
)

type UiModel struct {
	common      common.Common
	shared      *common.Shared
	tabs        []common.PageComponent
	activeTab   int
	searchField textfield.Model

	// httpClient *retryablehttp.Client
}

func New(httpClient *retryablehttp.Client) *UiModel {

	searchField := textfield.New(common.Common{
		Width:  10, // TODO get dimensions from args
		Height: 3,
	})

	searchField.CharLimit(80)
	searchField.Placeholder("(s)earch for movies...")

	return &UiModel{
		common: common.Common{
			// Width: ,
			Styles: styles.DefaultStyles(),
			KeyMap: keymap.DefaultKeyMap(),
		},
		shared: &common.Shared{
			HttpClient: *httpClient,
		},
		tabs:        make([]common.PageComponent, NUM_TABS),
		activeTab:   0,
		searchField: *searchField,
		// httpClient: httpClient,
	}
}

func (m *UiModel) SetSize(width, height int) {
	m.common.SetSize(width, height)

	m.searchField.SetSize(width-20, 3)

	for _, tab := range m.tabs {
		tab.SetSize(width, height)
	}

	// wm, hm := ui.getMargins()

	// SetSize(width - wm, height - hm)

}

func (m *UiModel) Init() tea.Cmd {

	m.tabs[searchTab] = search.New(m.common, m.shared)
	m.tabs[accountTab] = account.New(m.common, m.shared)

	m.SetSize(m.common.Width, m.common.Height)

	cmds := []tea.Cmd{
		m.tabs[searchTab].Init(),
		m.tabs[accountTab].Init(),
	}

	return tea.Batch(cmds...)
}

func (m *UiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case common.AuthState:
		m.shared.AuthState = msg
	case tea.WindowSizeMsg:
		frameW, frameH := m.common.Styles.App.GetFrameSize()

		viewW, viewH := msg.Width-frameW, msg.Height-frameH

		m.SetSize(viewW, viewH)

		// for _, tab := range m.tabs {
		// 	_, cmd := tab.Update(msg)
		// 	// m.tabs[i] = tabModel.(common.PageComponent)

		// 	tab.SetSize(viewW, viewH-4)

		// 	cmds = append(cmds, cmd)
		// }
	// Is it a key press?
	case tea.KeyMsg:
		switch m.activeTab {
		case searchTab:
		case accountTab:
		}

		if key.Matches(msg, m.common.KeyMap.NextTab) {
			m.activeTab = (m.activeTab + 1) % NUM_TABS
			// return m, nil
		} else if key.Matches(msg, m.common.KeyMap.PrevTab) {
			m.activeTab = (m.activeTab - 1 + NUM_TABS) % NUM_TABS
			// return m, nil
		} else if key.Matches(msg, m.common.KeyMap.Quit) {
			return m, tea.Quit
		}

	}

	_, cmd := m.tabs[m.activeTab].Update(msg)
	// m.tabs[m.activeTab] = tabModel.(common.PageComponent)

	cmds = append(cmds, cmd)
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, tea.Batch(cmds...)
}

var tabNames = []string{
	"watching",
	"plan to watch",
	"completed",
	"dropped",
}

func (m *UiModel) View() string {
	view := strings.Builder{}

	bar := lipgloss.JoinHorizontal(lipgloss.Center, titleStyle.Render("movielo"), m.searchField.View())
	view.WriteString(bar + "\n")

	names := []string{}

	for i, name := range tabNames {
		if i == m.activeTab {
			names = append(names, activeTabStyle.Render(name))
		} else {
			names = append(names, tabStyle.Render(name))
		}
	}

	tabs := lipgloss.JoinHorizontal(lipgloss.Top,
		names...,
	)

	view.WriteString(tabs + "\n\n")

	// view.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(m.tabs[m.activeTab]))
	// The footer

	// view = ui.
	view.WriteString(m.tabs[m.activeTab].View())
	// view = lipgloss.JoinVertical(lipgloss.Left, ui.)
	// Send the UI for rendering
	parent := m.common.Styles.App.Render(view.String())
	// return parent
	return util.RenderOverlay(parent, docStyle.Render("hello there\nthis should be an overlay\ndid it work?"), 5, 20)

}
