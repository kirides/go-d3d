package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"net/http"
	"runtime"

	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/kirides/go-d3d/d3d11"
	"github.com/kirides/go-d3d/examples/framelimiter"
	"github.com/kirides/go-d3d/outputduplication"
	"github.com/kirides/go-d3d/win"

	"github.com/kbinani/screenshot"
	"github.com/mattn/go-mjpeg"
)

func main() {
	n := screenshot.NumActiveDisplays()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	http.HandleFunc("/watch", func(w http.ResponseWriter, r *http.Request) {
		screen := r.URL.Query().Get("screen")
		if screen == "" {
			screen = "0"
		}
		screenNo, err := strconv.Atoi(screen)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		if screenNo >= n || screenNo < 0 {
			screenNo = 0
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title> Screen ` + strconv.Itoa(screenNo) + `</title>
	</head>
		<body style="margin:0">
	<img src="/mjpeg` + strconv.Itoa(screenNo) + `" style="max-width: 100vw; max-height: 100vh;object-fit: contain;display: block;margin: 0 auto;" />
</body>`))
	})

	framerate := 25
	for i := 0; i < screenshot.NumActiveDisplays(); i++ {
		fmt.Fprintf(os.Stderr, "Registering stream %d\n", i)
		stream := mjpeg.NewStream()
		defer stream.Close()
		go streamDisplayDXGI(ctx, i, framerate, stream)
		http.HandleFunc(fmt.Sprintf("/mjpeg%d", i), stream.ServeHTTP)
	}
	go func() {
		http.ListenAndServe("127.0.0.1:8023", nil)

	}()
	<-ctx.Done()
	<-time.After(time.Second)
}

// Capture using IDXGIOutputDuplication
//
//	https://docs.microsoft.com/en-us/windows/win32/api/dxgi1_2/nn-dxgi1_2-idxgioutputduplication
func streamDisplayDXGI(ctx context.Context, n int, framerate int, out *mjpeg.Stream) {
	max := screenshot.NumActiveDisplays()
	if n >= max {
		fmt.Printf("Not enough displays\n")
		return
	}

	// Keep this thread, so windows/d3d11/dxgi can use their threadlocal caches, if any
	runtime.LockOSThread()

	// Make thread PerMonitorV2 Dpi aware if supported on OS
	// allows to let windows handle BGRA -> RGBA conversion and possibly more things
	if win.IsValidDpiAwarenessContext(win.DpiAwarenessContextPerMonitorAwareV2) {
		_, err := win.SetThreadDpiAwarenessContext(win.DpiAwarenessContextPerMonitorAwareV2)
		if err != nil {
			fmt.Printf("Could not set thread DPI awareness to PerMonitorAwareV2. %v\n", err)
		} else {
			fmt.Printf("Enabled PerMonitorAwareV2 DPI awareness.\n")
		}
	}

	// Setup D3D11 stuff
	device, deviceCtx, err := d3d11.NewD3D11Device()
	if err != nil {
		fmt.Printf("Could not create D3D11 Device. %v\n", err)
		return
	}
	defer device.Release()
	defer deviceCtx.Release()

	var ddup *outputduplication.OutputDuplicator
	defer func() {
		if ddup != nil {
			ddup.Release()
			ddup = nil
		}
	}()

	buf := &bufferFlusher{Buffer: bytes.Buffer{}}
	opts := jpegQuality(50)
	limiter := framelimiter.New(framerate)

	lastBounds := image.Rectangle{}
	var imgBuf *image.RGBA
	for {
		select {
		case <-ctx.Done():
			return
		default:
			limiter.Wait()
		}
		// create output duplication if doesn't exist yet (maybe due to resolution change)
		if ddup == nil {
			ddup, err = outputduplication.NewIDXGIOutputDuplication(device, deviceCtx, uint(n))
			if err != nil {
				fmt.Printf("err: %v\n", err)
				continue
			}
			bounds, err := ddup.GetBounds()
			if err != nil {
				return
			}
			if bounds != lastBounds {
				lastBounds = bounds
				imgBuf = image.NewRGBA(lastBounds)
			}
		}

		// Grab an image.RGBA from the current output presenter
		err = ddup.GetImage(imgBuf, 999)
		if err != nil {
			if errors.Is(err, outputduplication.ErrNoImageYet) {
				// don't update
				continue
			}
			fmt.Printf("Err ddup.GetImage: %v\n", err)
			// Retry with new ddup, can occur when changing resolution
			ddup.Release()
			ddup = nil
			continue
		}
		buf.Reset()
		encodeJpeg(buf, imgBuf, opts)
		out.Update(buf.Bytes())
	}
}

// Workaround for jpeg.Encode(), which requires a Flush()
// method to not call `bufio.NewWriter`
type bufferFlusher struct {
	bytes.Buffer
}

func (*bufferFlusher) Flush() error { return nil }
