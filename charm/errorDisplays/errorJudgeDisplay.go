package errorDisplays

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"time"
)

var errorStyling = lipgloss.NewStyle().
	Width(GetWidth()).
	Align(lipgloss.Center).
	Foreground(lipgloss.Color("#e64553"))

func PrintErrorForJudge(ptype []string) {
	ptypeStr := strings.Join(ptype, ", ")

	fmt.Println(errorStyling.
		Render("You need to have " + ptypeStr + " judges in your settings file\n" +
			"Please restart the program"))

	sleepForever()
}

func PrintErrorForSettings(err error) {
	fmt.Println(errorStyling.
		Render("Could not read settings file\nPlease check settings.json for errors and restart the program\n")+
		"\ndetailed error message: ", err)

	sleepForever()
}

func sleepForever() {
	time.Sleep(time.Hour * 9999)
}

func GetWidth() int {
	width, _, _ := terminal.GetSize(int(os.Stdout.Fd()))
	return width + 5
}
