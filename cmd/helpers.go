package cmd

import (
	"fmt"
	"strconv"
	"strings"
)

func nodeHeader(addrStrings []string, port int) {
	header := ` 

▓▓▌ ▓▓ ▐▓▓ ▓▓▓▓▓▓▓▓▓▓▌▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▄
▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▌▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
  ▓▓▓▓▓▓    ▓▓▓▓▓▓▓▀    ▐▓▓▓▓▓▓    ▐▓▓▓▓▓   ▓▓▓▓▓▓     ▓▓▓▓▓   ▐▓▓▓▓▓▌   ▐▓▓▓▓▓▓
  ▓▓▓▓▓▓▄▄▓▓▓▓▓▓▓▀      ▐▓▓▓▓▓▓▄▄▄▄         ▓▓▓▓▓▓▄▄▄▄         ▐▓▓▓▓▓▌   ▐▓▓▓▓▓▓
  ▓▓▓▓▓▓▓▓▓▓▓▓▓▀        ▐▓▓▓▓▓▓▓▓▓▓         ▓▓▓▓▓▓▓▓▓▓▌        ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
  ▓▓▓▓▓▓▀▀▓▓▓▓▓▓▄       ▐▓▓▓▓▓▓▀▀▀▀         ▓▓▓▓▓▓▀▀▀▀         ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▀
  ▓▓▓▓▓▓   ▀▓▓▓▓▓▓▄     ▐▓▓▓▓▓▓     ▓▓▓▓▓   ▓▓▓▓▓▓     ▓▓▓▓▓   ▐▓▓▓▓▓▌
▓▓▓▓▓▓▓▓▓▓ █▓▓▓▓▓▓▓▓▓ ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  ▓▓▓▓▓▓▓▓▓▓
▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓ ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  ▓▓▓▓▓▓▓▓▓▓

Trust math, not hardware.
	
`

	prefix := "| "
	suffix := " |"

	maxLineLength := len(strconv.Itoa(port))

	for _, addrString := range addrStrings {
		if addrLength := len(addrString); addrLength > maxLineLength {
			maxLineLength = addrLength
		}
	}

	maxLineLength += len(prefix) + len(suffix) + 6
	dashes := strings.Repeat("-", maxLineLength) + "\n"

	fmt.Print(
		header,
		dashes,
		buildLine(maxLineLength, prefix, suffix, "Keep Client Node"),
		buildLine(maxLineLength, prefix, suffix, ""),
		buildLine(maxLineLength, prefix, suffix, fmt.Sprintf("Port: %d", port)),
		buildMultiLine(maxLineLength, prefix, suffix, "IPs : ", addrStrings),
		buildLine(maxLineLength, prefix, suffix, ""),
		dashes,
		"\n",
	)
}

func buildLine(lineLength int, prefix, suffix string, internalContent string) string {
	contentLength := len(prefix) + len(suffix) + len(internalContent)
	padding := lineLength - contentLength

	return fmt.Sprint(
		prefix,
		internalContent,
		strings.Repeat(" ", padding),
		suffix,
		"\n",
	)
}

func buildMultiLine(lineLength int, prefix, suffix, startPrefix string, lines []string) string {
	combinedLines := buildLine(lineLength, prefix+startPrefix, suffix, lines[0])

	startPadding := strings.Repeat(" ", len(startPrefix))
	for _, line := range lines[1:] {
		combinedLines += buildLine(lineLength, prefix+startPadding, suffix, line)
	}

	return combinedLines
}
