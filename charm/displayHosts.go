package charm

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
		data = append(data, []string{
			val.Host,
			val.ResponseTime.String(), // Format to display two decimal places
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

			if data[row-1][1] == "999.999999ms" {
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
		SetString(ReplaceAll(common.GetHosts()[0].Host, []string{"http://", "https://"}))
	fmt.Println("Fastest host:", fastestStyle)
}
