package lib

import (
	"context"
	"fmt"
	"net/rpc"

	"github.com/qri-io/qri/base"
	"github.com/qri-io/qri/logbook"
	"github.com/qri-io/qri/p2p"
	"github.com/qri-io/qri/repo"
)

// LogRequests encapsulates business logic for the log
// of changes to datasets, think "git log"
// TODO (b5): switch to using an Instance instead of separate fields
type LogRequests struct {
	node *p2p.QriNode
	cli  *rpc.Client
}

// CoreRequestsName implements the Requets interface
func (r LogRequests) CoreRequestsName() string { return "log" }

// NewLogRequests creates a LogRequests pointer from either a repo
// or an rpc.Client
func NewLogRequests(node *p2p.QriNode, cli *rpc.Client) *LogRequests {
	if node != nil && cli != nil {
		panic(fmt.Errorf("both node and client supplied to NewLogRequests"))
	}
	return &LogRequests{
		node: node,
		cli:  cli,
	}
}

// LogParams defines parameters for the Log method
type LogParams struct {
	ListParams
	// Reference to data to fetch history for
	Ref string
}

// DatasetLogItem is a line item in a dataset response
type DatasetLogItem = base.DatasetLogItem

// Log returns the history of changes for a given dataset
func (r *LogRequests) Log(params *LogParams, res *[]DatasetLogItem) (err error) {
	if r.cli != nil {
		return r.cli.Call("LogRequests.Log", params, res)
	}
	ctx := context.TODO()

	if params.Ref == "" {
		return repo.ErrEmptyRef
	}
	ref, err := repo.ParseDatasetRef(params.Ref)
	if err != nil {
		return fmt.Errorf("'%s' is not a valid dataset reference", params.Ref)
	}
	// we only canonicalize the profile here, full dataset canonicalization
	// currently relies on repo's refstore, and the logbook may be a superset
	// of the refstore
	if err = repo.CanonicalizeProfile(r.node.Repo, &ref); err != nil {
		return err
	}

	// ensure valid limit value
	if params.Limit <= 0 {
		params.Limit = 25
	}
	// ensure valid offset value
	if params.Offset < 0 {
		params.Offset = 0
	}

	*res, err = base.DatasetLog(ctx, r.node.Repo, ref, params.Limit, params.Offset, true)
	return
}

// RefListParams encapsulates parameters for requests to a single reference
// that will produce a paginated result
type RefListParams struct {
	// String value of a reference
	Ref string
	// Pagination Parameters
	Offset, Limit int
}

// LogEntry is a record in a log of operations on a dataset
type LogEntry = logbook.LogEntry

// Logbook lists log entries for actions taken on a given dataset
func (r *LogRequests) Logbook(p *RefListParams, res *[]LogEntry) error {
	if r.cli != nil {
		return r.cli.Call("LogRequests.Logbook", p, res)
	}
	ctx := context.TODO()

	ref, err := repo.ParseDatasetRef(p.Ref)
	if err != nil {
		return err
	}

	if err = repo.CanonicalizeDatasetRef(r.node.Repo, &ref); err != nil {
		return err
	}
	log.Debugf("%v", ref)

	book := r.node.Repo.Logbook()
	*res, err = book.LogEntries(ctx, repo.ConvertToDsref(ref), p.Offset, p.Limit)
	return err
}

// RawLogsParams enapsulates parameters for the RawLogs methods
type RawLogsParams struct {
	// no options yet
}

// RawLogs is an alias for a human representation of a raw logbook
type RawLogs = []logbook.Log

// RawLogs encodes the full logbook as human-oriented json
func (r *LogRequests) RawLogs(p *RawLogsParams, res *RawLogs) (err error) {
	if r.cli != nil {
		return r.cli.Call("LogRequests.RawLogs", p, res)
	}
	ctx := context.TODO()

	*res, err = r.node.Repo.Logbook().RawLogs(ctx)
	return err
}
