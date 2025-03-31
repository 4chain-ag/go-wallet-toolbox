package testabilities

import "testing"

func New(t testing.TB) (given FunderFixture, then FunderAssertion, cleanup func()) {
	given, cleanup = newFixture(t)
	fixture := given.(*funderFixture)
	then = newFunderAssertion(t, fixture)
	return
}
