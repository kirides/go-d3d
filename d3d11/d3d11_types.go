package d3d11

import (
	"structs"

	"github.com/kirides/go-d3d/dxgi"
)

type D3D11_BOX struct {
	_ structs.HostLayout

	Left, Top, Front, Right, Bottom, Back uint32
}

type D3D11_TEXTURE2D_DESC struct {
	_ structs.HostLayout

	Width          uint32
	Height         uint32
	MipLevels      uint32
	ArraySize      uint32
	Format         uint32
	SampleDesc     dxgi.DXGI_SAMPLE_DESC
	Usage          D3D11_USAGE
	BindFlags      uint32
	CPUAccessFlags uint32
	MiscFlags      uint32
}

type D3D11_USAGE uint32

const (
	D3D11_USAGE_DEFAULT   D3D11_USAGE = 0
	D3D11_USAGE_IMMUTABLE D3D11_USAGE = 1
	D3D11_USAGE_DYNAMIC   D3D11_USAGE = 2
	D3D11_USAGE_STAGING   D3D11_USAGE = 3
)
