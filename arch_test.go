package arch_test

import (
	"github.com/matthewmcnew/archtest"
	"testing"
)

func TestArch_NoDependencies(t *testing.T) {

	t.Run("Services do not depend on controllers", func(t *testing.T) {
		mockT := new(testingT)
		archtest.Package(mockT, "gses4_project/internal/services").
			ShouldNotDependOn("gses4_project/internal/server/controllers")
		assertNoError(t, mockT)
	})

	t.Run("Services do not depend on database", func(t *testing.T) {
		mockT := new(testingT)
		archtest.Package(mockT, "gses4_project/internal/services").ShouldNotDependOn(
			"gses4_project/internal/database")
		assertNoError(t, mockT)
	})

	t.Run("Services do not depend on providers", func(t *testing.T) {
		mockT := new(testingT)
		archtest.Package(mockT, "gses4_project/internal/services").ShouldNotDependOn(
			"gses4_project/internal/providers")
		assertNoError(t, mockT)
	})

	t.Run("Services do not depend on mailers", func(t *testing.T) {
		mockT := new(testingT)
		archtest.Package(mockT, "gses4_project/internal/services").ShouldNotDependOn(
			"gses4_project/internal/mailers")
		assertNoError(t, mockT)
	})

	t.Run("Services do not depend on crons", func(t *testing.T) {
		mockT := new(testingT)
		archtest.Package(mockT, "gses4_project/internal/services").ShouldNotDependOn(
			"gses4_project/internal/crons")
		assertNoError(t, mockT)
	})

	t.Run("Providers do not depend on mailers", func(t *testing.T) {
		mockT := new(testingT)
		archtest.Package(mockT, "gses4_project/internal/providers").
			ShouldNotDependOn("gses4_project/internal/providers")
		assertNoError(t, mockT)
	})

	t.Run("Providers do not depend on database", func(t *testing.T) {
		mockT := new(testingT)
		archtest.Package(mockT, "gses4_project/internal/providers").
			ShouldNotDependOn("gses4_project/internal/database")
		assertNoError(t, mockT)
	})

	t.Run("Providers do not depend on crons", func(t *testing.T) {
		mockT := new(testingT)
		archtest.Package(mockT, "gses4_project/internal/providers").
			ShouldNotDependOn("gses4_project/internal/crons")
		assertNoError(t, mockT)
	})

	t.Run("Crons do not depend on providers", func(t *testing.T) {
		mockT := new(testingT)
		archtest.Package(mockT, "gses4_project/internal/crons").
			ShouldNotDependOn("gses4_project/internal/providers")
		assertNoError(t, mockT)
	})

	t.Run("Crons do not depend on database", func(t *testing.T) {
		mockT := new(testingT)
		archtest.Package(mockT, "gses4_project/internal/crons").
			ShouldNotDependOn("gses4_project/internal/database")
		assertNoError(t, mockT)
	})

	t.Run("Crons do not depend on mailers", func(t *testing.T) {
		mockT := new(testingT)
		archtest.Package(mockT, "gses4_project/internal/crons").
			ShouldNotDependOn("gses4_project/internal/mailers")
		assertNoError(t, mockT)
	})

	t.Run("Mailers do not depend on providers", func(t *testing.T) {
		mockT := new(testingT)
		archtest.Package(mockT, "gses4_project/internal/mailers").
			ShouldNotDependOn("gses4_project/internal/providers")
		assertNoError(t, mockT)
	})

	t.Run("Mailers do not depend on database", func(t *testing.T) {
		mockT := new(testingT)
		archtest.Package(mockT, "gses4_project/internal/mailers").
			ShouldNotDependOn("gses4_project/internal/database")
		assertNoError(t, mockT)
	})

	t.Run("Mailers do not depend on crons", func(t *testing.T) {
		mockT := new(testingT)
		archtest.Package(mockT, "gses4_project/internal/mailers").
			ShouldNotDependOn("gses4_project/internal/crons")
		assertNoError(t, mockT)
	})

	t.Run("Controllers do not depend on providers", func(t *testing.T) {
		mockT := new(testingT)
		archtest.Package(mockT, "gses4_project/internal/server/controllers").
			ShouldNotDependOn("gses4_project/internal/providers")
		assertNoError(t, mockT)
	})

	t.Run("Controllers do not depend on database", func(t *testing.T) {
		mockT := new(testingT)
		archtest.Package(mockT, "gses4_project/internal/server/controllers").
			ShouldNotDependOn("gses4_project/internal/database")
		assertNoError(t, mockT)
	})

	t.Run("Controllers do not depend on crons", func(t *testing.T) {
		mockT := new(testingT)
		archtest.Package(mockT, "gses4_project/internal/server/controllers").
			ShouldNotDependOn("gses4_project/internal/crons")
		assertNoError(t, mockT)
	})
}

func assertNoError(t *testing.T, mockT *testingT) {
	t.Helper()
	if mockT.errored() {
		t.Fatalf("archtest should not have failed but, %s", mockT.message())
	}
}

type testingT struct {
	errors [][]interface{}
}

func (t *testingT) Error(args ...interface{}) {
	t.errors = append(t.errors, args)
}

func (t *testingT) errored() bool {
	return len(t.errors) != 0
}

func (t *testingT) message() interface{} {
	return t.errors[0][0]
}
