package d3d11

import (
	"fmt"
	"structs"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/kirides/go-d3d"
	"github.com/kirides/go-d3d/com"
	"github.com/kirides/go-d3d/dxgi"
)

var (
	modD3D11              = windows.NewLazySystemDLL("d3d11.dll")
	procD3D11CreateDevice = modD3D11.NewProc("D3D11CreateDevice")

	// iid_ID3D11Device1, _   = windows.GUIDFromString("{a04bfb29-08ef-43d6-a49c-a9bdbdcbe686}")
	IID_ID3D11Texture2D, _ = windows.GUIDFromString("{6f15aaf2-d208-4e89-9ab4-489535d34f9c}")
	IID_ID3D11Debug, _     = windows.GUIDFromString("{79cf2233-7536-4948-9d36-1e4692dc5760}")
	IID_ID3D11InfoQueue, _ = windows.GUIDFromString("{6543dbb6-1b48-42f5-ab82-e97ec74326f6}")
)

const (
	D3D11_CPU_ACCESS_READ = 0x20000

	D3D11_RLDO_SUMMARY         = 0x1
	D3D11_RLDO_DETAIL          = 0x2
	D3D11_RLDO_IGNORE_INTERNAL = 0x4

	D3D11_CREATE_DEVICE_DEBUG        = 0x2
	D3D11_CREATE_DEVICE_BGRA_SUPPORT = 0x20

	D3D11_SDK_VERSION = 7
)

func _D3D11CreateDevice(ppDevice **ID3D11Device, ppDeviceContext **ID3D11DeviceContext) error {
	var factory1 *dxgi.IDXGIFactory1
	if err := dxgi.CreateDXGIFactory1(&factory1); err != nil {
		return fmt.Errorf("CreateDXGIFactory1: %w", err)
	}
	defer factory1.Release()

	var adapter1 *dxgi.IDXGIAdapter1
	var desc dxgi.DXGI_ADAPTER_DESC1

	ai := uint32(0)
	for {
		hr := factory1.EnumAdapters1(ai, &adapter1)
		if d3d.HRESULT(hr).Failed() {
			break
		}
		ai++

		hr = int32(adapter1.GetDesc1(&desc))
		if d3d.HRESULT(hr).Failed() {
			adapter1.Release()
			adapter1 = nil
			continue
		}

		if (desc.Flags & dxgi.DXGI_ADAPTER_FLAG_SOFTWARE) == 0 {
			break
		}
		adapter1.Release()
		adapter1 = nil
	}

	if adapter1 == nil {
		hr := factory1.EnumAdapters1(0, &adapter1)
		if d3d.HRESULT(hr).Failed() {
			return fmt.Errorf("failed to fallback to default display adapter")
		}
	}

	if adapter1 == nil {
		return fmt.Errorf("no suitable adapter found")
	}

	defer adapter1.Release()

	fflags := [...]uint32{
		0xc100, // D3D_FEATURE_LEVEL_12_1
		0xc000, // D3D_FEATURE_LEVEL_12_0
		0xb100, // D3D_FEATURE_LEVEL_11_1
		0xb000, // D3D_FEATURE_LEVEL_11_0
		0xa100, // D3D_FEATURE_LEVEL_10_1
		0xa000, // D3D_FEATURE_LEVEL_10_0
		// 0x9300, // D3D_FEATURE_LEVEL_9_3
		// 0x9200, // D3D_FEATURE_LEVEL_9_2
		// 0x9100, // D3D_FEATURE_LEVEL_9_1
		// 0x1000, // D3D_FEATURE_LEVEL_1_0_CORE <-- unsupported!
	}
	featureLevel := 0x9100
	flags :=
		//  D3D11_CREATE_DEVICE_DEBUG |
		0

	const D3D_DRIVER_TYPE_UNKNOWN = 0
	ret, _, _ := syscall.SyscallN(
		procD3D11CreateDevice.Addr(),
		uintptr(unsafe.Pointer(adapter1)),   // pAdapter
		uintptr(D3D_DRIVER_TYPE_UNKNOWN),    // driverType: 1 = Hardware
		uintptr(0),                          // software
		uintptr(flags),                      // flags
		uintptr(unsafe.Pointer(&fflags[0])), // supported feature levels
		uintptr(len(fflags)),                // number of levels
		uintptr(D3D11_SDK_VERSION),
		uintptr(unsafe.Pointer(ppDevice)),        // *D3D11Device
		uintptr(unsafe.Pointer(&featureLevel)),   // feature level
		uintptr(unsafe.Pointer(ppDeviceContext)), // *D3D11DeviceContext
	)

	if ret != 0 {
		return d3d.HRESULT(ret)
	}
	return nil
}

func NewD3D11Device() (*ID3D11Device, *ID3D11DeviceContext, error) {
	var device *ID3D11Device
	var deviceCtx *ID3D11DeviceContext

	err := _D3D11CreateDevice(&device, &deviceCtx)

	if err != nil || device == nil || deviceCtx == nil {
		return nil, nil, err
	}

	return device, deviceCtx, nil
}

type ID3D11Texture2D struct {
	_    structs.HostLayout
	vtbl *ID3D11Texture2DVtbl
}

func (obj *ID3D11Texture2D) GetDesc(desc *D3D11_TEXTURE2D_DESC) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.GetDesc,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(desc)),
	)
	return int32(ret)
}
func (obj *ID3D11Texture2D) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}
func (obj *ID3D11Texture2D) QueryInterface(iid windows.GUID, pp interface{}) int32 {
	return com.ReflectQueryInterface(obj, obj.vtbl.QueryInterface, &iid, pp)
}

type ID3D11Device struct {
	_    structs.HostLayout
	vtbl *ID3D11DeviceVtbl
}

func (obj *ID3D11Device) QueryInterface(iid windows.GUID, pp interface{}) int32 {
	return com.ReflectQueryInterface(obj, obj.vtbl.QueryInterface, &iid, pp)
}

func (obj *ID3D11Device) CreateTexture2D(desc *D3D11_TEXTURE2D_DESC, ppTexture2D **ID3D11Texture2D) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.CreateTexture2D,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(desc)),
		uintptr(0),
		uintptr(unsafe.Pointer(ppTexture2D)),
	)
	return int32(ret)
}

func (obj *ID3D11Device) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

type ID3D11Device1 struct {
	_    structs.HostLayout
	vtbl *ID3D11DeviceVtbl
}

func (obj *ID3D11Device1) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

func (obj *ID3D11Device1) CreateTexture2D(desc *D3D11_TEXTURE2D_DESC, ppTexture2D **ID3D11Texture2D) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.CreateTexture2D,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(desc)),
		uintptr(0),
		uintptr(unsafe.Pointer(ppTexture2D)),
	)
	return int32(ret)
}

type ID3D11DeviceContext struct {
	_    structs.HostLayout
	vtbl *ID3D11DeviceContextVtbl
}

func (obj *ID3D11DeviceContext) CopyResourceDXGI(dst, src *dxgi.IDXGIResource) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.CopyResource,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(dst)),
		uintptr(unsafe.Pointer(src)),
	)
	return int32(ret)
}
func (obj *ID3D11DeviceContext) CopyResource2D(dst, src *ID3D11Texture2D) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.CopyResource,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(dst)),
		uintptr(unsafe.Pointer(src)),
	)
	return int32(ret)
}
func (obj *ID3D11DeviceContext) CopySubresourceRegion2D(dst *ID3D11Texture2D, dstSubResource, dstX, dstY, dstZ uint32, src *ID3D11Texture2D, srcSubResource uint32, pSrcBox *D3D11_BOX) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.CopySubresourceRegion,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(dst)),
		uintptr(dstSubResource),
		uintptr(dstX),
		uintptr(dstY),
		uintptr(dstZ),
		uintptr(unsafe.Pointer(src)),
		uintptr(srcSubResource),
		uintptr(unsafe.Pointer(pSrcBox)),
	)
	return int32(ret)
}

func (obj *ID3D11DeviceContext) CopySubresourceRegion(dst *ID3D11Resource, dstSubResource, dstX, dstY, dstZ uint32, src *ID3D11Resource, srcSubResource uint32, pSrcBox *D3D11_BOX) int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.CopySubresourceRegion,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(dst)),
		uintptr(dstSubResource),
		uintptr(dstX),
		uintptr(dstY),
		uintptr(dstZ),
		uintptr(unsafe.Pointer(src)),
		uintptr(srcSubResource),
		uintptr(unsafe.Pointer(pSrcBox)),
	)
	return int32(ret)
}
func (obj *ID3D11DeviceContext) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}

type ID3D11Resource struct {
	_    structs.HostLayout
	vtbl *ID3D11ResourceVtbl
}

func (obj *ID3D11Resource) Release() int32 {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return int32(ret)
}
