//go:build turbo

package main

import (
	"image"
	"io"

	"github.com/viam-labs/go-libjpeg/jpeg"
)

func jpegQuality(q int) *jpeg.EncoderOptions {
	return &jpeg.EncoderOptions{Quality: q}
}

func encodeJpeg(w io.Writer, src image.Image, opts *jpeg.EncoderOptions) {
	jpeg.Encode(w, src, opts)
}
