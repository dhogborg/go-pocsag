module github.com/dhogborg/go-pocsag

go 1.18

replace (
	github.com/dhogborg/go-pocsag/internal/datatypes => ./internal/datatypes
	github.com/dhogborg/go-pocsag/internal/pocsag => ./internal/pocsag
	github.com/dhogborg/go-pocsag/internal/util => ./internal/util
	github.com/dhogborg/go-pocsag/internal/wav => ./internal/wav
)

require (
	github.com/dhogborg/go-pocsag/internal/datatypes v0.0.0-00010101000000-000000000000
	github.com/dhogborg/go-pocsag/internal/pocsag v0.0.0-00010101000000-000000000000
	github.com/fatih/color v1.13.0
	github.com/urfave/cli v1.22.9
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.0-20190314233015-f79a8a8ca69d // indirect
	github.com/dhogborg/go-pocsag/internal/wav v0.0.0-00010101000000-000000000000 // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/kr/text v0.1.0 // indirect
	github.com/mattn/go-colorable v0.1.9 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/russross/blackfriday/v2 v2.0.1 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
)
