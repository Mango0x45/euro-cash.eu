package try

import (
	"log"

	"git.thomasvoss.com/euro-cash.eu/pkg/atexit"
)

func Try(e error) {
	if e != nil {
		atexit.Exec()
		log.Fatalln(e)
	}
}

func Try2[T any](x T, e error) T {
	Try(e)
	return x
}
