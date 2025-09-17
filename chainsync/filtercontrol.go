package chainsync

import (
    "fmt"

    "github.com/flokiorg/go-flokicoin/chaincfg"
    "github.com/flokiorg/go-flokicoin/chaincfg/chainhash"
    "github.com/flokiorg/go-flokicoin/wire"
)

// ErrCheckpointMismatch is returned if given filter headers don't pass our
// control check.
var ErrCheckpointMismatch = fmt.Errorf("checkpoint doesn't match")

// Deprecated: legacy hardcoded filter header checkpoints map.
// Left defined for backward compatibility with older tests that may mutate it.
// New logic derives checkpoints from chaincfg.Params.Checkpoints.
var filterHeaderCheckpoints = map[wire.FlokicoinNet]map[uint32]*chainhash.Hash{}

// checkpointHashAtHeight returns the checkpoint hash for the given height from
// the provided chain parameters, or nil if no checkpoint exists at that height.
func checkpointHashAtHeight(params chaincfg.Params, height uint32) *chainhash.Hash {
    for _, cp := range params.Checkpoints {
        if uint32(cp.Height) == height {
            return cp.Hash
        }
    }
    return nil
}

// ControlCFHeader controls the given filter header against our list of
// checkpoints. It returns ErrCheckpointMismatch if we have a checkpoint at the
// given height, and it doesn't match.
func ControlCFHeader(params chaincfg.Params, fType wire.FilterType,
    height uint32, filterHeader *chainhash.Hash) error {

    if fType != wire.GCSFilterRegular {
        return fmt.Errorf("unsupported filter type %v", fType)
    }

    // Prefer checkpoints provided by the active chain parameters.
    if hash := checkpointHashAtHeight(params, height); hash != nil {
        if *filterHeader != *hash {
            return ErrCheckpointMismatch
        }
        return nil
    }

    // Fallback to legacy hardcoded map if present (primarily for tests or
    // networks without declared checkpoints in params).
    if control, ok := filterHeaderCheckpoints[params.Net]; ok {
        if hash, ok := control[height]; ok {
            if *filterHeader != *hash {
                return ErrCheckpointMismatch
            }
        }
    }

    return nil
}

// hashFromStr makes a chainhash.Hash from a valid hex string. If the string is
// invalid, a nil pointer will be returned.
func hashFromStr(hexStr string) *chainhash.Hash {
	hash, _ := chainhash.NewHashFromStr(hexStr)
	return hash
}
