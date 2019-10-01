# graceful

- [Graceful Shutdown for your HTTP Server](#graceful-http-server)
- [Graceful Shutdown for your Go Routines](#graceful-go-routines)
- [Middlewares for your Go Routines](#middlewares-go-routines)
- [Cleanup after Shutdown](#cleanup-funcs)

## <a name="graceful-http-server"></a>Graceful Shutdown for your HTTP Server

This package provides a method to gracefully shutdown your golang http server.

When a SIGTERM signal is received, the shutdown method on your server is executed. That means that no new connections are accepted but all running requests continue until finished.

```golang
func main() {
    server := &http.Server{}

    if err := graceful.ListenAndServe(server); err != nil && err != http.ErrServerClosed {
        log.Println("Could not start server, error:", err.Error())
        os.Exit(1)
    }
}
```

## <a name="graceful-go-routines"></a>Graceful Shutdown for your Go Routines

After shutting down your server you can wait for registered go routines to finish before exiting.

```golang
func main() {

    ...

    log.Println("Waiting for unfinished go routines...")
    graceful.Wait() // blocks until all go routines are finished
}
```

To keep track of your started go routines you need to start them like this (otherwise you can't wait for them):

```golang
graceful.Run(func() {
    // do not start a go routine on your own
    // the go routine will be started in the package to keep track of it
})
```

## <a name="middlewares-go-routines"></a>Middlewares for your Go Routines

It is possible to register middlewares that are applied automatically to all your go routines.

```golang
graceful.AddMiddleware(func(next func(context.Context)) func(context.Context) {
    return func(ctx context.Context) {
        // do something before your go routine
        next(ctx)
        // do something after your go routine
    }
})
```

For middlewares to be useful you should provide a context that is passed down to your function. The context is also available in all middlewares.

```golang
graceful.RunContext(context.Background(), func(ctx context.Context) {
    // do not start a go routine on your own
    // the go routine will be started in the package to keep track of it
})
```

#### Example: Middleware to trace Go Routines with OpenCensus

```golang
import (
	"github.com/mattrx/graceful"
	"go.opencensus.io/trace"
)

graceful.AddMiddleware(func(next func(context.Context)) func(context.Context) {
    return func(ctx context.Context) {
        ctx, span := trace.StartSpan(ctx, "go-routine")
        defer span.End()

        next(ctx)
    }
})
```

## <a name="cleanup-funcs"></a>Cleanup after Shutdown

After shutting down you may want to run some cleanup functions to close connections or flush some internal memory buffer. You can run them like this:

```golang
func main() {

    ...

    log.Println("Running cleanup functions...")
    graceful.Cleanup()
}
```

In your code you just have to register the functions when you get them:

```golang
db, err := sqlx.Connect("mysql", config.FormatDSN())

...

graceful.AddCleanup(func() {
    if err := db.Close(); err != nil {
        log.Printf("Error closing db conn: %v", err)
    }
})
```
