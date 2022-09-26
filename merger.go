package merge

import "reflect"

type merger struct {
	*Options
	visited map[uintptr]*visit
}

// During deepMerge, must keep track of checks that are
// in progress.  The comparison algorithm assumes that all
// checks in progress are true when it reencounters them.
// Visited are stored in a map indexed by 17 * a1 + a2;
type visit struct {
	typ  reflect.Type
	next *visit
	ptr  uintptr
}

func newMerger(options *Options) *merger {
	return &merger{
		Options: options,
		visited: make(map[uintptr]*visit),
	}
}

func (*merger) merge(dst reflect.Value, src reflect.Value) error {
	
}
