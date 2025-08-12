package chainimport

import flog "github.com/flokiorg/go-flokicoin/log"

// log is a logger that is initialized with no output filters.  This
// means the package will not perform any logging by default until the caller
// requests it.
var log flog.Logger

// log is a logger that is initialized with no output filters. This means the
// package will not perform any logging by default until the caller requests it.
func init() {
	DisableLog()
}

// DisableLog disables all library log output. Logging output is disabled by
// default until either UseLogger or SetLogWriter are called.
func DisableLog() {
	log = flog.Disabled
}

// UseLogger uses a specified Logger to output package logging info.
func UseLogger(logger flog.Logger) {
	log = logger
}
