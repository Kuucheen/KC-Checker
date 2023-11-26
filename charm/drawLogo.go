package charm

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"os"
	"os/exec"
	"runtime"
)

var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func DrawLogo() {
	CallClear()

	var style = lipgloss.NewStyle().
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("#758ECD")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#7189FF")).
		SetString("_  _ ____    ____ _  _ ____ ____ _  _ ____ ____ \n|_/  |    __ |    |__| |___ |    |_/  |___ |__/ \n| \\_ |___    |___ |  | |___ |___ | \\_ |___ |  \\ \n                                                ")

	fmt.Print(style, "\n")
}
