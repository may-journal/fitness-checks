// Package req ties a Go test case to the requirement THEN it proves, by the
// THEN's ID and a hash of its text. A case that makes no assertion fails.
package req

import (
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/may-journal/fitness-checks/req/digest"
)

// Hash is the tag hash of a THEN, as digest.Of gives it.
func Hash(text string) string { return digest.Of(text) }

// running maps each test inside a Then to its assertion count.
var running sync.Map

// Asserted records that an assertion ran for t. Assertion helpers call it,
// and a Then body that records none fails.
func Asserted(t testing.TB) {
	if n, ok := running.Load(t); ok {
		n.(*atomic.Int64).Add(1)
	}
}

// A is what a Then body asserts with. Each method counts as an assertion.
type A struct{ t *testing.T }

// T is the test running the case, for helpers that take one. Helpers that
// assert call Asserted with it.
func (a A) T() *testing.T { return a.t }

// Equal fails the case unless got equals want.
func (a A) Equal(want, got any) {
	a.t.Helper()
	Asserted(a.t)
	if !reflect.DeepEqual(want, got) {
		a.t.Errorf("got %v, want %v", got, want)
	}
}

// True fails the case with the message unless ok.
func (a A) True(ok bool, format string, args ...any) {
	a.t.Helper()
	Asserted(a.t)
	if !ok {
		a.t.Errorf(format, args...)
	}
}

// Contains fails the case unless text holds want.
func (a A) Contains(text, want string) {
	a.t.Helper()
	Asserted(a.t)
	if !strings.Contains(text, want) {
		a.t.Errorf("%q is missing %q", text, want)
	}
}

// NoError fails the case if err is not nil.
func (a A) NoError(err error) {
	a.t.Helper()
	Asserted(a.t)
	if err != nil {
		a.t.Errorf("unexpected error: %v", err)
	}
}

// Body is what a Then runs: a function of the case's own test, or of the A
// that wraps it.
type Body interface{ func(t *testing.T) | func(a A) }

// A call of Then tags body with the ID and hash of its THEN and runs it in a
// nested case of t. A body of func(t *testing.T) gets that case's test, so
// code written for a plain test runs as it was. The nested case fails when
// body makes no assertion, and neither fails nor skips.
func Then[B Body](t *testing.T, id, hash string, body B) {
	t.Helper()
	t.Run("then", func(t *testing.T) {
		count := new(atomic.Int64)
		running.Store(t, count)
		defer running.Delete(t)
		defer func() {
			if count.Load() == 0 && !t.Failed() && !t.Skipped() {
				t.Errorf("THEN %s#%s makes no assertion; assert its result in the case", id, hash)
			}
		}()
		switch run := any(body).(type) {
		case func(t *testing.T):
			run(t)
		case func(a A):
			run(A{t: t})
		}
	})
}
