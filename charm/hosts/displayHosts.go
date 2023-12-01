package hosts

import (
	"KC-Checker/common"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"os"
	"strings"
)

func DisplayHosts(hosts []common.HostTime) {
	re := lipgloss.NewRenderer(os.Stdout)
	baseStyle := re.NewStyle().Padding(0, 1)
	headerStyle := baseStyle.Copy().Foreground(lipgloss.Color("252")).Bold(true)
	selectedStyle := baseStyle.Copy().Foreground(lipgloss.Color("#01BE85")).Background(lipgloss.Color("#00432F"))
	errorStyle := baseStyle.Copy().Foreground(lipgloss.Color("#BE0101")).Background(lipgloss.Color("#430000"))

	headers := []string{"Name", "Time"}
	var data [][]string

	for _, val := range hosts {
		response := val.ResponseTime.String()

		if response == "999h0m0s" {
			response = "error"
		}

		data = append(data, []string{
			val.Host,
			response,
		})
	}

	CapitalizeHeaders := func(data []string) []string {
		for i := range data {
			data[i] = strings.ToUpper(data[i])
		}
		return data
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(re.NewStyle().Foreground(lipgloss.Color("238"))).
		Headers(CapitalizeHeaders(headers)...).
		Width(80).
		Rows(data...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return headerStyle
			}

			if data[row-1][0] == common.GetHosts()[0].Host {
				return selectedStyle
			}

			if data[row-1][1] == "error" {
				return errorStyle
			}

			if row%2 == 0 {
				return baseStyle.Copy().Foreground(lipgloss.Color("245"))
			}
			return baseStyle.Copy().Foreground(lipgloss.Color("252"))
		})
	fmt.Println(t)

	fastestStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#01BE85")).
		SetString(replaceAll(common.GetHosts()[0].Host, []string{"http://", "https://"}))
	fmt.Println("Fastest host:", fastestStyle)
}

func replaceAll(str string, list []string) string {
	for _, val := range list {
		str = strings.ReplaceAll(str, val, "")
	}
	return str
}
