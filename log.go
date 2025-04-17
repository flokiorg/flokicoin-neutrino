package neutrino

import (
	"github.com/flokiorg/flokicoin-neutrino/blockntfns"
	"github.com/flokiorg/flokicoin-neutrino/chanutils"
	"github.com/flokiorg/flokicoin-neutrino/filterdb"
	"github.com/flokiorg/flokicoin-neutrino/pushtx"
	"github.com/flokiorg/flokicoin-neutrino/query"
	"github.com/flokiorg/go-flokicoin/blockchain"
	"github.com/flokiorg/go-flokicoin/connmgr"
	flog "github.com/flokiorg/go-flokicoin/log"
	"github.com/flokiorg/go-flokicoin/netaddr"
	"github.com/flokiorg/go-flokicoin/peer"
	"github.com/flokiorg/go-flokicoin/txscript"
)

// log is a logger that is initialized with no output filters.  This
// means the package will not perform any logging by default until the caller
// requests it.
var log flog.Logger

// The default amount of logging is none.
func init() {
	DisableLog()
}

// DisableLog disables all library log output.  Logging output is disabled
// by default until either UseLogger or SetLogWriter are called.
func DisableLog() {
	log = flog.Disabled
}

// UseLogger uses a specified Logger to output package logging info.
// This should be used in preference to SetLogWriter if the caller is also
// using flog.
func UseLogger(logger flog.Logger) {
	log = logger
	blockchain.UseLogger(logger)
	txscript.UseLogger(logger)
	peer.UseLogger(logger)
	netaddr.UseLogger(logger)
	blockntfns.UseLogger(logger)
	pushtx.UseLogger(logger)
	connmgr.UseLogger(logger)
	query.UseLogger(logger)
	filterdb.UseLogger(logger)
	chanutils.UseLogger(logger)
}
