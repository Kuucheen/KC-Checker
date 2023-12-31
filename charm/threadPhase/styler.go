package threadPhase

import (
	"KC-Checker/helper"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strconv"
	"strings"
)

var (
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

		times := strings.Repeat(" ", 21-len(value.Full))

		retString += fmt.Sprintf("[%s]%s %s\n", addString, times, value.Full)
	}
	if len(retString) >= 2 {
		retString = retString[:len(retString)-1]
	}

	return lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderRight(true).
		Align(lipgloss.Center).
		Height(5).
		Width(GetWidth() / 4).
		Render(retString)
}

func getStyledInfo(elite int, anon int, trans int) string {
	activeThreads := helper.GetActive()

	retString := getFormattedInfo("Threads:", activeThreads) + "\n" +
		getFormattedInfo("Elite:", elite) + "\n" +
		getFormattedInfo("Anonymous:", anon) + "\n" +
		getFormattedInfo("Transparent:", trans) + "\n" +
		getFormattedInfo("Invalid:", helper.GetInvalid())

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
	numStr := strconv.Itoa(num)
	return str + strings.Repeat(" ", 18-len(str)-len(numStr)) + numStr
}
