package medialocker

import "io"

// https://github.com/briandowns/spinner
// https://github.com/tj/go-spin
// https://github.com/rcrowley/go-metrics
// https://github.com/olekukonko/tablewriter
// https://github.com/k0kubun/pp
type Console interface {
	Say(...interface{})
	Warn(...interface{})
	Pp(interface{})
	Debug(...interface{})
	Spin(string) func()
}

type console struct {
	in           io.Reader
	out, err     io.Writer
	debug, color bool
}
