package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	flag.Usage = help
	flag.Parse()

	cmds := map[string]func(){
		"help":  help,
		"start": start,
	}

	if cmdFunc, ok := cmds[flag.Arg(0)]; ok {
		cmdFunc()
	} else {
		help()
		os.Exit(2)
	}
}

func help() {
	divider := "| %s | %s |\n"
	header := "| %-30s | %-50s |\n"
	row := "| %-30s | %-50s |\n"

	output :=
		fmt.Sprintf(divider, strings.Repeat("-", 30), strings.Repeat("-", 50)) +
			fmt.Sprintf(header, "Usage", "Description") +
			fmt.Sprintf(divider, strings.Repeat("-", 30), strings.Repeat("-", 50)) +
			fmt.Sprintf(row, "help", "show this help message") +
			fmt.Sprintf(row, "start", "start the server") +
			fmt.Sprintf(divider, strings.Repeat("_", 30), strings.Repeat("_", 50))

	fmt.Fprintln(os.Stderr, output)
}
