package neutrino

import "github.com/flokiorg/go-flokicoin/chainutil"

// CacheableBlock is a wrapper around the chainutil.Block type which provides a
// Size method used by the cache to target certain memory usage.
type CacheableBlock struct {
	*chainutil.Block
}

// Size returns size of this block in bytes.
func (c *CacheableBlock) Size() (uint64, error) {
	return uint64(c.Block.MsgBlock().SerializeSize()), nil
}
