package watch

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"time"

	. "git.thomasvoss.com/euro-cash.eu/pkg/try"
)

func File(path string, f func()) {
	impl(path, os.Stat, f)
}

func FileFS(dir fs.FS, path string, f func()) {
	impl(path, func(path string) (os.FileInfo, error) {
		return fs.Stat(dir, path)
	}, f)
}

func impl(path string, statfn func(string) (os.FileInfo, error), f func()) {
	ostat := Try2(statfn(path))

	for {
		nstat, err := statfn(path)
		switch {
		case errors.Is(err, os.ErrNotExist):
			return
		case err != nil:
			log.Println(err)
		}

		if nstat.ModTime() != ostat.ModTime() {
			f()
			ostat = nstat
		}
		time.Sleep(1 * time.Second)
	}
}
