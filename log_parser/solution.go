package log_parser

import (
	"strings"
)

func ExtractColumn(logContents string, column uint8) (result string) {
	var lines []string = strings.Split(logContents, "\n")

	var columnContent []string
	for _, line := range lines {
		var lineElements []string = strings.SplitN(line, " ", 4)
		if column >= 0 && int(column+1) < len(lineElements) {
			if column == 0 {
				columnContent = append(columnContent, strings.Join(lineElements[:2], " "))
			} else {
				columnContent = append(columnContent, lineElements[column+1])
			}
		}
	}

	if len(columnContent) > 0 {
		result = strings.Join(columnContent, "\n") + "\n"
	}

	return result
}
