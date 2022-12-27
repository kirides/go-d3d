This code allows to use D3D11 IDXGIOutputDuplication in Go
## Examples

- Encoding an mjpeg stream [examples/mjpegstream](./examples/mjpegstream)
- Recording an h264 video using ffmpeg for transcoding [examples/recording](./examples/recording)

## Libaries used

- `golang.org/x/exp/shiny/driver/internal/swizzle` for faster BGRA -> RGBA conversion (see [shiny LICENSE](./outputduplication/swizzle/LICENSE))

## app.manifest

To make use of `IDXGIOutput5::DuplicateOutput1`, an application has to provide support for PerMonitorV2 DPI-Awareness (Windows 10 1703+) This is usually done by providing an `my-executable.exe.manifest` file either next to the executable, or as an embedded resource.

In the examples there are calls to `IsValidDpiAwarenessContext` and `SetThreadDpiAwarenessContext` which circumvent the requirement.
