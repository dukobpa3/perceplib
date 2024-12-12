module perceptrail/api

go 1.23.2

require (
	github.com/dukobpa3/perceplib v0.0.0
)

require (
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
)

replace (
	github.com/dukobpa3/perceplib/api => ../api
	github.com/dukobpa3/perceplib/chain => ../chain
	github.com/dukobpa3/perceplib/exiftool => ../exiftool
	github.com/dukobpa3/perceplib/logger => ../logger
)
