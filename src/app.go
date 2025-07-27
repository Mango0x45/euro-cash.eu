package app

import (
	"os"
	"syscall"

	"git.thomasvoss.com/euro-cash.eu/pkg/atexit"
	. "git.thomasvoss.com/euro-cash.eu/pkg/try"
)

func Restart() {
	path := Try2(os.Executable())
	atexit.Exec()
	Try(syscall.Exec(path, append([]string{path}, os.Args[1:]...),
		os.Environ()))
}
