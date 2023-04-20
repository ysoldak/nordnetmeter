package main

const (
	_trace   = true
	_verbose = false
)

func trace(s string) string {
	if _trace {
		println(s)
	}
	return s
}

func un(s string) {
	if _trace {
		println(s + " Exit")
	}
}
