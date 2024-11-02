package charm

import (
	"KC-Checker/charm/threadPhase"
	"KC-Checker/common"
	"KC-Checker/helper"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"strconv"
	"time"
)

const (
	maxIndex = 3
)

var (
	selectedItems     []int
	currentIndexStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#B393EB"))
	selectedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#4F3793"))
	borderStyle       = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderBottom(true)

	helpStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#504A61")).Width(threadPhase.GetWidth() / 2).Align(lipgloss.Center).MarginTop(2).Render
	checkButtonStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#3E8262")).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("#55B08C")).
				MarginTop(2)

	topStyle          = lipgloss.NewStyle().Align(lipgloss.Center).Width(threadPhase.GetWidth() / 3).Foreground(lipgloss.Color("#e3dcf7"))
	middleStyle       = topStyle.BorderStyle(lipgloss.NormalBorder()).BorderLeft(true).BorderRight(true).BorderForeground(lipgloss.Color("#6C44BB"))
	bottomBorderStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).
				MarginBottom(2).BorderForeground(lipgloss.Color("#6C44BB"))

	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#D14D4D"))
)

type model struct {
	spinner   spinner.Model
	prevIndex int
	index     int
	finished  bool
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Spinner{
		Frames: []string{
			"▱▱▱",
			"▰▱▱",
			"▰▰▱",
			"▰▰▰",
			"▱▰▰",
			"▱▱▰",
		},
		FPS: time.Second / 5, //nolint:gomnd
	}

	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#624CAB")).MarginLeft(8)
	return model{spinner: s}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEnter.String():
			if m.index == -1 {
				if len(selectedItems) > 0 && helper.AllProxiesSum > 0 {
					m.finished = true
					return m, tea.Quit
				} else {
					break
				}
			}

			if !inSelectedItems(m.index) {
				selectedItems = append(selectedItems, m.index)
			} else {
				var newSelectedItems []int
				for _, v := range selectedItems {
					if v != m.index {
						newSelectedItems = append(newSelectedItems, v)
					}
				}
				selectedItems = newSelectedItems
			}
		case tea.KeyRight.String():
			if m.index < maxIndex && m.index != -1 {
				m.index++
			}
		case tea.KeyLeft.String():
			if m.index > 0 {
				m.index--
			}
		case tea.KeyDown.String():
			if m.index != -1 {
				m.prevIndex = m.index
				m.index = -1
			}
		case tea.KeyUp.String():
			if m.index == -1 {
				m.index = m.prevIndex
			}
		case tea.KeyCtrlC.String():
			os.Exit(1)
		}
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.finished {
		return ""
	}

	firstBox := topStyle.Render(threadPhase.GetFormattedInfoStr("Threads", strconv.Itoa(common.GetConfig().Threads)))

	secBoxString := ""

	if helper.AllProxiesSum == 0 {
		secBoxString = "Proxies" + m.spinner.View()
	} else if helper.AllProxiesSum == -1 {
		secBoxString = errorStyle.Render("No proxies detected")
	} else {
		secBoxString = threadPhase.GetFormattedInfoStr("Proxies", strconv.Itoa(int(helper.AllProxiesSum)))
	}

	secBox := middleStyle.Render(secBoxString)

	thirdBoxString := ""

	if common.GetAutoOutput() != "" {
		thirdBoxString = lipgloss.NewStyle().MarginLeft(6).Foreground(lipgloss.Color("#57CC99")).Render("enabled")
	} else {
		thirdBoxString = errorStyle.MarginLeft(2).Render("disabled")
	}

	thirdBox := topStyle.Render("Autosafe", thirdBoxString)

	combined := bottomBorderStyle.Render(lipgloss.JoinHorizontal(lipgloss.Right, firstBox, secBox, thirdBox))

	style := borderStyle.
		MarginRight(threadPhase.GetWidth() / 8).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true)

	title := lipgloss.NewStyle().
		Width(threadPhase.GetWidth()).
		Align(lipgloss.Center).
		MarginBottom(2).
		Render(lipgloss.NewStyle().Background(lipgloss.Color("#4B2D85")).
			Render("What type of proxies do you want to check?"))

	httpText := style.Render("HTTP")
	httpsText := style.Render("HTTPS")
	socks4Text := style.Render("SOCKS4")
	socks5Text := borderStyle.Render("SOCKS5")

	var options = []string{httpText, httpsText, socks4Text, socks5Text}

	var selectBar = ""

	for i := 0; i < len(options); i++ {
		if m.index%len(options) == i {
			options[i] = currentIndexStyle.Render(options[i])
		} else if inSelectedItems(i) {
			options[i] = selectedStyle.Render(options[i])
		} else {
			options[i] = lipgloss.NewStyle().Foreground(lipgloss.Color("#504A61")).Render(options[i])
		}

		selectBar = lipgloss.JoinHorizontal(lipgloss.Right, selectBar, options[i])
	}

	selectBar = lipgloss.NewStyle().Align(lipgloss.Center).Width(threadPhase.GetWidth()).Render(selectBar)

	color := ""

	if m.index == -1 {
		color = "#57CC99"
	} else {
		color = "#3E8262"
	}

	selectBar = lipgloss.JoinVertical(lipgloss.Center, selectBar,
		checkButtonStyle.Foreground(lipgloss.Color(color)).BorderForeground(lipgloss.Color(color)).Render("CHECK"))

	selectBar = lipgloss.JoinVertical(lipgloss.Bottom, title, selectBar)

	selectBar = lipgloss.JoinVertical(lipgloss.Center, selectBar, helpStyle("↑ up • ↓ down • → right • ← left • enter select"))

	return lipgloss.JoinVertical(lipgloss.Center, combined, selectBar)
}

func GetProxyType() []int {
	checkForAutoUpdateSetting()

	if len(selectedItems) > 0 {
		helper.GetCleanedProxies()
		return selectedItems
	}

	go helper.GetCleanedProxies()

	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
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

func checkForAutoUpdateSetting() {
	autoSelect := common.GetConfig().AutoSelect

	if autoSelect.Http {
		selectedItems = append(selectedItems, 0)
	}
	if autoSelect.Https {
		selectedItems = append(selectedItems, 1)
	}
	if autoSelect.Socks4 {
		selectedItems = append(selectedItems, 2)
	}
	if autoSelect.Socks5 {
		selectedItems = append(selectedItems, 3)
	}
}
