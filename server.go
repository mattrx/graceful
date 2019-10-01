package graceful

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// ListenAndServe on the server with graceful shutdown
func ListenAndServe(server *http.Server) error {

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGTERM)
	signal.Notify(signalChan, syscall.SIGINT)

	go func() {
		sig := <-signalChan
		log.Printf("Signal received: %+v\n", sig)

		log.Println("Shutting down server...")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Println("Could not shut down server: ", err)
		}
	}()

	return server.ListenAndServe()
}
