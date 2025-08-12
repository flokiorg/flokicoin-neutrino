package chainimport

import (
	"fmt"
	"io"

	"github.com/flokiorg/flokicoin-neutrino/headerfs"
	"github.com/flokiorg/go-flokicoin/wire"
)

// blockHeader represents a block header that can be imported into the chain
// store. It wraps a headerfs.BlockHeader with additional functionality needed
// for the import process.
type blockHeader struct {
	headerfs.BlockHeader
}

// Compile-time assertion to ensure blockHeader implements HeaderBase interface.
var _ Header = (*blockHeader)(nil)

// newBlockHeader creates and returns a new block header instance.
func newBlockHeader() Header {
	return &blockHeader{
		BlockHeader: headerfs.BlockHeader{
			BlockHeader: &wire.BlockHeader{},
		},
	}
}

// Deserialize reconstructs a block header from binary data at the specified
// height.
func (b *blockHeader) Deserialize(r io.Reader, height uint32) error {
	// Deserialize the fixed-size (80-byte) wire.BlockHeader portion.
	// Import sources store the base header only, without any AuxPoW
	// payload, matching headerfs's own on-disk flat-file format.
	if err := b.BlockHeader.BlockHeader.DeserializeHeader(r); err != nil {
		return fmt.Errorf("failed to deserialize wire.BlockHeader: "+
			"%w", err)
	}

	// Set block header height.
	b.BlockHeader.Height = height

	return nil
}

// filterHeader represents a filter header that can be imported into the chain
// store. It wraps a headerfs.FilterHeader with additional functionality needed
// for the import process.
type filterHeader struct {
	headerfs.FilterHeader
}

// Compile-time assertion to ensure filterHeader implements HeaderBase
// interface.
var _ Header = (*filterHeader)(nil)

// newFilterHeader creates and returns a new filter header instance.
func newFilterHeader() Header {
	return &filterHeader{
		FilterHeader: headerfs.FilterHeader{},
	}
}

// Deserialize reconstructs a filter header from binary data at the specified
// height.
func (f *filterHeader) Deserialize(r io.Reader, height uint32) error {
	// Read the filter hash (32 bytes).
	if _, err := io.ReadFull(r, f.FilterHash[:]); err != nil {
		return fmt.Errorf("failed to read filter hash: %w", err)
	}

	// Set filter header height.
	f.FilterHeader.Height = height

	return nil
}
