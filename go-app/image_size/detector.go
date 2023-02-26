package image_size

import "io"

type sizeDetector = func(reader io.ReaderAt) (Dimension, error)

var firstBytes = map[byte]sizeDetector{
	0x89: GetPngSize,
}

func getSizeDetector(reader io.ReaderAt) (sizeDetector, error) {
	buffer := make([]byte, 1)
	_, err := reader.ReadAt(buffer, 0)
	if err != nil {
		return nil, err
	}
	detector := firstBytes[buffer[0]]
	return detector, nil
}

func DetectSize(reader io.ReaderAt) (Dimension, error) {
	detector, err := getSizeDetector(reader)
	if err != nil || detector == nil {
		return Dimension{}, err
	}

	return detector(reader)
}
