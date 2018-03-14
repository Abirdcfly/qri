// Package p2p implements qri peer-to-peer communication.
// This is very, very early days, with message passing sorely in need of a
// rewrite, but hey it's a start.
package p2p

import (

	// gologging "gx/ipfs/QmQvJiADDe7JR4m968MwXobTCCzUqQkP87aRHe29MEBGHV/go-logging"
	// golog "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"
	logger "github.com/ipfs/go-log"
	protocol "gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
	identify "gx/ipfs/QmefgzMbKZYsmHFkLqxgaTBG9ypeEjrdWRD5WXH4j1cWDL/go-libp2p/p2p/protocol/identify"
)

var log = logger.Logger("qri_p2p")

// QriProtocolID is the top level Protocol Identifier
const QriProtocolID = protocol.ID("/qri")

// QriServiceTag tags the type & version of the qri service
const QriServiceTag = "qri/0.0.1"

func init() {
	// LibP2P code uses golog to log messages. They log with different
	// string IDs (i.e. "swarm"). We can control the verbosity level for
	// all loggers with:
	// golog.SetAllLoggers(gologging.INFO) // Change to DEBUG for extra info
	// golog.SetLogLevel("swarm2", "error")

	// ipfs core includes a client version. seems like a good idea.
	// TODO - understand where & how client versions are used
	identify.ClientVersion = QriServiceTag
}
