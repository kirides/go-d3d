package win

//go:generate mkwinsyscall -output zsyscall_windows.go syscall_windows.go

type (
	BOOL          uint32
	BOOLEAN       byte
	BYTE          byte
	DWORD         uint32
	DWORD64       uint64
	HANDLE        uintptr
	HLOCAL        uintptr
	LARGE_INTEGER int64
	LONG          int32
	LPVOID        uintptr
	SIZE_T        uintptr
	UINT          uint32
	ULONG_PTR     uintptr
	ULONGLONG     uint64
	WORD          uint16

	HWND uintptr
)

type BITMAPINFOHEADER struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}
type RGBQUAD struct {
	RgbBlue     byte
	RgbGreen    byte
	RgbRed      byte
	RgbReserved byte
}

type BITMAPINFO struct {
	BmiHeader BITMAPINFOHEADER
	BmiColors *RGBQUAD
}

const (
	OBJ_BITMAP = 7
)

const (
	DpiAwarenessContextUndefined         = 0
	DpiAwarenessContextUnaware           = -1
	DpiAwarenessContextSystemAware       = -2
	DpiAwarenessContextPerMonitorAware   = -3
	DpiAwarenessContextPerMonitorAwareV2 = -4
	DpiAwarenessContextUnawareGdiScaled  = -5
)

//sys	SetThreadDpiAwarenessContext(value int32) (n int, err error) = User32.SetThreadDpiAwarenessContext
//sys	IsValidDpiAwarenessContext(value int32) (n bool) = User32.IsValidDpiAwarenessContext
