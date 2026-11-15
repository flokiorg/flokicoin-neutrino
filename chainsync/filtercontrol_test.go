package chainsync

import (
	"testing"

	"github.com/flokiorg/go-flokicoin/chaincfg"
	"github.com/flokiorg/go-flokicoin/wire"
	"github.com/stretchr/testify/require"
)

// TestValidateCFHeader tests filter header validation against checkpoints,
// including successful validation, mismatch detection, and unknown height
// handling.
func TestValidateCFHeader(t *testing.T) {
	t.Parallel()

	// We'll construct params with a custom checkpoint for this test.
	height := uint32(999)
	header := hashFromStr(
		"4a242283a406a7c089f671bb8df7671e5d5e9ba577cea1047d30a7f4919d" +
			"f193",
	)
	params := chaincfg.MainNetParams
	params.Checkpoints = []chaincfg.Checkpoint{{
		Height: int32(height),
		Hash:   header,
	}}

	// Expect the control at height to succeed.
	err := ValidateCFHeader(
		params, wire.GCSFilterRegular, height, header,
	)
	require.NoError(t, err)

	// Pass an invalid header, this should return an error.
	header = hashFromStr(
		"000000000006a7c089f671bb8df7671e5d5e9ba577cea1047d30a7f4919d" +
			"f193",
	)
	err = ValidateCFHeader(
		params, wire.GCSFilterRegular, height, header,
	)
	require.ErrorIs(t, err, ErrCheckpointMismatch)

	// Finally, control an unknown height. This should also pass since we
	// don't have the checkpoint stored.
	err = ValidateCFHeader(
		params, wire.GCSFilterRegular, 99, header,
	)
	require.NoError(t, err)
}
