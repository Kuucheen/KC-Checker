package threadPhase

import (
	"KC-Checker/charm/errorDisplays"
	"KC-Checker/common"
	"KC-Checker/helper"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
	"math"
	"strconv"
	"strings"
	"time"
)

var (
	startTime        = time.Now()
	stopTime         = time.Now()
	stoppedTime      = false
	typeStyle        = lipgloss.NewStyle().Italic(true)
	EliteStyle       = typeStyle.Foreground(lipgloss.Color("#624CAB"))
	AnonymousStyle   = typeStyle.Foreground(lipgloss.Color("#57CC99"))
	TransparentStyle = typeStyle.Foreground(lipgloss.Color("#4F4F4F"))

	tempStyle        = lipgloss.NewStyle().Width(GetWidth() / 4).Render
	borderRightStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderRight(true).Width(GetWidth() / 2).Align(lipgloss.Center)

	protocolMap     map[string]map[int]int
	itemStyle       = lipgloss.NewStyle().MarginRight(1)
	enumeratorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).MarginRight(1)
)

func getStyledQueue() string {
	getQueue := helper.GetQueue()
	queue := getQueue.Data()

	retString := ""

	for _, value := range queue {
		addString := ""

		switch value.Level {
		case 3:
			addString += EliteStyle.Render("E")
		case 2:
			addString += AnonymousStyle.Render("A")
		case 1:
			addString += TransparentStyle.Render("T")
		}

		ip := value.Full

		if common.GetPrivacyMode() {
			ip = GetPrivateFull(value)
		}

		times := strings.Repeat(" ", int(math.Abs(float64(21-len(ip)))))

		retString += fmt.Sprintf("[%s]%s %s\n", addString, times, ip)
	}

	retString += strings.Repeat("─", GetWidth()/4) + "\n" +
		getFormattedInfo("CPM:", int(helper.GetCPM()))

	return lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderRight(true).
		Align(lipgloss.Center).
		Height(7).
		Width(GetWidth() / 4).
		Render(retString)
}

func GetPrivateFull(proxy *helper.Proxy) string {
	splitted := strings.Split(proxy.Ip, ".")
	return splitted[0] + "." +
		strings.Repeat("*", len(splitted[1])) + "." +
		strings.Repeat("*", len(splitted[2])) + "." +
		strings.Repeat("*", len(splitted[3])) + ":" +
		strings.Repeat("*", 5)
}

func getStyledInfo(elite int, anon int, trans int) string {
	activeThreads := helper.GetThreadsActive()

	timeSince := GetTimeSince()

	ms := fmt.Sprintf("%02d", int(timeSince.Milliseconds()%1000/10))

	retString := getFormattedInfo("Threads:", activeThreads) + "\n" +
		getFormattedInfo("Elite:", elite) + "\n" +
		getFormattedInfo("Anonymous:", anon) + "\n" +
		getFormattedInfo("Transparent:", trans) + "\n" +
		getFormattedInfo("Invalid:", helper.GetInvalid()) + "\n" +
		strings.Repeat("─", GetWidth()/4-6) + "\n" +
		GetFormattedInfoStr("Time:", strconv.Itoa(int(timeSince.Seconds()))+"."+ms+"s")

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Align(lipgloss.Center).
		PaddingRight(6).
		Width(GetWidth() / 4).
		Render(retString)
}

func GetTimeSince() time.Duration {
	timeCalc := time.Now()

	if stoppedTime {
		timeCalc = stopTime
	}

	return timeCalc.Sub(startTime)
}

func GetWidth() int {
	return errorDisplays.GetWidth()
}

func getFormattedInfo(str string, num int) string {
	MAXLENGTH := 18

	numStr := strconv.Itoa(num)
	length := len(str) + len(numStr)

	if length > MAXLENGTH {
		length = MAXLENGTH
	}

	return str + strings.Repeat(" ", MAXLENGTH-length) + numStr
}

func GetFormattedInfoStr(str string, value string) string {
	return str + strings.Repeat(" ", 18-len(str)-len(value)) + value
}

func SetTime() {
	startTime = time.Now()
}

func SetStopTime() {
	if !stoppedTime {
		stopTime = time.Now()
		stoppedTime = true
	}
}

func getProxyTree() string {
	protocolMap = helper.GetProxyProtocolCountMap()

	t := tree.Root("Proxies").
		Child(
			tree.Root("HTTP").
				Child(
					GetFormattedInfoStr("Elite", strconv.Itoa(protocolMap["http"][3])),
					GetFormattedInfoStr("Anonymous", strconv.Itoa(protocolMap["http"][2])),
					GetFormattedInfoStr("Transparent", strconv.Itoa(protocolMap["http"][1])),
				),
			tree.Root("HTTPS").
				Child(
					GetFormattedInfoStr("Elite", strconv.Itoa(protocolMap["https"][3])),
					GetFormattedInfoStr("Anonymous", strconv.Itoa(protocolMap["https"][2])),
					GetFormattedInfoStr("Transparent", strconv.Itoa(protocolMap["https"][1])),
				),
			tree.Root("SOCKS4").
				Child(
					GetFormattedInfoStr("Elite", strconv.Itoa(protocolMap["socks4"][3])),
					GetFormattedInfoStr("Anonymous", strconv.Itoa(protocolMap["socks4"][2])),
					GetFormattedInfoStr("Transparent", strconv.Itoa(protocolMap["socks4"][1])),
				),
			tree.Root("SOCKS5").
				Child(
					GetFormattedInfoStr("Elite", strconv.Itoa(protocolMap["socks5"][3])),
					GetFormattedInfoStr("Anonymous", strconv.Itoa(protocolMap["socks5"][2])),
					GetFormattedInfoStr("Transparent", strconv.Itoa(protocolMap["socks5"][1])),
				),
		).ItemStyle(itemStyle).EnumeratorStyle(enumeratorStyle).Enumerator(tree.RoundedEnumerator)

	return borderRightStyle.Render("\n" + tempStyle(t.String()+"\n"))
}
