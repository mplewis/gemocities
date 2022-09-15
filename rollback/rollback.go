// Package rollback provides a friendly interface for rolling back partial state when parent tasks fail.
package rollback

// Rollback wraps a deferred cleanup task that can be cancelled if the parent task succeeds.
type Rollback struct {
	runnable func()
}

// New creates a new Rollback with code to be run to cleanup a dirty state.
func New(runnable func()) *Rollback {
	return &Rollback{runnable}
}

// OK should be called when the parent task succeeds.
// It will prevent the cleanup task from running.
func (r *Rollback) OK() {
	r.runnable = nil
}

// Done should be deferred immediately after creating a Rollback.
// It will run the cleanup task if the parent task fails to call OK.
func (r *Rollback) Done() {
	if r.runnable != nil {
		r.runnable()
	}
}
