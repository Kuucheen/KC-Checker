package charm

import (
	"KC-Checker/charm/threadPhase"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
)

const (
	maxIndex = 3
)

var (
	prevIndex         = 0
	index             = 0
	selectedItems     []int
	finished          = false
	currentIndexStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8A6EE3"))
	selectedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#624CAB"))
	borderStyle       = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderBottom(true)

	checkButtonStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#57CC99")).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("#419972")).
				MarginTop(2).
		//Background(lipgloss.Color("#419972")).
		Blink(true)
	//BorderStyle(lipgloss.ThickBorder())
)

type model struct {
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("KC-Checker - github.com/Kuucheen")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEnter.String():
			if index == -1 && len(selectedItems) > 0 {
				finished = true
				return m, tea.Quit
			}

			if !inSelectedItems(index) {
				selectedItems = append(selectedItems, index)
			} else {
				newSelectedItems := []int{}
				for _, v := range selectedItems {
					if v != index {
						newSelectedItems = append(newSelectedItems, v)
					}
				}
				selectedItems = newSelectedItems
			}
			break
		case tea.KeyRight.String():
			if index < maxIndex && index != -1 {
				index++
			}
			break
		case tea.KeyLeft.String():
			if index > 0 || index == -1 {
				index--
			}
			break
		case tea.KeyDown.String():
			if index != -1 {
				prevIndex = index
				index = -1
			}
			break
		case tea.KeyUp.String():
			if index == -1 {
				index = prevIndex
			}
			break
		case tea.KeyCtrlC.String():
			os.Exit(1)
		}
	}

	var cmd tea.Cmd
	return m, cmd
}

func (m model) View() string {
	if finished {
		return ""
	}

	style := borderStyle.Copy().
		MarginRight(threadPhase.GetWidth() / 8).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true)

	title := lipgloss.NewStyle().
		Width(threadPhase.GetWidth()).
		Align(lipgloss.Center).
		MarginBottom(2).
		Render(lipgloss.NewStyle().Background(lipgloss.Color("#5f5fd7")).
			Render("What type of proxies do you want to check?"))

	httpText := style.Render("HTTP")
	httpsText := style.Render("HTTPS")
	socks4Text := style.Render("SOCKS4")
	socks5Text := borderStyle.Render("SOCKS5")

	var options = []string{httpText, httpsText, socks4Text, socks5Text}

	var selectBar = ""

	for i := 0; i < len(options); i++ {
		if index%len(options) == i {
			options[i] = currentIndexStyle.Render(options[i])
		} else if inSelectedItems(i) {
			options[i] = selectedStyle.Render(options[i])
		} else {
			options[i] = lipgloss.NewStyle().Foreground(lipgloss.Color("#777777")).Render(options[i])
		}

		selectBar = lipgloss.JoinHorizontal(lipgloss.Right, selectBar, options[i])
	}

	selectBar = lipgloss.NewStyle().Align(lipgloss.Center).Width(threadPhase.GetWidth()).Render(selectBar)

	selectBar = lipgloss.JoinVertical(lipgloss.Center, selectBar, checkButtonStyle.Render("CHECK"))

	return lipgloss.JoinVertical(lipgloss.Bottom, title, selectBar)
}

func GetProxyType() []int {
	m := model{}
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	return selectedItems
}

func inSelectedItems(index int) bool {
	for _, item := range selectedItems {
		if item == index {
			return true
		}
	}

	return false
}
