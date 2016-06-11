package gen

import "os"

func FindAllDir(from string, m map[string]struct{}) {
	file, err := os.Open(from)
	if err != nil {
		panic(err)
	}

	fi, _ := file.Readdir(-1)
	for _, v := range fi {
		if v.IsDir() {
			m[from + L + v.Name()] = struct{}{}
			FindAllDir(from + L + v.Name(), m)
		}
	}
	return
}