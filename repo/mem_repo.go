package repo

import (
	"context"
	"fmt"
	"sync"

	"github.com/qri-io/qfs"
	"github.com/qri-io/qfs/muxfs"
	"github.com/qri-io/qri/auth/key"
	"github.com/qri-io/qri/dsref"
	"github.com/qri-io/qri/event"
	"github.com/qri-io/qri/logbook"
	"github.com/qri-io/qri/profile"
)

// MemRepo is an in-memory implementation of the Repo interface
type MemRepo struct {
	*MemRefstore

	bus        event.Bus
	filesystem *muxfs.Mux
	refCache   *MemRefstore
	logbook    *logbook.Book

	profiles profile.Store

	doneWg  sync.WaitGroup
	doneCh  chan struct{}
	doneErr error
}

var _ Repo = (*MemRepo)(nil)

// NewMemRepoWithProfile creates a new in-memory repository and an empty profile
// store owned by the given profile
func NewMemRepoWithProfile(ctx context.Context, owner *profile.Profile, fs *muxfs.Mux, bus event.Bus) (*MemRepo, error) {
	keyStore, err := key.NewMemStore()
	if err != nil {
		return nil, err
	}
	pros, err := profile.NewMemStore(owner, keyStore)
	if err != nil {
		return nil, err
	}
	return NewMemRepo(ctx, pros, fs, bus)
}

// NewMemRepo creates a new in-memory repository
func NewMemRepo(ctx context.Context, pros profile.Store, fs *muxfs.Mux, bus event.Bus) (*MemRepo, error) {
	if fs.Filesystem(qfs.MemFilestoreType) == nil {
		fs.SetFilesystem(qfs.NewMemFS())
	}

	p := pros.Owner()
	book, err := logbook.NewJournal(p.PrivKey, p.Peername, bus, fs, "/mem/logbook.qfb")
	if err != nil {
		return nil, err
	}

	mr := &MemRepo{
		bus:         bus,
		filesystem:  fs,
		MemRefstore: &MemRefstore{},
		refCache:    &MemRefstore{},
		logbook:     book,
		profiles:    pros,

		doneCh: make(chan struct{}),
	}

	mr.doneWg.Add(1)
	go func() {
		<-fs.Done()
		mr.doneErr = fs.DoneErr()
		mr.doneWg.Done()
	}()

	go func() {
		mr.doneWg.Wait()
		close(mr.doneCh)
	}()

	return mr, nil
}

// ResolveRef implements the dsref.Resolver interface
func (r *MemRepo) ResolveRef(ctx context.Context, ref *dsref.Ref) (string, error) {
	if r == nil {
		return "", dsref.ErrRefNotFound
	}

	// TODO (b5) - not totally sure why, but memRepo doesn't seem to be wiring up
	// dscache correctly in in tests
	// if r.dscache != nil {
	// 	return r.dscache.ResolveRef(ctx, ref)
	// }

	if r.logbook == nil {
		return "", fmt.Errorf("cannot resolve local references without logbook")
	}
	return r.logbook.ResolveRef(ctx, ref)
}

// Bus accesses the repo's event bus
func (r *MemRepo) Bus() event.Bus {
	return r.bus
}

// Filesystem gives access to the underlying filesystem
func (r *MemRepo) Filesystem() *muxfs.Mux {
	return r.filesystem
}

// Logbook accesses the mem repo logbook
func (r *MemRepo) Logbook() *logbook.Book {
	return r.logbook
}

// RemoveLogbook drops a MemRepo's logbook pointer. MemRepo gets used in tests
// a bunch, where logbook manipulation is helpful
func (r *MemRepo) RemoveLogbook() {
	r.logbook = nil
}

// SetLogbook assigns MemRepo's logbook. MemRepo gets used in tests a bunch,
// where logbook manipulation is helpful
func (r *MemRepo) SetLogbook(book *logbook.Book) {
	r.logbook = book
}

// SetFilesystem implements QFSSetter, currently used during lib contstruction
func (r *MemRepo) SetFilesystem(fs *muxfs.Mux) {
	r.filesystem = fs
}

// RefCache gives access to the ephemeral Refstore
func (r *MemRepo) RefCache() Refstore {
	return r.refCache
}

// Profiles gives this repo's Peer interface implementation
func (r *MemRepo) Profiles() profile.Store {
	return r.profiles
}

// Done returns a channel that the repo will send on when the repo is closed
func (r *MemRepo) Done() <-chan struct{} {
	return r.doneCh
}

// DoneErr gives an error that occurred during the shutdown process
func (r *MemRepo) DoneErr() error {
	return r.doneErr
}
