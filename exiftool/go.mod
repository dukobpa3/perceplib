module perceptrail/exiftool

go 1.23.2

require (
	github.com/dukobpa3/perceplib v0.0.0
)

replace (
	github.com/dukobpa3/perceplib/api => ../api
	github.com/dukobpa3/perceplib/chain => ../chain
	github.com/dukobpa3/perceplib/exiftool => ../exiftool
	github.com/dukobpa3/perceplib/logger => ../logger
)
