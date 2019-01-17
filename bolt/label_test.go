package bolt_test

import (
	"context"
	"testing"

	platform "github.com/influxdata/influxdb"
	"github.com/influxdata/influxdb/bolt"
	platformtesting "github.com/influxdata/influxdb/testing"
)

func initLabelService(f platformtesting.LabelFields, t *testing.T) (platform.LabelService, string, func()) {
	c, closeFn, err := NewTestClient()
	if err != nil {
		t.Fatalf("failed to create new bolt client: %v", err)
	}
	ctx := context.Background()
	for _, l := range f.Labels {
		if err := c.CreateLabel(ctx, l); err != nil {
			t.Fatalf("failed to populate labels")
		}
	}

	return c, bolt.OpPrefix, func() {
		defer closeFn()
		for _, l := range f.Labels {
			if err := c.DeleteLabel(ctx, *l); err != nil {
				t.Logf("failed to remove label: %v", err)
			}
		}
	}
}

func TestLabelService_LabelService(t *testing.T) {
	platformtesting.LabelService(initLabelService, t)
}
