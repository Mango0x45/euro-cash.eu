package try

import "log"

func Try(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func Try2[T any](x T, e error) T {
	Try(e)
	return x
}
