module nordnetmeter

go 1.19

require tinygo.org/x/drivers v0.24.0

require tinygo.org/x/tinyfont v0.3.0

require (
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/text v0.3.6 // indirect
)

// bacause of wifinina Stop() function
replace tinygo.org/x/drivers => github.com/ysoldak/tinygo-drivers v0.15.2-0.20230421090656-00144a59a758
