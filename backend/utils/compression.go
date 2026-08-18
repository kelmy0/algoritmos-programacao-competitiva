package utils

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"
)

func CompressText(text string) ([]byte, error) {
	if text == "" {
		return []byte{}, nil
	}

	var buf bytes.Buffer
	writer, err := flate.NewWriter(&buf, flate.DefaultCompression)
	if err != nil {
		return nil, fmt.Errorf("failed to create compression writer: %w", err)
	}

	if _, err := writer.Write([]byte(text)); err != nil {
		_ = writer.Close()
		return nil, fmt.Errorf("failed to write data to compressor: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close compression writer: %w", err)
	}

	return buf.Bytes(), nil
}

func DecompressText(compressed []byte) (string, error) {
	if len(compressed) == 0 {
		return "", nil
	}

	reader := flate.NewReader(bytes.NewReader(compressed))
	defer reader.Close()

	decompressed, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to decompress data: %w", err)
	}

	return string(decompressed), nil
}
