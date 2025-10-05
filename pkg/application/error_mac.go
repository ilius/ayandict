//go:build darwin

package application

import (
	"os/exec"
	"strings"
)

func showErrorMessage(msg string) {
	defer func() {
		r := recover()
		if r != nil {
			stdLogger.Printf("Panic: %v", r)
		}
	}()
	escaped := strings.ReplaceAll(strings.ReplaceAll(msg, `\`, `\\`), `"`, `\"`)
	escaped = strings.ReplaceAll(escaped, "\n", `\n`)
	escaped = strings.ReplaceAll(escaped, "<br>", `\n`)
	cmd := exec.Command(
		"/usr/bin/osascript",
		"-e",
		`display dialog "`+escaped+`" with title "Error" buttons {"OK"} with icon stop`,
	)
	err := cmd.Run()
	if err != nil {
		stdLogger.Printf(msg)
		stdLogger.Printf(err.Error())
	}
}
