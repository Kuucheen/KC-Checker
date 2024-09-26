package threadPhase

import (
	"KC-Checker/common"
	"KC-Checker/helper"
	"github.com/charmbracelet/bubbles/viewport"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

var (
	DurationBetweenRefresh time.Duration

	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	eliteCount       = 0
	anonymousCount   = 0
	transparentCount = 0
	banCheckCount    = 0

	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Width(GetWidth() / 2).Align(lipgloss.Center).Render
	greenStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#01BE85")).Render
	barStyle   = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			Align(lipgloss.Center).
			Width((GetWidth() - 10) / 2).
			BorderBottom(true)

	extraString = helpStyle("Press q to stop")
)

type tickMsg time.Time

type model struct {
	//Displays % of good proxies that passed the bancheck
	bancheck progress.Model
	//Displays % of progress made while checking
	percentage progress.Model

	viewport viewport.Model
}

var (
	threadPhase = true

	finished = false

	outputPath = ""
)

func RunBars() {
	m := model{
		bancheck:   progress.New(progress.WithDefaultGradient()),
		percentage: progress.New(progress.WithDefaultGradient()),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	DurationBetweenRefresh = time.Duration(common.GetConfig().TimeBetweenRefresh) * time.Millisecond

	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if helper.HasFinished {
		threadPhase = false
		SetStopTime()
		finished = true
		return m, tea.Quit
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			os.Exit(200)
			return m, tea.Quit
		}
		if threadPhase {
			if msg.String() == "q" {
				threadPhase = false
				helper.StopThreads()
				finished = true
				SetStopTime()
				return m, tea.Quit
			}
		}
		return m, nil

	case tea.WindowSizeMsg:
		width := GetWidth()/2 - 14
		m.bancheck.Width = width - 10
		m.percentage.Width = width - 10
		return m, nil

	case tickMsg:
		sum := float64(getSum())

		eliteCount = helper.ProxyCountMap[3]
		anonymousCount = helper.ProxyCountMap[2]
		transparentCount = helper.ProxyCountMap[1]
		banCheckCount = helper.ProxyCountMap[-1]
		banCheckCmd := m.bancheck.SetPercent(float64(banCheckCount) / sum)
		percentageCmd := m.percentage.SetPercent(float64(getSumWithInvalid()) / helper.AllProxiesSum)

		return m, tea.Batch(tickCmd(), banCheckCmd, percentageCmd)

	case progress.FrameMsg:
		bancheckModel, bancheckCmd := m.bancheck.Update(msg)
		m.bancheck = bancheckModel.(progress.Model)

		percentageModel, percentageCmd := m.percentage.Update(msg)
		m.percentage = percentageModel.(progress.Model)

		return m, tea.Batch(bancheckCmd, percentageCmd)
	default:
		return m, nil
	}
}

func (m model) View() string {
	if finished {
		return ""
	}

	percentageBar := barStyle.Render("Progress   " + m.percentage.View())

	percentageBar = lipgloss.JoinVertical(lipgloss.Left, percentageBar, barStyle.Render("Banchecked  "+m.bancheck.View()))

	return lipgloss.JoinHorizontal(lipgloss.Top, getProxyTree(),
		lipgloss.JoinVertical(lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Left, lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.JoinHorizontal(lipgloss.Top, getStyledQueue(),
					getStyledInfo(eliteCount, anonymousCount, transparentCount)),
				percentageBar), extraString), greenStyle(outputPath)),
	)
}

func (m model) renderLine(str string) string {
	title := titleStyle.Render(str)
	line := strings.Repeat("─", max(0, GetWidth()/2-lipgloss.Width(title)+2))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func tickCmd() tea.Cmd {
	return tea.Tick(DurationBetweenRefresh, func(t time.Time) tea.Msg {
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
