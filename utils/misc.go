package utils

import (
	"image/gif"
	"image/png"
	"io"
)

func PngToGif(gifWriter io.Writer, pngReader io.Reader) error {
	pngImg, err := png.Decode(pngReader)
	if err != nil {
		return err
	}

	return gif.Encode(gifWriter, pngImg, nil)
}
