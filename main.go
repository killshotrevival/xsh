package main

import (
	cmd "xsh/cmd"

	"github.com/charmbracelet/log"
)

func main() {
	log.SetReportCaller(false)    // No file/line info
	log.SetReportTimestamp(false) // Optional: remove timestamps too
	cmd.Execute()
}
