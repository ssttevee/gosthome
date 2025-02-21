package component

import (
	"context"
	"errors"

	"github.com/gosthome/gosthome/core/component/cid"
)

type testPoll struct {
	cid.CID
	*PollingComponent[testPoll, *testPoll]
	c int
}

func newTestComponent(ctx context.Context, cfg *PollingComponentConfig) (ret *testPoll, err error) {
	ret = &testPoll{
		CID: cid.MakeID("test"),
	}
	ret.PollingComponent, err = NewPollingComponent(ctx, ret, cfg)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// Setup implements Poller.
func (t *testPoll) Setup() {
	t.PollingComponent.Setup()
}

// Poll implements Poller.
func (t *testPoll) Poll() {
	t.c += 1
}

// Close implements Poller.
func (t *testPoll) Close() error {
	return errors.Join(t.PollingComponent.Close())
}

// InitializationPriority implements Poller.
func (t *testPoll) InitializationPriority() InitializationPriority {
	return InitializationPriorityProcessor
}

var _ Poller = (*testPoll)(nil)
var _ Component = (*testPoll)(nil)
