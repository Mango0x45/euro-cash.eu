package mintage

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func Parse(country string) ([3]Data, error) {
	if data, ok := cache[country]; ok {
		return data, nil
	}

	var (
		data [3]Data
		err  error
		path = filepath.Join("data", "mintages", country)
	)

	data[TypeCirc].Standard, err = parseS(path + "-s-circ.csv")
	if err != nil {
		return data, err
	}
	data[TypeNifc].Standard, err = parseS(path + "-s-nifc.csv")
	if err != nil {
		return data, err
	}
	data[TypeProof].Standard, err = parseS(path + "-s-proof.csv")
	if err != nil {
		return data, err
	}
	data[TypeCirc].Commemorative, err = parseC(path + "-c-circ.csv")
	if err != nil {
		return data, err
	}
	data[TypeNifc].Commemorative, err = parseC(path + "-c-nifc.csv")
	if err != nil {
		return data, err
	}
	data[TypeProof].Commemorative, err = parseC(path + "-c-proof.csv")
	if err == nil {
		cache[country] = data
	}
	return data, err
}

func parseS(path string) ([]SRow, error) {
	rows := make([]SRow, 0, guessRows(false))

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comment = '#'
	r.FieldsPerRecord = 11
	r.ReuseRecord = true

	/* Skip header */
	if _, err := r.Read(); err != nil {
		return nil, err
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		data := SRow{
			Mintmark:  record[1],
			Reference: record[10],
		}

		data.Year, err = strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}

		for i, s := range record[2:10] {
			if s == "" {
				data.Mintages[i] = Unknown
			} else {
				data.Mintages[i], err = strconv.Atoi(s)
				if err != nil {
					data.Mintages[i] = Invalid
				}
			}
		}

		rows = append(rows, data)
	}

	return rows, nil
}

func parseC(path string) ([]CRow, error) {
	rows := make([]CRow, 0, guessRows(true))

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comment = '#'
	r.FieldsPerRecord = 5
	r.ReuseRecord = true

	/* Skip header */
	if _, err := r.Read(); err != nil {
		return nil, err
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		data := CRow{
			Name:      record[1],
			Mintmark:  record[2],
			Reference: record[4],
		}

		data.Year, err = strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}

		s := record[3]
		if s == "" {
			data.Mintage = Unknown
		} else {
			data.Mintage, err = strconv.Atoi(s)
			if err != nil {
				data.Mintage = Invalid
			}
		}

		rows = append(rows, data)
	}

	return rows, nil
}

func guessRows(commemorativep bool) int {
	/* Try to guess the number of rows for Germany, because nobody needs more
	   rows than Germany. */
	n := (time.Now().Year() - 2002) * 5
	if commemorativep {
		return n * 2
	}
	return n
}
