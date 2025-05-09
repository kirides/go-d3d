package dxgi

import (
	"structs"
	"syscall"
	"unsafe"

	"github.com/kirides/go-d3d"
	"github.com/kirides/go-d3d/com"
	"golang.org/x/sys/windows"
)

var (
	modDXGI                = windows.NewLazySystemDLL("dxgi.dll")
	procCreateDXGIFactory1 = modDXGI.NewProc("CreateDXGIFactory1")

	// iid_IDXGIDevice, _   = windows.GUIDFromString("{54ec77fa-1377-44e6-8c32-88fd5f44c84c}")
	IID_IDXGIDevice1, _ = windows.GUIDFromString("{77db970f-6276-48ba-ba28-070143b4392c}")
	// IID_IDXGIAdapter, _  = windows.GUIDFromString("{2411E7E1-12AC-4CCF-BD14-9798E8534DC0}")
	IID_IDXGIAdapter1, _ = windows.GUIDFromString("{29038f61-3839-4626-91fd-086879011a05}")
	// IID_IDXGIOutput, _   = windows.GUIDFromString("{ae02eedb-c735-4690-8d52-5a8dc20213aa}")
	IID_IDXGIOutput1, _  = windows.GUIDFromString("{00cddea8-939b-4b83-a340-a685226666cc}")
	IID_IDXGIOutput5, _  = windows.GUIDFromString("{80A07424-AB52-42EB-833C-0C42FD282D98}")
	IID_IDXGIFactory1, _ = windows.GUIDFromString("{770aae78-f26f-4dba-a829-253c83d1b387}")
	// IID_IDXGIResource, _ = windows.GUIDFromString("{035f3ab4-482e-4e50-b41f-8a7f8bd8960b}")
	IID_IDXGISurface, _ = windows.GUIDFromString("{cafcb56c-6ac3-4889-bf47-9e23bbd260ec}")
)

const (
	DXGI_MAP_READ    = 1 << 0
	DXGI_MAP_WRITE   = 1 << 1
	DXGI_MAP_DISCARD = 1 << 2
)

type IDXGIFactory1 struct {
	_    structs.HostLayout
	vtbl *IDXGIFactory1Vtbl
}

func (obj *IDXGIFactory1) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

func (obj *IDXGIFactory1) EnumAdapters1(adapter uint32, pp **IDXGIAdapter1) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.EnumAdapters1,
		uintptr(unsafe.Pointer(obj)),
		uintptr(adapter),
		uintptr(unsafe.Pointer(pp)),
	)
	return int32(ret)
}

func CreateDXGIFactory1(ppFactory **IDXGIFactory1) error {
	ret, _, _ := syscall.SyscallN(
		procCreateDXGIFactory1.Addr(),
		uintptr(unsafe.Pointer(&IID_IDXGIFactory1)),
		uintptr(unsafe.Pointer(ppFactory)),
	)
	if ret != 0 {
		return d3d.HRESULT(ret)
	}

	return nil
}

type IDXGIAdapter1 struct {
	_    structs.HostLayout
	vtbl *IDXGIAdapter1Vtbl
}

func (obj *IDXGIAdapter1) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

func (obj *IDXGIAdapter1) EnumOutputs(output uint32, pp **IDXGIOutput) uint32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.EnumOutputs,
		uintptr(unsafe.Pointer(obj)),
		uintptr(output),
		uintptr(unsafe.Pointer(pp)),
	)
	return uint32(ret)
}

type IDXGIAdapter struct {
	_    structs.HostLayout
	vtbl *IDXGIAdapterVtbl
}

func (obj *IDXGIAdapter) EnumOutputs(output uint32, pp **IDXGIOutput) uint32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.EnumOutputs,
		uintptr(unsafe.Pointer(obj)),
		uintptr(output),
		uintptr(unsafe.Pointer(pp)),
	)
	return uint32(ret)
}

func (obj *IDXGIAdapter) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

type IDXGIDevice struct {
	_    structs.HostLayout
	vtbl *IDXGIDeviceVtbl
}

func (obj *IDXGIDevice) GetGPUThreadPriority(priority *int) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetGPUThreadPriority,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(priority)),
	)
	return int32(ret)
}
func (obj *IDXGIDevice) QueryInterface(iid windows.GUID, pp interface{}) int32 {
	return com.ReflectQueryInterface(obj, obj.vtbl.QueryInterface, &iid, pp)
}
func (obj *IDXGIDevice) GetParent(iid windows.GUID, pp *unsafe.Pointer) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetParent,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&iid)),
		uintptr(unsafe.Pointer(pp)),
	)
	return int32(ret)
}
func (obj *IDXGIDevice) GetAdapter(pAdapter **IDXGIAdapter) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetAdapter,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(pAdapter)),
	)
	return int32(ret)
}
func (obj *IDXGIDevice) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

type IDXGIDevice1 struct {
	_    structs.HostLayout
	vtbl *IDXGIDevice1Vtbl
}

func (obj *IDXGIDevice1) QueryInterface(iid windows.GUID, pp interface{}) int32 {
	return com.ReflectQueryInterface(obj, obj.vtbl.QueryInterface, &iid, pp)
}

func (obj *IDXGIDevice1) GetParent(iid windows.GUID, pp *unsafe.Pointer) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetParent,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&iid)),
		uintptr(unsafe.Pointer(pp)),
	)

	return int32(ret)
}
func (obj *IDXGIDevice1) GetAdapter(pAdapter *IDXGIAdapter) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetAdapter,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&pAdapter)),
	)

	return int32(ret)
}
func (obj *IDXGIDevice1) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

type IDXGIOutput struct {
	_    structs.HostLayout
	vtbl *IDXGIOutputVtbl
}

func (obj *IDXGIOutput) QueryInterface(iid windows.GUID, pp interface{}) int32 {
	return com.ReflectQueryInterface(obj, obj.vtbl.QueryInterface, &iid, pp)
}

func (obj *IDXGIOutput) GetParent(iid windows.GUID, pp *unsafe.Pointer) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetParent,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&iid)),
		uintptr(unsafe.Pointer(pp)),
	)
	return int32(ret)
}

func (obj *IDXGIOutput) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

type IDXGIOutput1 struct {
	_    structs.HostLayout
	vtbl *IDXGIOutput1Vtbl
}

func (obj *IDXGIOutput1) DuplicateOutput(device1 *IDXGIDevice1, ppOutputDuplication **IDXGIOutputDuplication) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.DuplicateOutput,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(device1)),
		uintptr(unsafe.Pointer(ppOutputDuplication)),
	)
	return int32(ret)
}

func (obj *IDXGIOutput1) GetParent(iid windows.GUID, pp *unsafe.Pointer) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetParent,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&iid)),
		uintptr(unsafe.Pointer(pp)),
	)
	return int32(ret)
}

func (obj *IDXGIOutput1) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

type IDXGIOutput5 struct {
	_    structs.HostLayout
	vtbl *IDXGIOutput5Vtbl
}

type DXGI_FORMAT uint32

func (obj *IDXGIOutput5) GetDesc(desc *DXGI_OUTPUT_DESC) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetDesc,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(desc)),
	)
	return int32(ret)
}

func (obj *IDXGIOutput5) DuplicateOutput1(device1 *IDXGIDevice1, flags uint32, pSupportedFormats []DXGI_FORMAT, ppOutputDuplication **IDXGIOutputDuplication) int32 {
	pFormats := &pSupportedFormats[0]
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.DuplicateOutput1,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(device1)),
		uintptr(flags),
		uintptr(len(pSupportedFormats)),
		uintptr(unsafe.Pointer(pFormats)),
		uintptr(unsafe.Pointer(ppOutputDuplication)),
	)
	return int32(ret)
}

func (obj *IDXGIOutput5) GetParent(iid windows.GUID, pp *unsafe.Pointer) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetParent,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&iid)),
		uintptr(unsafe.Pointer(pp)),
	)
	return int32(ret)
}

func (obj *IDXGIOutput5) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

type IDXGIResource struct {
	_    structs.HostLayout
	vtbl *IDXGIResourceVtbl
}

func (obj *IDXGIResource) QueryInterface(iid windows.GUID, pp interface{}) int32 {
	return com.ReflectQueryInterface(obj, obj.vtbl.QueryInterface, &iid, pp)
}
func (obj *IDXGIResource) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

type IDXGISurface struct {
	_    structs.HostLayout
	vtbl *IDXGISurfaceVtbl
}

func (obj *IDXGISurface) QueryInterface(iid windows.GUID, pp interface{}) int32 {
	return com.ReflectQueryInterface(obj, obj.vtbl.QueryInterface, &iid, pp)
}
func (obj *IDXGISurface) Map(pLockedRect *DXGI_MAPPED_RECT, mapFlags uint32) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Map,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(pLockedRect)),
		uintptr(mapFlags),
	)
	return int32(ret)
}
func (obj *IDXGISurface) Unmap() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Unmap,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}
func (obj *IDXGISurface) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

type IDXGIOutputDuplication struct {
	_    structs.HostLayout
	vtbl *IDXGIOutputDuplicationVtbl
}

func (obj *IDXGIOutputDuplication) GetFrameMoveRects(buffer []DXGI_OUTDUPL_MOVE_RECT, rectsRequired *uint32) int32 {
	var buf *DXGI_OUTDUPL_MOVE_RECT
	if len(buffer) > 0 {
		buf = &buffer[0]
	}
	size := uint32(len(buffer) * 24)
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetFrameMoveRects,
		uintptr(unsafe.Pointer(obj)),
		uintptr(size),
		uintptr(unsafe.Pointer(buf)),
		uintptr(unsafe.Pointer(rectsRequired)),
	)
	*rectsRequired = *rectsRequired / 24
	return int32(ret)
}
func (obj *IDXGIOutputDuplication) GetFrameDirtyRects(buffer []RECT, rectsRequired *uint32) int32 {
	var buf *RECT
	if len(buffer) > 0 {
		buf = &buffer[0]
	}
	size := uint32(len(buffer) * 16)
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetFrameDirtyRects,
		uintptr(unsafe.Pointer(obj)),
		uintptr(size),
		uintptr(unsafe.Pointer(buf)),
		uintptr(unsafe.Pointer(rectsRequired)),
	)
	*rectsRequired = *rectsRequired / 16
	return int32(ret)
}

func (obj *IDXGIOutputDuplication) GetFramePointerShape(pointerShapeBufferSize uint32,
	pPointerShapeBuffer []byte,
	pPointerShapeBufferSizeRequired *uint32,
	pPointerShapeInfo *DXGI_OUTDUPL_POINTER_SHAPE_INFO) int32 {

	var buf *byte
	if len(pPointerShapeBuffer) > 0 {
		buf = &pPointerShapeBuffer[0]
	}

	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetFramePointerShape,
		uintptr(unsafe.Pointer(obj)),
		uintptr(pointerShapeBufferSize),
		uintptr(unsafe.Pointer(buf)),
		uintptr(unsafe.Pointer(pPointerShapeBufferSizeRequired)),
		uintptr(unsafe.Pointer(pPointerShapeInfo)),
	)

	return int32(ret)
}
func (obj *IDXGIOutputDuplication) GetDesc(desc *DXGI_OUTDUPL_DESC) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetDesc,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(desc)),
	)
	return int32(ret)
}

func (obj *IDXGIOutputDuplication) MapDesktopSurface(pLockedRect *DXGI_MAPPED_RECT) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.MapDesktopSurface,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(pLockedRect)),
	)
	return int32(ret)
}
func (obj *IDXGIOutputDuplication) UnMapDesktopSurface() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.UnMapDesktopSurface,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}
func (obj *IDXGIOutputDuplication) AddRef() uint32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.AddRef,
		uintptr(unsafe.Pointer(obj)),
	)
	return uint32(ret)
}

func (obj *IDXGIOutputDuplication) Release() uint32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return uint32(ret)
}

func (obj *IDXGIOutputDuplication) AcquireNextFrame(timeoutMs uint32, pFrameInfo *DXGI_OUTDUPL_FRAME_INFO, ppDesktopResource **IDXGIResource) uint32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.AcquireNextFrame,    // function address
		uintptr(unsafe.Pointer(obj)), // always pass the COM object address first
		uintptr(timeoutMs),           // then all function parameters follow
		uintptr(unsafe.Pointer(pFrameInfo)),
		uintptr(unsafe.Pointer(ppDesktopResource)),
	)
	return uint32(ret)
}

func (obj *IDXGIOutputDuplication) ReleaseFrame() uint32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.ReleaseFrame,
		uintptr(unsafe.Pointer(obj)),
	)
	return uint32(ret)
}
