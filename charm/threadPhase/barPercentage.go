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
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
)

func RunBars() {
	m := model{
		elite:       progress.New(progress.WithDefaultGradient()),
		anonymous:   progress.New(progress.WithDefaultGradient()),
		transparent: progress.New(progress.WithDefaultGradient()),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}

type tickMsg time.Time

type model struct {
	elite       progress.Model
	anonymous   progress.Model
	transparent progress.Model
	viewport    viewport.Model
}

var finished = false

func (m model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if finished {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.elite.Width = msg.Width - padding*2 - 4
		if m.elite.Width > maxWidth {
			m.elite.Width = maxWidth
		}
		m.anonymous.Width = msg.Width - padding*2 - 4
		if m.anonymous.Width > maxWidth {
			m.anonymous.Width = maxWidth
		}
		m.transparent.Width = msg.Width - padding*2 - 4
		if m.transparent.Width > maxWidth {
			m.transparent.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		sum := float64(getSum())

		eliteCmd := m.elite.SetPercent(float64(helper.EliteCount) / sum)
		anonymousCmd := m.anonymous.SetPercent(float64(helper.AnonymousCount) / sum)
		transparentCmd := m.transparent.SetPercent(float64(helper.TransparentCount) / sum)

		return m, tea.Batch(tickCmd(), eliteCmd, anonymousCmd, transparentCmd)
	// FrameMsg is sent when the elite bar wants to animate itself
	case progress.FrameMsg:
		eliteModel, eliteCmd := m.elite.Update(msg)
		m.elite = eliteModel.(progress.Model)

		anonymousModel, anonymousCmd := m.anonymous.Update(msg)
		m.anonymous = anonymousModel.(progress.Model)

		transparentModel, transparentCmd := m.transparent.Update(msg)
		m.transparent = transparentModel.(progress.Model)

		return m, tea.Batch(eliteCmd, anonymousCmd, transparentCmd)

	default:
		return m, nil
	}
}

func (m model) View() string {
	var style = lipgloss.NewStyle().
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("#758ECD")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#7189FF")).
		SetString("_  _ ____    ____ _  _ ____ ____ _  _ ____ ____ \n|_/  |    __ |    |__| |___ |    |_/  |___ |__/ \n| \\_ |___    |___ |  | |___ |___ | \\_ |___ |  \\ \n                                                ")

	pad := strings.Repeat(" ", padding)
	return style.String() + "\n" + m.renderLine("Elite") + "\n\n" +
		pad + m.elite.View() + "\n\n" + m.renderLine("Anonymous") + "\n\n" +
		pad + m.anonymous.View() + "\n\n" + m.renderLine("Transparent") + "\n\n" +
		pad + m.transparent.View() + "\n\nElite: " +
		strconv.Itoa(int(helper.EliteCount)) + "\n\nAnonymous: " +
		strconv.Itoa(int(helper.AnonymousCount)) + "\n\nTransparent: " +
		strconv.Itoa(int(helper.TransparentCount)) + "\n\nSum: " +
		strconv.Itoa(getSum()) + "\n\nInvalid: " +
		strconv.Itoa(int(helper.Invalid)) + "\n\n" +
		pad + helpStyle("Press any key to quit")
}

func (m model) renderLine(str string) string {
	title := titleStyle.Render(str)
	line := strings.Repeat("─", max(0, 50-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func getSum() int {
	sum := int(helper.EliteCount + helper.AnonymousCount + helper.TransparentCount)

	if sum == 0 {
		sum = 1
	}

	return sum
}

func SetFinished() {
	finished = true
}
