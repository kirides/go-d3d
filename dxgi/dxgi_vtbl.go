package dxgi

import (
	"structs"

	"github.com/kirides/go-d3d/com"
)

type IDXGIObjectVtbl struct {
	_ structs.HostLayout
	com.IUnknownVtbl

	SetPrivateData          uintptr
	SetPrivateDataInterface uintptr
	GetPrivateData          uintptr
	GetParent               uintptr
}

type IDXGIAdapterVtbl struct {
	_ structs.HostLayout
	IDXGIObjectVtbl

	EnumOutputs           uintptr
	GetDesc               uintptr
	CheckInterfaceSupport uintptr
}
type IDXGIAdapter1Vtbl struct {
	_ structs.HostLayout
	IDXGIAdapterVtbl

	GetDesc1 uintptr
}

type IDXGIDeviceVtbl struct {
	_ structs.HostLayout
	IDXGIObjectVtbl

	CreateSurface          uintptr
	GetAdapter             uintptr
	GetGPUThreadPriority   uintptr
	QueryResourceResidency uintptr
	SetGPUThreadPriority   uintptr
}

type IDXGIDevice1Vtbl struct {
	_ structs.HostLayout
	IDXGIDeviceVtbl

	GetMaximumFrameLatency uintptr
	SetMaximumFrameLatency uintptr
}

type IDXGIDeviceSubObjectVtbl struct {
	_ structs.HostLayout
	IDXGIObjectVtbl

	GetDevice uintptr
}

type IDXGISurfaceVtbl struct {
	_ structs.HostLayout
	IDXGIDeviceSubObjectVtbl

	GetDesc uintptr
	Map     uintptr
	Unmap   uintptr
}

type IDXGIResourceVtbl struct {
	_ structs.HostLayout
	IDXGIDeviceSubObjectVtbl

	GetSharedHandle     uintptr
	GetUsage            uintptr
	SetEvictionPriority uintptr
	GetEvictionPriority uintptr
}

type IDXGIOutputVtbl struct {
	_ structs.HostLayout
	IDXGIObjectVtbl

	GetDesc                     uintptr
	GetDisplayModeList          uintptr
	FindClosestMatchingMode     uintptr
	WaitForVBlank               uintptr
	TakeOwnership               uintptr
	ReleaseOwnership            uintptr
	GetGammaControlCapabilities uintptr
	SetGammaControl             uintptr
	GetGammaControl             uintptr
	SetDisplaySurface           uintptr
	GetDisplaySurfaceData       uintptr
	GetFrameStatistics          uintptr
}

type IDXGIOutput1Vtbl struct {
	_ structs.HostLayout
	IDXGIOutputVtbl

	GetDisplayModeList1      uintptr
	FindClosestMatchingMode1 uintptr
	GetDisplaySurfaceData1   uintptr
	DuplicateOutput          uintptr
}

type IDXGIOutput2Vtbl struct {
	_ structs.HostLayout
	IDXGIOutput1Vtbl

	SupportsOverlays uintptr
}

type IDXGIOutput3Vtbl struct {
	_ structs.HostLayout
	IDXGIOutput2Vtbl

	CheckOverlaySupport uintptr
}

type IDXGIOutput4Vtbl struct {
	_ structs.HostLayout
	IDXGIOutput3Vtbl

	CheckOverlayColorSpaceSupport uintptr
}
type IDXGIOutput5Vtbl struct {
	_ structs.HostLayout
	IDXGIOutput4Vtbl

	DuplicateOutput1 uintptr
}

type IDXGIOutputDuplicationVtbl struct {
	_ structs.HostLayout
	IDXGIObjectVtbl

	GetDesc              uintptr
	AcquireNextFrame     uintptr
	GetFrameDirtyRects   uintptr
	GetFrameMoveRects    uintptr
	GetFramePointerShape uintptr
	MapDesktopSurface    uintptr
	UnMapDesktopSurface  uintptr
	ReleaseFrame         uintptr
}
type IDXGIFactoryVtbl struct {
	_ structs.HostLayout
	IDXGIObjectVtbl

	EnumAdapters          uintptr
	MakeWindowAssociation uintptr
	GetWindowAssociation  uintptr
	CreateSwapChain       uintptr
	CreateSoftwareAdapter uintptr
}
type IDXGIFactory1Vtbl struct {
	_ structs.HostLayout
	IDXGIFactoryVtbl

	EnumAdapters1 uintptr
	IsCurrent     uintptr
}
