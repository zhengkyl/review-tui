package film

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zhengkyl/review-ssh/ui/common"
)

type Model struct {
	common common.Common
}

func New(c common.Common) *Model {
	return &Model{
		common: c,
	}
}

func (m *Model) SetSize(width, height int) {
	m.common.Width = width
	m.common.Height = height
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) View() string {
	return ""
}