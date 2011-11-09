package main

func nextField(line string) (string, string, bool) {
	fs := strings.SplitN(line, " ", 2)
	if len(fs) > 1 {
		return fs[0], fs[1], true
	}
	return fs[0], "", false
}

