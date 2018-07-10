package cmd

import (
	"fmt"
	"strings"

	multiaddr "github.com/multiformats/go-multiaddr"
)

func nodeHeader(isBootstrapNode bool, multiaddrs []multiaddr.Multiaddr, port int) {
	prefix := "| "
	suffix := " |"

	nodeName := "node"
	if isBootstrapNode {
		nodeName = "BOOTSTRAP node"
	}
	maxLineLength := len(nodeName)

	ipStrings := make([]string, 0, len(multiaddrs))
	for _, multiaddr := range multiaddrs {
		multiaddrString := multiaddr.String()
		ipStrings = append(ipStrings, multiaddr.String())
		if len(multiaddrString) > maxLineLength {
			maxLineLength = len(multiaddrString)
		}
	}

	maxLineLength += len(prefix) + len(suffix) + 6
	dashes := strings.Repeat("-", maxLineLength)

	fmt.Printf(
		"%s\n"+
			"%s\n"+
			"%s\n"+
			"%s"+
			"%s\n",
		dashes,
		buildLine(maxLineLength, prefix, suffix, fmt.Sprintf("Node: %s", nodeName)),
		buildLine(maxLineLength, prefix, suffix, fmt.Sprintf("Port: %d", port)),
		buildMultiLine(maxLineLength, prefix, suffix, "IPs : ", ipStrings),
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
	firstLine := lines[0]
	contentLength := len(prefix) + len(startPrefix) + len(suffix) + len(firstLine)
	endPadding := strings.Repeat(" ", lineLength-contentLength)
	combinedLines :=
		fmt.Sprintf(
			"%s%s%s%s%s\n",
			prefix,
			startPrefix,
			firstLine,
			endPadding,
			suffix,
		)

	startPadding := strings.Repeat(" ", len(startPrefix))
	for _, line := range lines[1:] {
		contentLength = len(prefix) + len(startPadding) + len(suffix) + len(line)
		endPadding = strings.Repeat(" ", lineLength-contentLength)

		combinedLines +=
			fmt.Sprintf(
				"%s%s%s%s%s\n",
				prefix,
				startPadding,
				line,
				endPadding,
				suffix,
			)
	}

	return combinedLines
}
