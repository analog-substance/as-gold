package util

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func UniqueSlice(slice []string) []string {
	r := map[string]bool{}
	for _, entry := range slice {
		if entry != "" {
			r[entry] = true
		}
	}

	s := make([]string, 0, len(r))
	for k, _ := range r {
		s = append(s, k)
	}
	return s
}
