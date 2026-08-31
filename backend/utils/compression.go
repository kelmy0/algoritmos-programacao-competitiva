package utils

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"sync"
)

var flateReaderPool = sync.Pool{
	New: func() any {
		return flate.NewReader(bytes.NewReader(nil))
	},
}

var flateWriterPool = sync.Pool{
	New: func() any {
		w, _ := flate.NewWriter(io.Discard, flate.DefaultCompression)
		return w
	},
}

func CompressText(text string) ([]byte, error) {
	if text == "" {
		return []byte{}, nil
	}

	var buf bytes.Buffer
	buf.Grow(len(text))

	fw := flateWriterPool.Get().(*flate.Writer)

	fw.Reset(&buf)

	defer func() {
		_ = fw.Close()
		flateWriterPool.Put(fw)
	}()

	if _, err := io.WriteString(fw, text); err != nil {
		return nil, fmt.Errorf("failed to write data to compressor: %w", err)
	}

	if err := fw.Close(); err != nil {
		return nil, fmt.Errorf("failed to flush compressor data: %w", err)
	}

	return buf.Bytes(), nil
}

func DecompressText(compressed []byte) (string, error) {
	if len(compressed) == 0 {
		return "", nil
	}

	fr := flateReaderPool.Get().(flate.Resetter)
	if err := fr.Reset(bytes.NewReader(compressed), nil); err != nil {
		return "", fmt.Errorf("failed to reset flate reader: %w", err)
	}

	defer flateReaderPool.Put(fr)

	buf := bytes.NewBuffer(make([]byte, 0, len(compressed)*4))

	if _, err := io.Copy(buf, fr.(io.Reader)); err != nil {
		return "", fmt.Errorf("failed to decompress data: %w", err)
	}

	return buf.String(), nil
}
