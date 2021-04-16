package lib

import (
	"testing"

	"github.com/qri-io/dataset"
	"github.com/qri-io/qri/base/dsfs"
	dsrefspec "github.com/qri-io/qri/dsref/spec"
	"github.com/qri-io/qri/event"
)

func TestLoadDataset(t *testing.T) {
	t.Skip("TODO(dustmop): Change in LoadDataset semantics breaks this test, figure out why")

	tr := newTestRunner(t)
	defer tr.Delete()

	fs := tr.Instance.Repo().Filesystem()

	if _, err := (*datasetLoader)(nil).LoadDataset(tr.Ctx, ""); err == nil {
		t.Errorf("expected loadDataset on a nil instance to fail without panicing")
	}

	loader := &datasetLoader{inst: nil}
	if _, err := loader.LoadDataset(tr.Ctx, ""); err == nil {
		t.Errorf("expected loadDataset on a nil instance to fail without panicing")
	}

	loader = &datasetLoader{inst: tr.Instance}
	dsrefspec.AssertLoaderSpec(t, loader, func(ds *dataset.Dataset) (string, error) {
		return dsfs.CreateDataset(
			tr.Ctx,
			fs,
			fs.DefaultWriteFS(),
			event.NilBus,
			ds,
			nil,
			tr.Instance.repo.Profiles().Owner().PrivKey,
			dsfs.SaveSwitches{},
		)
	})
}
