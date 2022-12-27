package main

import (
	"context"
	"fmt"
	"net/http"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kbinani/screenshot"
	"github.com/mattn/go-mjpeg"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	framerate := 25
	for i := 0; i < screenshot.NumActiveDisplays(); i++ {
		fmt.Fprintf(os.Stderr, "Registering stream %d\n", i)
		stream := mjpeg.NewStream()
		defer stream.Close()
		go captureScreenTranscode(ctx, i, framerate)
		http.HandleFunc(fmt.Sprintf("/mjpeg%d", i), stream.ServeHTTP)
	}
	<-ctx.Done()
	<-time.After(time.Second)
}
