package threadPhase

// A simple example that shows how to render an animated elite bar. In this
// example we bump the elite by 25% every two seconds, animating our
// elite bar to its new target state.
//
// It's also possible to render a elite bar in a more static fashion without
// transitions. For details on that approach see the elite-static example.

import (
	"KC-Checker/helper"
	"fmt"
	"github.com/charmbracelet/bubbles/viewport"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding = 2
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	eliteCount       = 0
	anonymousCount   = 0
	transparentCount = 0

	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
	greenStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#01BE85")).Render
)

type tickMsg time.Time

type model struct {
	elite       progress.Model
	anonymous   progress.Model
	transparent progress.Model
	percentage  progress.Model
	viewport    viewport.Model
	list        list.Model
}

var (
	finished    = false
	threadPhase = true

	outputPath = ""
)

func RunBars() {
	items := []list.Item{
		item{"ip:port", ""},
		item{title: "type://ip:port", desc: ""},
	}

	m := model{
		elite:       progress.New(progress.WithDefaultGradient()),
		anonymous:   progress.New(progress.WithDefaultGradient()),
		transparent: progress.New(progress.WithDefaultGradient()),
		percentage:  progress.New(progress.WithDefaultGradient()),
		list:        list.New(items, list.NewDefaultDelegate(), 75, 0),
	}
	m.list.Title = "How do you want to save the proxies?"
	m.list.SetShowStatusBar(false)
	m.list.SetFilteringEnabled(false)
	m.list.SetShowHelp(false)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if finished {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if threadPhase {
			if msg.String() == "q" {
				threadPhase = false
				helper.StopThreads()

			}
		} else {
			switch msg.String() {
			case tea.KeyEnter.String():
				outputPath = "\n\n" + helper.Write(helper.ProxyMap, m.list.Index())
			case tea.KeyRight.String():
				m.list.CursorDown()
			case tea.KeyLeft.String():
				m.list.CursorUp()
			}
		}
		return m, nil

	case tea.WindowSizeMsg:
		width := 55
		m.elite.Width = width
		m.anonymous.Width = width
		m.transparent.Width = width
		m.percentage.Width = width - 10
		m.list.SetWidth(msg.Width)
		return m, nil

	case tickMsg:
		sum := float64(getSum())

		eliteCount = helper.ProxyCountMap[3]
		anonymousCount = helper.ProxyCountMap[2]
		transparentCount = helper.ProxyCountMap[1]

		eliteCmd := m.elite.SetPercent(float64(eliteCount) / sum)
		anonymousCmd := m.anonymous.SetPercent(float64(anonymousCount) / sum)
		transparentCmd := m.transparent.SetPercent(float64(transparentCount) / sum)
		percentageCmd := m.percentage.SetPercent(float64(getSumWithInvalid()) / helper.ProxySum)

		return m, tea.Batch(tickCmd(), eliteCmd, anonymousCmd, transparentCmd, percentageCmd)

	case progress.FrameMsg:
		eliteModel, eliteCmd := m.elite.Update(msg)
		m.elite = eliteModel.(progress.Model)

		anonymousModel, anonymousCmd := m.anonymous.Update(msg)
		m.anonymous = anonymousModel.(progress.Model)

		transparentModel, transparentCmd := m.transparent.Update(msg)
		m.transparent = transparentModel.(progress.Model)

		percentageModel, percentageCmd := m.percentage.Update(msg)
		m.percentage = percentageModel.(progress.Model)

		return m, tea.Batch(eliteCmd, anonymousCmd, transparentCmd, percentageCmd)
	default:
		return m, nil
	}
}

func (m model) View() string {
	pad := strings.Repeat(" ", padding)
	bars := lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderRight(true).
		SetString(m.renderLine("Elite") + "\n\n" +
			pad + m.elite.View() + "\n\n" + m.renderLine("Anonymous") + "\n\n" +
			pad + m.anonymous.View() + "\n\n" + m.renderLine("Transparent") + "\n\n" +
			pad + m.transparent.View())

	percentageBar := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		PaddingLeft(3).
		PaddingRight(3).
		Width(64).
		BorderBottom(true).
		SetString("Progress  " + m.percentage.View())

	extraString := ""

	if threadPhase {
		extraString = helpStyle("Press q to stop")
	} else {
		extraString = m.list.View() + "\n" + helpStyle("→ right • ← left")
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, bars.String(),
		lipgloss.JoinVertical(lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Center, lipgloss.JoinVertical(lipgloss.Left, lipgloss.JoinHorizontal(lipgloss.Top, getStyledQueue(),
				getStyledInfo(eliteCount, anonymousCount, transparentCount)),
				percentageBar.String()), extraString), greenStyle(outputPath)),
	)
}

func (m model) renderLine(str string) string {
	title := titleStyle.Render(str)
	line := strings.Repeat("─", max(0, 57-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func getSum() int {
	sum := eliteCount + anonymousCount + transparentCount

	if sum == 0 {
		sum = 1
	}

	return sum
}

func getSumWithInvalid() int {
	return getSum() + helper.GetInvalid()
}

func SetFinished() {
	finished = true
}
