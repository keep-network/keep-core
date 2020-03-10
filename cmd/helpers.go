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
	dashes := strings.Repeat("-", maxLineLength)

	fmt.Printf(
		"%s%s\n%s\n%s\n%s\n%s%s\n\n",
		header,
		dashes,
		buildLine(maxLineLength, prefix, suffix, "Keep Random Beacon Node"),
		buildLine(maxLineLength, prefix, suffix, ""),
		buildLine(maxLineLength, prefix, suffix, fmt.Sprintf("Port: %d", port)),
		buildMultiLine(maxLineLength, prefix, suffix, "IPs : ", addrStrings),
		dashes,
	)
}

func buildLine(lineLength int, prefix, suffix string, internalContent string) string {
	contentLength := len(prefix) + len(suffix) + len(internalContent)
	padding := lineLength - contentLength

	return fmt.Sprintf(
		"%s%s%s%s",
		prefix,
		internalContent,
		strings.Repeat(" ", padding),
		suffix,
	)
}

func buildMultiLine(lineLength int, prefix, suffix, startPrefix string, lines []string) string {
	combinedLines := buildLine(lineLength, prefix+startPrefix, suffix, lines[0]) + "\n"

	startPadding := strings.Repeat(" ", len(startPrefix))
	for _, line := range lines[1:] {
		combinedLines += buildLine(lineLength, prefix+startPadding, suffix, line) + "\n"
	}

	return combinedLines
}
