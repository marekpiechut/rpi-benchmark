package image_size

import (
	"encoding/binary"
	"errors"
	"io"
)

const pngSignature = "PNG\r\n\x1a\n"
const pngImageHeaderChunkName = "IHDR"
const pngFriedChunkName = "CgBI"

func isPng(reader io.ReaderAt) bool {
	magicChars := make([]byte, 8)
	bytes, err := reader.ReadAt(magicChars, 1)

	if err != nil || string(magicChars[:7]) != pngSignature {
		return false
	}

	bytes, err = reader.ReadAt(magicChars, 12)
	if err != nil {
		return false
	}

	if string(magicChars[:4]) == pngFriedChunkName {
		bytes, err = reader.ReadAt(magicChars, 28)
		if err != nil || bytes < 4 {
			return false
		}
	}

	return string(magicChars[:4]) == pngImageHeaderChunkName
}

func GetPngSize(reader io.ReaderAt) (Dimension, error) {
	if !isPng(reader) {
		return Dimension{}, errors.New("not a png")
	}
	buffer := make([]byte, 8)
	_, err := reader.ReadAt(buffer, 12)
	if err != nil {
		return Dimension{}, err
	}

	var offset int64 = 16
	if string(buffer) == pngFriedChunkName {
		offset = 32
	}

	var width uint32
	var height uint32

	_, err = reader.ReadAt(buffer, offset)
	if err != nil {
		return Dimension{}, err
	}
	width = binary.BigEndian.Uint32(buffer[:4])
	height = binary.BigEndian.Uint32(buffer[4:])

	return Dimension{int(width), int(height)}, nil
}
