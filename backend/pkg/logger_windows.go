package pkg

import (
	"os"

	"golang.org/x/sys/windows"
)

func init() {
	// https://github.com/fatih/color/blob/main/color.go

	// Opt-in for ansi color support for current process.
	// https://learn.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences#output-sequences
	for _, fd := range []uintptr{os.Stdout.Fd(), os.Stderr.Fd()} {
		var outMode uint32
		out := windows.Handle(fd)
		if err := windows.GetConsoleMode(out, &outMode); err != nil {
			return
		}
		outMode |= windows.ENABLE_PROCESSED_OUTPUT | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		_ = windows.SetConsoleMode(out, outMode)
	}
}
