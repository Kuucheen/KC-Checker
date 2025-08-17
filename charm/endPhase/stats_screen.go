package endPhase

import (
	"KC-Checker/charm/threadPhase"
	"KC-Checker/common"
	"KC-Checker/helper"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type tickMsg time.Time

var (
	titleStyle       = lipgloss.NewStyle().Background(lipgloss.Color("#4B2D85"))
	successStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#57CC99"))
	warnStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#F4A261"))
	errorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#D14D4D"))
	borderRightStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderRight(true)
	centerStyle      = lipgloss.NewStyle().Width(getWidth()).Align(lipgloss.Center).MarginTop(1)

	notSelectedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#504A61"))
	borderBottomStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderBottom(true)

	currentIndexStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#B393EB"))
	prevIndex         = 0
	index             = 0
	maxIndex          = 3
	maxItems          = 6
	customWindowStart int

	customEnabled = false

	savedText     = ""
	options       []item
	outputBuilder []string
)

type item struct {
	title  string
	format string

	separatorIndicator string
	separators         []string
	separatorIndex     int
}

type model struct {
}

func RunEndScreen() {
	m := model{}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	setOptions()

	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		switch msg.String() {
		case tea.KeyRight.String():
			if index < maxIndex {
				index++
			} else {
				index = 0
				customWindowStart = 0
			}
			setOptions()

		case tea.KeyLeft.String():
			if index > 0 {
				index--
			} else {
				index = maxIndex
			}
			setOptions()

		case tea.KeyDown.String():
			if customEnabled && index >= 0 {
				prevIndex = index
				index = -1
			}

		case tea.KeyUp.String():
			if customEnabled {
				index = prevIndex
			}

		case tea.KeyEnter.String():
			if index == -1 {
				if len(outputBuilder) == 0 {
					return m, nil
				}

				writeToFile(strings.Join(outputBuilder, ""))
				savedText = successStyle.Width(getWidth() / 3).Align(lipgloss.Center).Render("Saved proxies in output folder!")
				return m, nil
			}

			if index < maxIndex {
				if !customEnabled {
					writeToFile(options[index].format)

					savedText = successStyle.Width(getWidth() / 3).Align(lipgloss.Center).Render("Saved proxies in output folder!")
				} else {
					setOutputBuilder()
				}

			} else {
				customEnabled = !customEnabled
				setOptions()
				savedText = successStyle.Width(getWidth() / 3).Align(lipgloss.Center).Render("")

				if !customEnabled {
					outputBuilder = []string{}
				}
				index = 0
			}
		}

		return m, nil

	}
	return m, tickCmd()
}

func (m model) View() string {
	proxyTable := lipgloss.NewStyle().
		Width(getWidth() / 2).
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(true).
		Align(lipgloss.Center).
		Render(getProxyTable())

	leftMerged := lipgloss.JoinVertical(lipgloss.Left, getTopLeftInfo(), proxyTable)
	leftMerged = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderRight(true).Render(leftMerged)

	fastestProxies := lipgloss.NewStyle().
		Width(getWidth()/2 + 1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(true).
		Align(lipgloss.Center).
		Render(getFastestProxies())

	rightMerged := lipgloss.JoinVertical(lipgloss.Left, getTopRightInfo(), fastestProxies)

	merged := lipgloss.JoinHorizontal(lipgloss.Left, leftMerged, rightMerged)
	merged = borderBottomStyle.Render(merged)

	bottom := getSelection()

	merged = lipgloss.JoinVertical(lipgloss.Left, merged, bottom)

	savedBottomText := successStyle.Width(getWidth() / 2).Align(lipgloss.Center).Render(strings.Join(outputBuilder, ""))
	saveButton := ""

	if customEnabled {
		if index == -1 {
			saveButton = currentIndexStyle.Align(lipgloss.Center).Width(getWidth()/5 + getWidth()/30).Render("SAVE")
		} else {
			saveButton = notSelectedStyle.Align(lipgloss.Center).Width(getWidth()/5 + getWidth()/30).Render("SAVE")
		}
	}

	savedBottomText = lipgloss.JoinHorizontal(lipgloss.Left, savedBottomText, savedText, saveButton)
	savedBottomText = lipgloss.NewStyle().MarginTop(1).Render(savedBottomText)

	merged = lipgloss.JoinVertical(lipgloss.Left, merged, savedBottomText)

	return merged
}

func tickCmd() tea.Cmd {
	return tea.Tick(threadPhase.DurationBetweenRefresh, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func writeToFile(format string) {
	if common.DoBanCheck() {
		helper.Write(helper.ProxyMapFiltered, format, true, false)
	}
	helper.Write(helper.ProxyMap, format, false, false)
}

func setOutputBuilder() {
	option := &options[index]

	if outputContains(option.format) && len(option.separators) == 0 {
		outputBuilder = removeElementByValue(outputBuilder, options[index].format)
		option.separatorIndex = 0
	} else {
		shouldAppend := true
		if len(outputBuilder) > 0 {
			if option.separatorIndicator != "" {
				removedItem := false
				if outputContains(option.format) {
					if option.separatorIndex >= len(option.separators) || len(outputBuilder) < 2 {
						outputBuilder = removeElementByValue(outputBuilder, option.format)
						options[index].separatorIndex = 0
						removedItem = true
					} else {
						if getIndexOf(outputBuilder, option.format) > 0 {
							outputBuilder[getIndexOf(outputBuilder, option.format)-1] = option.separators[option.separatorIndex]
						}
					}
					shouldAppend = false

				} else if getIndexOf(outputBuilder, option.separatorIndicator)+1 == getIndexOf(outputBuilder, option.format) ||
					outputBuilder[len(outputBuilder)-1] == option.separatorIndicator {
					outputBuilder = append(outputBuilder, option.separators[option.separatorIndex])
				} else {
					outputBuilder = append(outputBuilder, ";")
					option.separatorIndex = len(option.separators) - 1
				}

				if !removedItem {
					option.separatorIndex++
				}
			} else {
				outputBuilder = append(outputBuilder, ";")
			}
		}

		if shouldAppend {
			outputBuilder = append(outputBuilder, option.format)
		}
	}
}

func getTopRightInfo() string {
	privacyMode := ""
	if common.GetConfig().PrivacyMode {
		privacyMode = successStyle.Render("enabled")
	} else {
		privacyMode = errorStyle.Render("disabled")
	}

	privacyMode = getTopItemInfoRatio("Privacy Mode", privacyMode, -1)

	copyToClipboard := ""
	if common.GetConfig().CopyToClipboard {
		copyToClipboard = successStyle.Render("enabled")
	} else {
		copyToClipboard = errorStyle.Render("disabled")
	}
	copyToClipboard = getTopItemInfoRatio("Clipboard Copy", copyToClipboard, -1)

	leftMerged := lipgloss.JoinVertical(lipgloss.Left, privacyMode, copyToClipboard)

	proxiesActiveText := ""
	if helper.GetThreadsActive() > 0 {
		proxiesActiveText = warnStyle.Render("Threads are still active")
	} else {
		proxiesActiveText = successStyle.Render("Threads finished")
	}
	proxyCounter := getTopItemInfo("Threads", strconv.Itoa(helper.GetThreadsActive())+"/"+strconv.Itoa(common.GetConfig().Threads))

	rightMerged := lipgloss.JoinVertical(lipgloss.Left, proxyCounter, proxiesActiveText)
	rightMerged = lipgloss.NewStyle().MarginRight(1).Render(rightMerged)
	rightMerged = borderRightStyle.Render(rightMerged)

	return lipgloss.JoinHorizontal(lipgloss.Right, rightMerged, leftMerged)
}

func getTopLeftInfo() string {
	totalChecked := getTopItemInfo("Total Alive Proxies", strconv.Itoa(getLenOfProxies()))
	totalChecks := getTopItemInfo("Total Checks", strconv.FormatInt(helper.GetChecksCompleted(), 10))

	totalTime := getTopItemInfo("Total Time", threadPhase.GetTimeSince().String())
	averageCPM := getTopItemInfo("Average CPM", strconv.FormatInt(helper.GetCPM(), 10))

	leftMerged := lipgloss.JoinVertical(lipgloss.Left, totalChecked, totalChecks)
	leftMerged = lipgloss.NewStyle().MarginRight(1).Render(leftMerged)
	rightMerged := lipgloss.JoinVertical(lipgloss.Left, totalTime, averageCPM)

	leftMerged = borderRightStyle.Render(leftMerged)

	bothMerged := lipgloss.JoinHorizontal(lipgloss.Right, leftMerged, rightMerged)

	return bothMerged
}

func getTopItemInfo(leftStr string, rightStr string) string {
	return getTopItemInfoRatio(leftStr, rightStr, 0)
}

func getTopItemInfoRatio(leftStr string, rightStr string, minus int) string {
	halfStyle := lipgloss.NewStyle().Width(getWidth()/8 - minus + (getWidth()/8 - len(leftStr)))
	leftStyle := leftStr
	rightStyle := halfStyle.Align(lipgloss.Right).Render(rightStr)

	merged := lipgloss.JoinHorizontal(lipgloss.Right, leftStyle, rightStyle)

	return merged
}

func getProxyTable() string {
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	p := helper.GetProxyProtocolCountMap()

	protocols := []string{"http", "https", "socks4", "socks5"}

	indexes := []int{3, 2, 1, -1}

	board := make([][]string, len(indexes))

	for i, idx := range indexes {
		row := make([]string, len(protocols))
		for j, protocol := range protocols {
			row[j] = getTableItemString(strconv.Itoa(p[protocol][idx]))
		}
		board[i] = row
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(false).
		BorderBottom(false).
		BorderRight(false).
		Rows(board...).
		StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().Padding(0, 1)
		})

	ranks := labelStyle.Render(strings.Join([]string{"              Http   ", "Https  ", "Socks4 ", "Socks5"}, " "))
	files := labelStyle.Render(strings.Join([]string{" Elite", "Anonymous", "Transparent", "Banchecked "}, "\n "))

	return "\n" + lipgloss.JoinVertical(lipgloss.Left, ranks, lipgloss.JoinHorizontal(lipgloss.Center, files, t.Render())) + "\n"
}

func getTableItemString(str string) string {
	return strings.Repeat(" ", 6-len(str)) + str
}

func getFastestProxies() string {
	title := lipgloss.NewStyle().
		Width(getWidth() + 6).
		Align(lipgloss.Center).
		MarginBottom(1).
		Render(titleStyle.Render("Fastest Proxies"))

	var allProxies []*helper.Proxy

	for _, proxyLevel := range helper.ProxyMap {
		sort.Slice(proxyLevel, func(i, j int) bool {
			if proxyLevel[i].Time == proxyLevel[j].Time {
				return proxyLevel[i].Full < proxyLevel[j].Full
			}
			return proxyLevel[i].Time < proxyLevel[j].Time
		})

		allProxies = append(allProxies, proxyLevel...)
	}

	sort.Slice(allProxies, func(i, j int) bool {
		if allProxies[i].Time == allProxies[j].Time {
			return allProxies[i].Full < allProxies[j].Full
		}
		return allProxies[i].Time < allProxies[j].Time
	})

	retString := title

	mimimumValue := min(4, len(allProxies))

	for i := 0; i < mimimumValue; i++ {
		addString := ""
		switch allProxies[i].Level {
		case 3:
			addString += threadPhase.EliteStyle.Render("E")
		case 2:
			addString += threadPhase.AnonymousStyle.Render("A")
		case 1:
			addString += threadPhase.TransparentStyle.Render("T")
		}

		ip := allProxies[i].Full

		if common.GetPrivacyMode() {
			ip = threadPhase.GetPrivateFull(allProxies[i])
		}

		times := strings.Repeat(" ", int(math.Abs(float64(21-len(ip)))))

		strTime := strconv.Itoa(allProxies[i].Time)

		ms := successStyle.Render(strings.Repeat(" ", len(strconv.Itoa(allProxies[mimimumValue-1].Time))-len(strTime)) + strTime + "ms")

		retString += fmt.Sprintf("[%s]%s %s %s\n", addString, times, ip, ms)
	}

	return retString
}

func getSelection() string {
	title := centerStyle.
		MarginBottom(1).
		Render(titleStyle.Render("How do you want to save the proxies?"))

	optionBar := ""

	for i := 0; i < len(options); i++ {
		str := ""

		if index%len(options) == i {
			str = currentIndexStyle.Render(options[i].title)
		} else {
			str = notSelectedStyle.Render(options[i].title)
		}

		optionBar = lipgloss.JoinHorizontal(lipgloss.Right, optionBar, str)
	}

	optionBar = lipgloss.NewStyle().Align(lipgloss.Center).Width(getWidth()).Render(optionBar)

	title = lipgloss.JoinVertical(lipgloss.Center, title, optionBar)

	return title
}

func setOptions() {
	style := borderBottomStyle.
		Width(getWidth() / 4).
		Align(lipgloss.Center).Render

	customStyle := borderBottomStyle.
		Width(getWidth() / 7).
		Align(lipgloss.Center)

	if customEnabled {
		allOptions := []item{
			{title: customStyle.Render("Protocol"), format: "protocol"},
			{title: customStyle.Render("Ip"), format: "ip", separatorIndicator: "protocol", separators: []string{"://", ";", " "}},
			{title: customStyle.Render("Port"), format: "port", separatorIndicator: "ip", separators: []string{":", ";", " "}},
			{title: customStyle.Render("Email"), format: "email"},
			{title: customStyle.Render("Password"), format: "password", separatorIndicator: "email", separators: []string{":", "@", ";", " "}},
			{title: customStyle.Render("Time"), format: "time"},
			{title: customStyle.Render("Country"), format: "country"},
			{title: customStyle.Render("Type"), format: "type"},
			{title: customStyle.Render("HttpVersion"), format: "httpVersion"},
		}

		total := len(allOptions)
		visible := maxItems
		if total < visible {
			visible = total
		}

		if customWindowStart < 0 {
			customWindowStart = 0
		}
		if customWindowStart > total-visible {
			customWindowStart = total - visible
		}

		lastReal := visible - 1

		if index >= lastReal && index != visible && customWindowStart+visible < total {
			customWindowStart++
			if index > 0 {
				index--
			}
		} else if index <= 0 && customWindowStart > 0 {
			customWindowStart--
			index++
		}

		start := customWindowStart
		end := start + visible
		options = append([]item{}, allOptions[start:end]...)
		if len(options)+1 == index {
			options = allOptions[:maxItems]
			customWindowStart = 0
			index = 0
		} else if len(options) == index {
			options = allOptions[len(allOptions)-maxItems:]
			customWindowStart = index
		}

		options = append(options, item{
			title:  customStyle.BorderBottom(false).MarginBottom(1).Render("CANCEL"),
			format: "cancel",
		})

		maxIndex = len(options) - 1

	} else {
		options = []item{
			{title: style("ip:port"), format: "ip:port"},
			{title: style("protocol://ip:port"), format: "protocol://ip:port"},
			{title: style("ip:port;time"), format: "ip:port;time"},
			{title: style("CUSTOM"), format: "protocol://ip:port"},
		}
		maxIndex = len(options) - 1
		customWindowStart = 0
	}
}

func getLenOfProxies() int {
	sum := 0

	for _, count := range helper.ProxyCountMap {
		sum += count
	}

	return sum
}

func getWidth() int {
	return threadPhase.GetWidth() - 6
}

func outputContains(str string) bool {
	for _, s := range outputBuilder {
		if s == str {
			return true
		}
	}

	return false
}

func removeElementByValue(slice []string, value string) []string {
	for i, v := range slice {
		if v == value {
			if i > 0 {
				return append(slice[:i-1], slice[i+1:]...)
			} else if len(slice) >= 2 {
				return append(slice[:i], slice[i+2:]...)
			}
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice // return original slice if value not found
}

func getIndexOf(slice []string, str string) int {
	for i, s := range slice {
		if s == str {
			return i
		}
	}

	return -1
}
