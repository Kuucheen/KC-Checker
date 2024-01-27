package threadPhase

import (
	"KC-Checker/common"
	"KC-Checker/helper"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/crypto/ssh/terminal"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	startTime        = time.Now()
	stopTime         = time.Now()
	stoppedTime      = false
	typeStyle        = lipgloss.NewStyle().Italic(true)
	eliteStyle       = typeStyle.Copy().Foreground(lipgloss.Color("#624CAB"))
	anonymousStyle   = typeStyle.Copy().Foreground(lipgloss.Color("#57CC99"))
	transparentStyle = typeStyle.Copy().Foreground(lipgloss.Color("#4F4F4F"))
)

func getStyledQueue() string {
	getQueue := helper.GetQueue()
	queue := getQueue.Data()

	retString := ""

	for _, value := range queue {
		addString := ""

		switch value.Level {
		case 3:
			addString += eliteStyle.Render("E")
		case 2:
			addString += anonymousStyle.Render("A")
		case 1:
			addString += transparentStyle.Render("T")
		}

		ip := value.Full

		if common.GetPrivacyMode() {
			ip = getPrivateFull(value)
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

func getPrivateFull(proxy *helper.Proxy) string {
	splitted := strings.Split(proxy.Ip, ".")
	return splitted[0] + "." +
		strings.Repeat("*", len(splitted[1])) + "." +
		strings.Repeat("*", len(splitted[2])) + "." +
		strings.Repeat("*", len(splitted[3])) + ":" +
		strings.Repeat("*", 5)
}

func getStyledInfo(elite int, anon int, trans int) string {
	activeThreads := helper.GetActive()

	timeCalc := time.Now()

	if stoppedTime {
		timeCalc = stopTime
	}

	timeSince := timeCalc.Sub(startTime)

	ms := fmt.Sprintf("%02d", int(timeSince.Milliseconds()%1000/10))

	retString := getFormattedInfo("Threads:", activeThreads) + "\n" +
		getFormattedInfo("Elite:", elite) + "\n" +
		getFormattedInfo("Anonymous:", anon) + "\n" +
		getFormattedInfo("Transparent:", trans) + "\n" +
		getFormattedInfo("Invalid:", helper.GetInvalid()) + "\n" +
		strings.Repeat("─", GetWidth()/4-6) + "\n" +
		getFormattedInfoStr("Time:", strconv.Itoa(int(timeSince.Seconds()))+"."+ms+"s")

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Align(lipgloss.Center).
		PaddingRight(6).
		Width(GetWidth() / 4).
		Render(retString)
}

func GetWidth() int {
	width, _, _ := terminal.GetSize(int(os.Stdout.Fd()))
	return width + 5
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

func getFormattedInfoStr(str string, value string) string {
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
