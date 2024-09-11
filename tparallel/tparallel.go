//go:build !solution

package tparallel

type T struct {
	finished chan struct{}
	started  bool

	bChannel chan struct{}

	parent   *T
	subTests []*T
}

func newT(parent *T) *T {
	return &T{
		finished: make(chan struct{}),

		bChannel: make(chan struct{}),

		parent:   parent,
		subTests: make([]*T, 0),
	}
}

func (t *T) Parallel() {
	if t.started {
		panic("already paralleled")
	}
	t.started = true
	t.parent.subTests = append(t.parent.subTests, t)

	t.finished <- struct{}{}
	<-t.parent.bChannel
}

func (t *T) tRunner(subtest func(t *T)) {
	subtest(t)
	if len(t.subTests) > 0 {
		close(t.bChannel)

		for _, sub := range t.subTests {
			<-sub.finished
		}
	}
	if t.started {
		t.parent.finished <- struct{}{}
	}
	t.finished <- struct{}{}

}

func (t *T) Run(subtest func(t *T)) {
	subT := newT(t)
	go subT.tRunner(subtest)
	<-subT.finished
}

func Run(tsts []func(t *T)) {
	root := newT(nil)
	for _, fn := range tsts {
		root.Run(fn)
	}
	close(root.bChannel)

	if len(root.subTests) > 0 {
		<-root.finished
	}
}
