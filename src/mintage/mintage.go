package mintage

type Data struct {
	Standard      []SRow
	Commemorative []CRow
}

type SRow struct {
	Year     int
	Mintmark string
	Mintages [denoms]int
}

type CRow struct {
	Year     int
	Name     string
	Mintmark string
	Mintage  int
}

const (
	TypeCirc = iota
	TypeNifc
	TypeProof
)

const (
	Unknown = -iota - 1
	Invalid
)

const denoms = 8

var cache map[string][3]Data = make(map[string][3]Data)

func ClearCache(country string) {
	if _, ok := cache[country]; ok {
		delete(cache, country)
	}
}
