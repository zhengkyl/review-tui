package search

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/components/poster"
	"golang.org/x/exp/slices"
)

var (
	// titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle       = lipgloss.NewStyle().PaddingLeft(1).PaddingRight(2).MarginBottom(1).BorderStyle(lipgloss.Border{Left: " "}).BorderLeft(true)
	activeItemStyle = lipgloss.NewStyle().PaddingLeft(1).PaddingRight(2).MarginBottom(1).Foreground(
		lipgloss.Color("170")).BorderStyle(lipgloss.Border{Left: "┃"}).
		BorderForeground(lipgloss.Color("176")).BorderLeft(true)
	// paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	// helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	// quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

// NOTE: Fullwidth spaces are 2 wide
const POSTER_WIDTH = 4 * 2
const POSTER_HEIGHT = 6

type item struct {
	id           int
	title        string
	overview     string
	release_date string
	poster       *poster.PosterModel
	buttons      *ButtonsModel
}

// implement list.Item
func (i item) FilterValue() string { return i.title }

type itemDelegate struct{}

// implement list.ItemDelegate
func (d itemDelegate) Height() int { return POSTER_HEIGHT }

// implement list.ItemDelegate
func (d itemDelegate) Spacing() int { return 0 }

// implement list.ItemDelegate
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	var cmds []tea.Cmd

	for _, listItem := range m.Items() {
		i := listItem.(item)

		var cmd tea.Cmd
		_, cmd = i.poster.Update(msg)
		cmds = append(cmds, cmd)

		_, cmd = i.buttons.Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

var textStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
var titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#fff"))
var subtitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
var contentStyle = lipgloss.NewStyle().MarginLeft(2)

var ellipsisPos = []rune{' ', '.', ','}

func ellipsisText(s string, max int) string {
	if max >= len(s) {
		return s
	}

	chars := []rune(s)

	// end is an exclusive bound
	var end int
	for end = max - 3; end >= 1; end-- {
		c := chars[end]
		prevC := chars[end-1]

		if slices.Contains(ellipsisPos, c) && !slices.Contains(ellipsisPos, prevC) {
			break
		}
	}

	if end == 0 {
		end = max - 3
	}

	return string(chars[:end]) + "..."
}

// implement list.ItemDelegate
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i := listItem.(item)

	contentWidth := m.Width() - itemStyle.GetHorizontalFrameSize() - POSTER_WIDTH - contentStyle.GetHorizontalFrameSize()

	// Subtract 15 to account for long word causing early newline.
	desc := ellipsisText(i.overview, contentWidth*2-15)

	var releaseYear string
	if len(i.release_date) > 4 {
		releaseYear = i.release_date[:4]
	}

	str := lipgloss.JoinHorizontal(lipgloss.Top, titleStyle.Render(i.title), " ", subtitleStyle.Render(releaseYear))

	str = lipgloss.JoinVertical(lipgloss.Left, str, textStyle.Width(contentWidth).Render(desc))

	str += "\n\n"
	str += i.buttons.View()

	str = contentStyle.Render(str)

	str = lipgloss.JoinHorizontal(lipgloss.Top, i.poster.View(), str)

	if index == m.Index() {
		// str = lipgloss.JoinHorizontal(lipgloss.Left, "> ", str)
		str = activeItemStyle.Render(str)
	} else {
		str = itemStyle.Render(str)
	}

	fmt.Fprint(w, str)
}