package graceful

import (
	"net/http"
	"syscall"
	"testing"
	"time"
)

func TestListenAndServe(t *testing.T) {

	server := &http.Server{}
	finished := false

	serverChan := make(chan struct{})

	go func() {
		if err := ListenAndServe(server); err != nil && err != http.ErrServerClosed {
			t.Fail()
		}

		finished = true
		serverChan <- struct{}{}
	}()

	time.Sleep(time.Second * 2)

	signalChan <- syscall.SIGINT

	<-serverChan

	if finished != true {
		t.Fail()
	}
}
