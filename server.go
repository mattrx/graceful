package graceful

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var signalChan = make(chan os.Signal)

// ListenAndServe on the server with graceful shutdown
func ListenAndServe(server *http.Server) error {

	signal.Notify(signalChan, syscall.SIGTERM)
	signal.Notify(signalChan, syscall.SIGINT)

	shutdownChan := make(chan struct{})

	go func() {
		<-signalChan

		if err := server.Shutdown(context.Background()); err != nil {
			panic(fmt.Errorf("Could not shut down server: %w", err))
		}

		shutdownChan <- struct{}{}
	}()

	err := server.ListenAndServe()

	<-shutdownChan

	return err
}
