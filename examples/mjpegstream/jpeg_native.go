//go:build !turbo

package main

import (
	"image"
	"io"

	"image/jpeg"
)

func jpegQuality(q int) *jpeg.Options {
	return &jpeg.Options{Quality: q}
}

func encodeJpeg(w io.Writer, src image.Image, opts *jpeg.Options) {
	jpeg.Encode(w, src, opts)
}
