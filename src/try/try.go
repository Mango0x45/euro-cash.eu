package try

import (
	"log"

	"git.thomasvoss.com/euro-cash.eu/src/atexit"
)

func Try(e error) {
	if e != nil {
		log.Fatal(e)
		atexit.Exec()
	}
}

func Try2[T any](x T, e error) T {
	Try(e)
	return x
}
