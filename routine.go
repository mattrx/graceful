package graceful

import (
	"context"
	"fmt"
	"runtime"
	"sync"
)

type middlewareFunc func(next func(context.Context)) func(context.Context)
type logFunc func(message string)

var middlewares []middlewareFunc
var wg = sync.WaitGroup{}

var debugLogger logFunc

// EnableDebugLogging will print the file and line number for each started go routine
func EnableDebugLogging(logger logFunc) {
	debugLogger = logger
}

// DisableDebugLogging removes the logger
func DisableDebugLogging() {
	debugLogger = nil
}

// Wait for all go routines to finish
func Wait() {
	wg.Wait()
}

// Run a new tracked go routine
func Run(routine func()) {
	doRun(context.Background(), func(context.Context) {
		routine()
	})
}

// RunContext runs a new tracked go routine with a context
func RunContext(ctx context.Context, routine func(context.Context)) {
	doRun(ctx, routine)
}

func doRun(ctx context.Context, routine func(context.Context)) {
	wg.Add(1)

	if debugLogger != nil {
		if _, file, line, ok := runtime.Caller(2); ok {
			debugLogger(fmt.Sprintf("Go routine started in: %v:%v", file, line))
		}
	}

	go func() {
		for _, mw := range middlewares {
			routine = mw(routine)
		}

		routine(ctx)
		wg.Done()
	}()
}

// AddMiddleware to the stack
func AddMiddleware(mw middlewareFunc) {
	// append new middleware to the front to preserve order of execution
	middlewares = append([]middlewareFunc{mw}, middlewares...)
}

// ClearMiddlewares from the stack
func ClearMiddlewares() {
	middlewares = []middlewareFunc{}
}
