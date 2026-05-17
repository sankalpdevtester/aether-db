// Package feature provides data compression functionality for AetherDB.
package feature

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/pierrec/lz4/v4"
)

// CompressionType represents the type of compression used.
type CompressionType uint8

const (
	// None represents no compression.
	None CompressionType = iota
	// LZ4 represents LZ4 compression.
	LZ4
)

// Compressor is an interface for compressing data.
type Compressor interface {
	Compress(data []byte) ([]byte, error)
	Decompress(data []byte) ([]byte, error)
}

// LZ4Compressor is an implementation of the Compressor interface using LZ4.
type LZ4Compressor struct{}

// NewLZ4Compressor returns a new LZ4Compressor instance.
func NewLZ4Compressor() *LZ4Compressor {
	return &LZ4Compressor{}
}

// Compress compresses the given data using LZ4.
func (c *LZ4Compressor) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	encoder := lz4.NewWriter(&buf)
	_, err := encoder.Write(data)
	if err != nil {
		return nil, err
	}
	err = encoder.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decompress decompresses the given data using LZ4.
func (c *LZ4Compressor) Decompress(data []byte) ([]byte, error) {
	decoder := lz4.NewReader(bytes.NewReader(data))
	defer decoder.Close()
	return io.ReadAll(decoder)
}

// CompressData compresses the given data using the specified compression type.
func CompressData(data []byte, compressionType CompressionType) ([]byte, error) {
	switch compressionType {
	case LZ4:
		compressor := NewLZ4Compressor()
		return compressor.Compress(data)
	case None:
		return data, nil
	default:
		return nil, errors.New("unsupported compression type")
	}
}

// DecompressData decompresses the given data using the specified compression type.
func DecompressData(data []byte, compressionType CompressionType) ([]byte, error) {
	switch compressionType {
	case LZ4:
		compressor := NewLZ4Compressor()
		return compressor.Decompress(data)
	case None:
		return data, nil
	default:
		return nil, errors.New("unsupported compression type")
	}
}

// CompressionHeader represents the header of a compressed data block.
type CompressionHeader struct {
	CompressionType CompressionType
	OriginalSize    uint32
}

// Encode encodes the compression header into a byte slice.
func (h *CompressionHeader) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, h.CompressionType)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, h.OriginalSize)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode decodes the compression header from a byte slice.
func (h *CompressionHeader) Decode(data []byte) error {
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.BigEndian, &h.CompressionType)
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.BigEndian, &h.OriginalSize)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// Example usage:
	data := []byte("Hello, World!")
	compressionType := LZ4
	compressedData, err := CompressData(data, compressionType)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Compressed data:", compressedData)

	decompressedData, err := DecompressData(compressedData, compressionType)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Decompressed data:", decompressedData)
}