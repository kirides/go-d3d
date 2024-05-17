package outputduplication

import (
	"errors"
	"fmt"
	"image"

	"unsafe"

	"github.com/kirides/go-d3d"
	"github.com/kirides/go-d3d/d3d11"
	"github.com/kirides/go-d3d/dxgi"
	"github.com/kirides/go-d3d/outputduplication/swizzle"
)

type PointerInfo struct {
	pos dxgi.POINT

	size           dxgi.POINT
	shapeInBuffer  []byte
	shapeOutBuffer *image.RGBA
	visible        bool
}

type OutputDuplicator struct {
	device            *d3d11.ID3D11Device
	deviceCtx         *d3d11.ID3D11DeviceContext
	outputDuplication *dxgi.IDXGIOutputDuplication
	dxgiOutput        *dxgi.IDXGIOutput5

	stagedTex  *d3d11.ID3D11Texture2D
	surface    *dxgi.IDXGISurface
	mappedRect dxgi.DXGI_MAPPED_RECT
	size       dxgi.POINT

	pointerInfo PointerInfo
	DrawPointer bool

	// TODO: handle DPI? Do we need it?
	dirtyRects    []dxgi.RECT
	movedRects    []dxgi.DXGI_OUTDUPL_MOVE_RECT
	acquiredFrame bool
	needsSwizzle  bool // in case we use DuplicateOutput1, swizzle is not neccessery
}

func (dup *OutputDuplicator) initializeStage(texture *d3d11.ID3D11Texture2D) int32 {

	/*
		TODO: Only do this on changes!
	*/
	var hr int32
	desc := d3d11.D3D11_TEXTURE2D_DESC{}
	hr = texture.GetDesc(&desc)
	if d3d.HRESULT(hr).Failed() {
		return hr
	}

	desc.Usage = d3d11.D3D11_USAGE_STAGING
	desc.CPUAccessFlags = d3d11.D3D11_CPU_ACCESS_READ
	desc.BindFlags = 0
	desc.MipLevels = 1
	desc.ArraySize = 1
	desc.MiscFlags = 0
	desc.SampleDesc.Count = 1

	hr = dup.device.CreateTexture2D(&desc, &dup.stagedTex)
	if d3d.HRESULT(hr).Failed() {
		return hr
	}

	hr = dup.stagedTex.QueryInterface(dxgi.IID_IDXGISurface, &dup.surface)
	if d3d.HRESULT(hr).Failed() {
		return hr
	}
	dup.size = dxgi.POINT{X: int32(desc.Width), Y: int32(desc.Height)}

	return 0
}

func (dup *OutputDuplicator) Release() {
	dup.ReleaseFrame()
	if dup.stagedTex != nil {
		dup.stagedTex.Release()
		dup.stagedTex = nil
	}
	if dup.surface != nil {
		dup.surface.Release()
		dup.surface = nil
	}
	if dup.outputDuplication != nil {
		dup.outputDuplication.Release()
		dup.outputDuplication = nil
	}
	if dup.dxgiOutput != nil {
		dup.dxgiOutput.Release()
		dup.dxgiOutput = nil
	}
}

var ErrNoImageYet = errors.New("no image yet")

type unmapFn func() int32

func (dup *OutputDuplicator) ReleaseFrame() {
	if dup.acquiredFrame {
		dup.outputDuplication.ReleaseFrame()
		dup.acquiredFrame = false
	}
}

// returns DXGI_FORMAT_B8G8R8A8_UNORM data
func (dup *OutputDuplicator) Snapshot(timeoutMs uint) (unmapFn, *dxgi.DXGI_MAPPED_RECT, *dxgi.POINT, error) {
	var hr int32
	desc := dxgi.DXGI_OUTDUPL_DESC{}
	hr = dup.outputDuplication.GetDesc(&desc)
	if hr := d3d.HRESULT(hr); hr.Failed() {
		return nil, nil, nil, fmt.Errorf("failed to get the description. %w", hr)
	}

	if desc.DesktopImageInSystemMemory != 0 {
		// TODO: Figure out WHEN exactly this can occur, and if we can make use of it
		dup.size = dxgi.POINT{X: int32(desc.ModeDesc.Width), Y: int32(desc.ModeDesc.Height)}
		hr = dup.outputDuplication.MapDesktopSurface(&dup.mappedRect)
		if hr := d3d.HRESULT(hr); !hr.Failed() {
			return dup.outputDuplication.UnMapDesktopSurface, &dup.mappedRect, &dup.size, nil
		}
	}

	var desktop *dxgi.IDXGIResource
	var frameInfo dxgi.DXGI_OUTDUPL_FRAME_INFO

	// Release a possible previous frame
	// TODO: Properly use ReleaseFrame...

	dup.ReleaseFrame()
	hrF := dup.outputDuplication.AcquireNextFrame(uint32(timeoutMs), &frameInfo, &desktop)
	dup.acquiredFrame = true
	if hr := d3d.HRESULT(hrF); hr.Failed() {
		if hr == d3d.DXGI_ERROR_WAIT_TIMEOUT {
			return nil, nil, nil, ErrNoImageYet
		}
		return nil, nil, nil, fmt.Errorf("failed to AcquireNextFrame. %w", d3d.HRESULT(hrF))
	}
	// If we do not release the frame ASAP, we only get FPS / 2 frames :/
	// Something wrong here?
	// defer dup.ReleaseFrame()
	defer desktop.Release()

	if dup.DrawPointer {
		if err := dup.updatePointer(&frameInfo); err != nil {
			return nil, nil, nil, err
		}
	}

	if frameInfo.AccumulatedFrames == 0 {
		return nil, nil, nil, ErrNoImageYet
	}
	var desktop2d *d3d11.ID3D11Texture2D
	hr = desktop.QueryInterface(d3d11.IID_ID3D11Texture2D, &desktop2d)
	if hr := d3d.HRESULT(hr); hr.Failed() {
		return nil, nil, nil, fmt.Errorf("failed to QueryInterface(iid_ID3D11Texture2D, ...). %w", hr)
	}
	defer desktop2d.Release()

	if dup.stagedTex == nil {
		hr = dup.initializeStage(desktop2d)
		if hr := d3d.HRESULT(hr); hr.Failed() {
			return nil, nil, nil, fmt.Errorf("failed to InitializeStage. %w", hr)
		}
	}

	// NOTE: we could use a single, large []byte buffer and use it as storage for moved rects & dirty rects
	if frameInfo.TotalMetadataBufferSize > 0 {
		// Handling moved / dirty rects, to reduce GPU<->CPU memory copying
		moveRectsRequired := uint32(1)
		for {
			if len(dup.movedRects) < int(moveRectsRequired) {
				dup.movedRects = make([]dxgi.DXGI_OUTDUPL_MOVE_RECT, moveRectsRequired)
			}
			hr = dup.outputDuplication.GetFrameMoveRects(dup.movedRects, &moveRectsRequired)
			if hr := d3d.HRESULT(hr); hr.Failed() {
				if hr == d3d.DXGI_ERROR_MORE_DATA {
					continue
				}
				return nil, nil, nil, fmt.Errorf("failed to GetFrameMoveRects. %w", d3d.HRESULT(hr))
			}
			dup.movedRects = dup.movedRects[:moveRectsRequired]
			break
		}

		dirtyRectsRequired := uint32(1)
		for {
			if len(dup.dirtyRects) < int(dirtyRectsRequired) {
				dup.dirtyRects = make([]dxgi.RECT, dirtyRectsRequired)
			}
			hr = dup.outputDuplication.GetFrameDirtyRects(dup.dirtyRects, &dirtyRectsRequired)
			if hr := d3d.HRESULT(hr); hr.Failed() {
				if hr == d3d.DXGI_ERROR_MORE_DATA {
					continue
				}
				return nil, nil, nil, fmt.Errorf("failed to GetFrameDirtyRects. %w", d3d.HRESULT(hr))
			}
			dup.dirtyRects = dup.dirtyRects[:dirtyRectsRequired]
			break
		}

		box := d3d11.D3D11_BOX{
			Front: 0,
			Back:  1,
		}
		if len(dup.movedRects) == 0 {
			for i := 0; i < len(dup.dirtyRects); i++ {
				box.Left = uint32(dup.dirtyRects[i].Left)
				box.Top = uint32(dup.dirtyRects[i].Top)
				box.Right = uint32(dup.dirtyRects[i].Right)
				box.Bottom = uint32(dup.dirtyRects[i].Bottom)

				dup.deviceCtx.CopySubresourceRegion2D(dup.stagedTex, 0, box.Left, box.Top, 0, desktop2d, 0, &box)
			}
		} else {
			// TODO: handle moved rects, then dirty rects
			// for now, just update the whole image instead
			dup.deviceCtx.CopyResource2D(dup.stagedTex, desktop2d)
		}
	} else {
		// no frame metadata, copy whole image
		dup.deviceCtx.CopyResource2D(dup.stagedTex, desktop2d)
		if !dup.needsSwizzle {
			dup.needsSwizzle = true
		}
		print("no frame metadata\n")
	}

	hr = dup.surface.Map(&dup.mappedRect, dxgi.DXGI_MAP_READ)
	if hr := d3d.HRESULT(hr); hr.Failed() {
		return nil, nil, nil, fmt.Errorf("failed to surface_.Map(...). %v", hr)
	}
	return dup.surface.Unmap, &dup.mappedRect, &dup.size, nil
}

func (dup *OutputDuplicator) GetImage(img *image.RGBA, timeoutMs uint) error {
	unmap, mappedRect, size, err := dup.Snapshot(timeoutMs)
	if err != nil {
		return err
	}
	defer unmap()

	// docs are unclear, but pitch is the total width of each row
	dataSize := int(mappedRect.Pitch) * int(size.Y)
	data := unsafe.Slice((*byte)(unsafe.Pointer(mappedRect.PBits)), dataSize)

	contentWidth := int(size.X) * 4
	dataWidth := int(mappedRect.Pitch)

	var imgStart, dataStart, dataEnd int
	// copy source bytes into image.RGBA.Pix, skipping padding
	for i := 0; i < int(size.Y); i++ {
		dataEnd = dataStart + contentWidth
		copy(img.Pix[imgStart:], data[dataStart:dataEnd])
		imgStart += contentWidth
		dataStart += dataWidth
	}

	dup.drawPointer(img)
	if dup.needsSwizzle {
		swizzle.BGRA(img.Pix)
	}

	return nil
}

func (dup *OutputDuplicator) updatePointer(info *dxgi.DXGI_OUTDUPL_FRAME_INFO) error {
	if info.LastMouseUpdateTime == 0 {
		return nil
	}
	dup.pointerInfo.visible = info.PointerPosition.Visible != 0
	dup.pointerInfo.pos = info.PointerPosition.Position

	if info.PointerShapeBufferSize != 0 {
		// new shape
		if len(dup.pointerInfo.shapeInBuffer) < int(info.PointerShapeBufferSize) {
			dup.pointerInfo.shapeInBuffer = make([]byte, info.PointerShapeBufferSize)
		}
		var requiredSize uint32
		var pointerInfo dxgi.DXGI_OUTDUPL_POINTER_SHAPE_INFO

		hr := dup.outputDuplication.GetFramePointerShape(info.PointerShapeBufferSize,
			dup.pointerInfo.shapeInBuffer,
			&requiredSize,
			&pointerInfo,
		)
		if hr != 0 {
			return fmt.Errorf("unable to obtain frame pointer shape")
		}
		neededSize := pointerInfo.Width * pointerInfo.Height * 4
		dup.pointerInfo.shapeOutBuffer = image.NewRGBA(image.Rect(0, 0, int(pointerInfo.Width), int(pointerInfo.Height)))
		if len(dup.pointerInfo.shapeOutBuffer.Pix) < int(neededSize) {
			dup.pointerInfo.shapeOutBuffer.Pix = make([]byte, neededSize)
		}

		switch pointerInfo.Type {
		case dxgi.DXGI_OUTDUPL_POINTER_SHAPE_TYPE_MONOCHROME:
			dup.pointerInfo.size = dxgi.POINT{X: int32(pointerInfo.Width), Y: int32(pointerInfo.Height)}

			xor_offset := pointerInfo.Pitch * (pointerInfo.Height / 2)
			andMap := dup.pointerInfo.shapeInBuffer
			xorMap := dup.pointerInfo.shapeInBuffer[:xor_offset]
			out_pixels := dup.pointerInfo.shapeOutBuffer.Pix
			widthBytes := (pointerInfo.Width + 7) / 8

			imgHeight := pointerInfo.Height / 2

			for j := 0; j < int(imgHeight); j++ {
				bit := byte(0x80)

				for i := 0; i < int(pointerInfo.Width); i++ {
					andByte := andMap[j*int(widthBytes)+i/8]
					xorByte := xorMap[j*int(widthBytes)+i/8]
					andBit := 0
					if (andByte & bit) != 0 {
						andBit = 1
					}
					xorBit := 0
					if (xorByte & bit) != 0 {
						xorBit = 1
					}
					outDx := j*int(pointerInfo.Width)*4 + i*4
					if andBit == 0 {
						if xorBit == 0 {
							out_pixels[outDx+0] = 0x00
							out_pixels[outDx+1] = 0x00
							out_pixels[outDx+2] = 0x00
							out_pixels[outDx+3] = 0x00
						} else {
							out_pixels[outDx+0] = 0xFF
							out_pixels[outDx+1] = 0xFF
							out_pixels[outDx+2] = 0xFF
							out_pixels[outDx+3] = 0xFF
						}
					} else {
						if xorBit == 0 {
							out_pixels[outDx+0] = 0x00
							out_pixels[outDx+1] = 0x00
							out_pixels[outDx+2] = 0x00
							out_pixels[outDx+3] = 0x00
						} else {
							out_pixels[outDx+0] = 0x00
							out_pixels[outDx+1] = 0x00
							out_pixels[outDx+2] = 0x00
							out_pixels[outDx+3] = 0xFF
						}
					}
					if bit == 0x01 {
						bit = 0x80
					} else {
						bit = bit >> 1
					}
				}
			}
		case dxgi.DXGI_OUTDUPL_POINTER_SHAPE_TYPE_COLOR:
			dup.pointerInfo.size = dxgi.POINT{X: int32(pointerInfo.Width), Y: int32(pointerInfo.Height)}

			out, in := dup.pointerInfo.shapeOutBuffer.Pix, dup.pointerInfo.shapeInBuffer
			for j := 0; j < int(pointerInfo.Height); j++ {
				tout := out[j*int(pointerInfo.Pitch):]
				tin := in[j*int(pointerInfo.Pitch):]
				copy(tout, tin[:pointerInfo.Pitch])
			}
		case dxgi.DXGI_OUTDUPL_POINTER_SHAPE_TYPE_MASKED_COLOR:
			dup.pointerInfo.size = dxgi.POINT{X: int32(pointerInfo.Width), Y: int32(pointerInfo.Height)}

			// TODO: Properly add mask
			out, in := dup.pointerInfo.shapeOutBuffer.Pix, dup.pointerInfo.shapeInBuffer
			for j := 0; j < int(pointerInfo.Height); j++ {
				tout := out[j*int(pointerInfo.Pitch):]
				tin := in[j*int(pointerInfo.Pitch):]
				copy(tout, tin[:pointerInfo.Pitch])
			}
		default:
			dup.pointerInfo.size = dxgi.POINT{X: 0, Y: 0}
			return fmt.Errorf("unsupported type %v", pointerInfo.Type)
		}
	}
	return nil
}

func (dup *OutputDuplicator) drawPointer(img *image.RGBA) error {
	if !dup.DrawPointer {
		return nil
	}

	for j := 0; j < int(dup.pointerInfo.size.Y); j++ {
		for i := 0; i < int(dup.pointerInfo.size.X); i++ {
			col := dup.pointerInfo.shapeOutBuffer.At(i, j)
			_, _, _, a := col.RGBA()
			if a == 0 {
				// just dont draw invisible pixel?
				// TODO: correctly apply mask
				continue
			}

			img.Set(int(dup.pointerInfo.pos.X)+i, int(dup.pointerInfo.pos.Y)+j, col)
		}
	}
	return nil
}

func (ddup *OutputDuplicator) GetBounds() (image.Rectangle, error) {
	desc := dxgi.DXGI_OUTPUT_DESC{}
	hr := ddup.dxgiOutput.GetDesc(&desc)
	if hr := d3d.HRESULT(hr); hr.Failed() {
		return image.Rectangle{}, fmt.Errorf("failed at dxgiOutput.GetDesc. %w", hr)
	}

	return image.Rect(int(desc.DesktopCoordinates.Left), int(desc.DesktopCoordinates.Top), int(desc.DesktopCoordinates.Right), int(desc.DesktopCoordinates.Bottom)), nil
}

// NewIDXGIOutputDuplication creates a new OutputDuplicator
func NewIDXGIOutputDuplication(device *d3d11.ID3D11Device, deviceCtx *d3d11.ID3D11DeviceContext, output uint) (*OutputDuplicator, error) {
	// DEBUG

	var d3dDebug *d3d11.ID3D11Debug
	hr := device.QueryInterface(d3d11.IID_ID3D11Debug, &d3dDebug)
	if hr := d3d.HRESULT(hr); !hr.Failed() {
		defer d3dDebug.Release()

		var d3dInfoQueue *d3d11.ID3D11InfoQueue
		hr := d3dDebug.QueryInterface(d3d11.IID_ID3D11InfoQueue, &d3dInfoQueue)
		if hr := d3d.HRESULT(hr); hr.Failed() {
			return nil, fmt.Errorf("failed at device.QueryInterface. %w", hr)
		}
		defer d3dInfoQueue.Release()
		// defer d3dDebug.ReportLiveDeviceObjects(D3D11_RLDO_SUMMARY | D3D11_RLDO_DETAIL)
	}

	var dxgiDevice1 *dxgi.IDXGIDevice1
	hr = device.QueryInterface(dxgi.IID_IDXGIDevice1, &dxgiDevice1)
	if hr := d3d.HRESULT(hr); hr.Failed() {
		return nil, fmt.Errorf("failed at device.QueryInterface. %w", hr)
	}
	defer dxgiDevice1.Release()

	var pdxgiAdapter unsafe.Pointer
	hr = dxgiDevice1.GetParent(dxgi.IID_IDXGIAdapter1, &pdxgiAdapter)
	if hr := d3d.HRESULT(hr); hr.Failed() {
		return nil, fmt.Errorf("failed at dxgiDevice1.GetAdapter. %w", hr)
	}
	dxgiAdapter := (*dxgi.IDXGIAdapter1)(pdxgiAdapter)
	defer dxgiAdapter.Release()

	var dxgiOutput *dxgi.IDXGIOutput
	// const DXGI_ERROR_NOT_FOUND = 0x887A0002
	hr = int32(dxgiAdapter.EnumOutputs(uint32(output), &dxgiOutput))
	if hr := d3d.HRESULT(hr); hr.Failed() {
		return nil, fmt.Errorf("failed at dxgiAdapter.EnumOutputs. %w", hr)
	}
	defer dxgiOutput.Release()

	var dxgiOutput5 *dxgi.IDXGIOutput5
	hr = dxgiOutput.QueryInterface(dxgi.IID_IDXGIOutput5, &dxgiOutput5)
	if hr := d3d.HRESULT(hr); hr.Failed() {
		return nil, fmt.Errorf("failed at dxgiOutput.QueryInterface. %w", hr)
	}

	var dup *dxgi.IDXGIOutputDuplication
	hr = dxgiOutput5.DuplicateOutput1(dxgiDevice1, 0, []dxgi.DXGI_FORMAT{
		dxgi.DXGI_FORMAT_R8G8B8A8_UNORM,
		// using the former, we don't have to swizzle ourselves
		// DXGI_FORMAT_B8G8R8A8_UNORM,
	}, &dup)
	needsSwizzle := false
	if hr := d3d.HRESULT(hr); hr.Failed() {
		needsSwizzle = true
		// fancy stuff not supported :/
		// fmt.Printf("Info: failed to use dxgiOutput5.DuplicateOutput1, falling back to dxgiOutput1.DuplicateOutput. Missing manifest with DPI awareness set to \"PerMonitorV2\"? %v\n", _DXGI_ERROR(hr))
		var dxgiOutput1 *dxgi.IDXGIOutput1
		hr := dxgiOutput.QueryInterface(dxgi.IID_IDXGIOutput1, &dxgiOutput1)
		if hr := d3d.HRESULT(hr); hr.Failed() {
			dxgiOutput5.Release()
			return nil, fmt.Errorf("failed at dxgiOutput.QueryInterface. %w", hr)
		}
		defer dxgiOutput1.Release()
		hr = dxgiOutput1.DuplicateOutput(dxgiDevice1, &dup)
		if hr := d3d.HRESULT(hr); hr.Failed() {
			dxgiOutput5.Release()
			return nil, fmt.Errorf("failed at dxgiOutput1.DuplicateOutput. %w", hr)
		}
	}

	return &OutputDuplicator{device: device, deviceCtx: deviceCtx, outputDuplication: dup, needsSwizzle: needsSwizzle, dxgiOutput: dxgiOutput5}, nil
}
