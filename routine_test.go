package graceful_test

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	graceful "github.com/mattrx/graceful"
)

func TestRun(t *testing.T) {

	routineExecuted := false

	graceful.Run(func() {
		time.Sleep(time.Second * 1)
		routineExecuted = true
	})

	graceful.Wait()

	if routineExecuted == false {
		t.Fatalf("Expected routine to have been executed")
	}
}

func TestRunContext(t *testing.T) {

	routineExecuted := false
	expectedCtx := context.Background()

	graceful.RunContext(expectedCtx, func(ctx context.Context) {
		if !reflect.DeepEqual(expectedCtx, ctx) {
			t.Fatalf("Contexts does not match")
		}

		time.Sleep(time.Second * 1)
		routineExecuted = true
	})

	graceful.Wait()

	if routineExecuted == false {
		t.Fatalf("Expected routine to have been executed")
	}
}

func TestRunWithDebug(t *testing.T) {

	logFuncCalled := 0
	expectedLogMessagePart := "routine_test.go:63" // must be adjusted on changes to this file

	graceful.EnableDebugLogging(func(msg string) {
		logFuncCalled++

		if !strings.Contains(msg, expectedLogMessagePart) {
			t.Fatalf("Expected log message to contain '%v', got '%v'", expectedLogMessagePart, msg)
		}
	})

	graceful.Run(func() {})
	graceful.Wait()
	graceful.DisableDebugLogging()

	if logFuncCalled != 1 {
		t.Fatalf("Expected log func to have been called 1 time, called %v times", logFuncCalled)
	}
}

func TestRunContextWithDebug(t *testing.T) {

	logFuncCalled := 0
	expectedLogMessagePart := "routine_test.go:85" // must be adjusted on changes to this file

	graceful.EnableDebugLogging(func(msg string) {
		logFuncCalled++

		if !strings.Contains(msg, expectedLogMessagePart) {
			t.Fatalf("Expected log message to contain '%v', got '%v'", expectedLogMessagePart, msg)
		}
	})

	graceful.RunContext(context.Background(), func(ctx context.Context) {})
	graceful.Wait()
	graceful.DisableDebugLogging()

	if logFuncCalled != 1 {
		t.Fatalf("Expected log func to have been called 1 time, called %v times", logFuncCalled)
	}
}

func TestRunContextWithInlineMiddlewares(t *testing.T) {

	step := 0

	assertStep := func(s int) {
		if step != s {
			t.Fatalf("Exptected step %v, got %v", s, step)
		}
		step++
	}

	mw1 := func(next func(context.Context)) func(context.Context) {
		return func(ctx context.Context) {
			assertStep(0)
			next(ctx)
			assertStep(4)
		}
	}

	mw2 := func(next func(context.Context)) func(context.Context) {
		return func(ctx context.Context) {
			assertStep(1)
			next(ctx)
			assertStep(3)
		}
	}

	graceful.RunContext(context.Background(), mw1(mw2(func(ctx context.Context) {
		assertStep(2)
	})))

	graceful.Wait()

	if step == 0 {
		t.Fatalf("Expected step to not be 0")
	}
}

func TestRunContextWithRegisteredMiddlewares(t *testing.T) {

	step := 0

	assertStep := func(s int) {
		if step != s {
			t.Fatalf("Exptected step %v, got %v", s, step)
		}
		step++
	}

	graceful.AddMiddleware(func(next func(context.Context)) func(context.Context) {
		return func(ctx context.Context) {
			assertStep(0)
			next(ctx)
			assertStep(4)
		}
	})

	graceful.AddMiddleware(func(next func(context.Context)) func(context.Context) {
		return func(ctx context.Context) {
			assertStep(1)
			next(ctx)
			assertStep(3)
		}
	})

	graceful.RunContext(context.Background(), func(ctx context.Context) {
		assertStep(2)
	})

	graceful.Wait()
	graceful.ClearMiddlewares()
}
