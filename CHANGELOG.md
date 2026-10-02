# Changelog

## [0.17.2]

### Fixed

- `UserAgentVersion` was pinned at `0.16.4`, so every consumer of this library
  advertised that version to its peers regardless of which release it was built
  against. It now tracks the released version, and a test asserts it matches
  the topmost `CHANGELOG.md` heading so the two cannot drift apart again.

### Changed

- Releases no longer attach binaries. This module has no `main` package, so the
  27 per-platform "binaries" previous releases published were Go object
  archives (`!<arch>` / `__.PKGDEF`), roughly 10 MB each and not runnable. The
  release is now sources-only, which is what a library release should be.
- Built with Go 1.26.5. (#5)
- Updated `go-flokicoin` to
  [v0.26.2](https://github.com/flokiorg/go-flokicoin/releases/tag/v0.26.2) and
  `walletd` to v0.2.1-beta, from v0.25.13-alpha and v0.1.5-beta. That also
  pulled `google.golang.org/grpc` from v1.71.0 to v1.79.3, closing
  GHSA-p77j-4mvh-x3m3.
- The release workflow now runs `go test -short`, the same command CI runs. It
  was running the full suite, which needs a chain to sync against and so could
  never pass on a runner -- it failed the first 0.17.2 release attempt.

## [0.17.1-beta]

flokicoin-neutrino v0.17.1-beta adds CI and fixes three real `go vet` findings, plus skips a couple of tests that don't belong in a fast CI gate.

### Fixes

- **`filterdb/db.go`**: `&filterData.BlockHash` double-pointered an already-pointer field, breaking `%s`/Stringer formatting. Removed the extra `&`.
- **`query.go`**: two log lines used `%d` on `float32` values (peer counts kept as float32 for a ratio comparison, and a fractional threshold). Cast counts to `int`, switched the threshold to `%.2f`.
- **`verification.go`**: `log.Debug(...)` doesn't interpret format verbs (needed `Debugf`); fixing that surfaced a real arg-count bug too — the format string references `tx %v` but `tx.Hash()` was never actually passed, so every subsequent arg was shifted by one position.

### Tests

- Skipped `TestNeutrinoSyncWithHeadersImport`, `TestNeutrinoImportThenP2PSync`, `TestNeutrinoSyncWithoutHeadersImport`, and `TestHandleHeaders` under `testing.Short()` — each spawns a real `lokid` process and syncs over actual P2P networking, not a good fit for a fast CI gate.
- Skipped `TestBlockCache` with a `TODO`: the test's `PowLimit` override doesn't actually bypass `CheckBlockSanity`'s check against a block's own encoded `Bits` target, so it fails deterministically against real historical difficulty. This is a real gap between the intended design (skip PoW post-header-sync) and actual behavior in a block-validation path — needs a maintainer's dedicated look, not a quick patch.

### CI

- Added `.github/workflows/ci.yaml`: runs `go build`, `go vet`, and `go test -short` on push to `main` and on pull requests.

Commit range: `0.17.0-beta..0.17.1-beta` (3 commits + 1 merge).

## [0.17.0-beta]

flokicoin-neutrino v0.17.0-beta introduces a new header-import fast-sync path and a batch of sync-engine and storage improvements ported from upstream neutrino.

### Header Import (new)

A new `chainimport` package allows a node to bootstrap its header chain from a local file or an HTTP source instead of syncing headers entirely over P2P:

- File and HTTP header import sources, with pluggable `HeaderImportSource` / `HttpClient` interfaces.
- Block and filter header validators, chain-continuity validation, and target-height verification.
- Divergence-region handling so an import can be reconciled against headers already present in the store.
- Configurable `WriteBatchSizePerRegion` and a larger default header-import batch size (65536).
- Extensive new test coverage (`headers_import_test.go`, benchmark comparisons between P2P sync and headers-import sync).

### Sync Engine

- Added `ProgressTimeout`, a query option that bounds batch idle timeouts instead of relying on wall-clock timeouts.
- `blockmanager` now bounds cfheader sync by an idle timeout rather than wall-clock time, and resets header state cleanly after a chain import.
- Skip-list ancestor lookups added to `headerlist.Node` and used in `lightHeaderCtx`, improving lookup performance on long header chains.

### Storage

- `headerfs` refactor: new `File` interface, pre-created header index sub-buckets, multi-index rollback and deletion, and a `runtime.GOOS`-based file implementation replacing the old build-flag split.
- `filterdb` falls back to `Update` when the backend lacks `BatchDB`.
- LRU cache gained a `Size` method with regression coverage.

### Fixes & Cleanup

- Fixed mocks for the fork's `Update`/`View` DB interface signature.
- `headerfs` now fails gracefully instead of panicking on write errors.
- Removed the deprecated `OnAlert` method; `ControlCFHeader` is deprecated in favor of `ValidateCFHeader`.
- Various linter fixes and long-line wraps in `chainimport`.

### Dependency Security

- Bumped `golang.org/x/crypto` (indirect, in both `go.mod` and `tools/go.mod`) from `v0.45.0` to `v0.52.0`, closing the same 7 `ssh`/`ssh/agent` advisories as the go-flokicoin bump. Not directly imported by this module either way.

Commit range: `0.16.6-beta..0.17.0-beta` (80 commits, plus a tools/ dependency bump).

## [0.16.6-beta]

### Dependency Updates

- **go-flokicoin**: Updated to `v0.25.13-alpha` for MuSig2 support and TestNet4 port fixes.
- Routine `go mod tidy` cleanup.

## [0.16.5-beta]

### Dependency Updates

- Bumped `github.com/stretchr/testify` to `v1.11.1`.
- Bumped `github.com/decred/dcrd/dcrec/secp256k1/v4` to `v4.4.0`.
- Updated `github.com/davecgh/go-spew` and `github.com/pmezard/go-difflib` to latest patch revisions.
- Routine `go mod tidy` cleanup following workspace sync.

## [0.16.4-beta]

- Updated identity to `lokitrino`/0.16.4 for network peers.
- Swapped build helper to install `lokid` tooling.
- Tightened walletdb test setup to pass observer flag explicitly.
- Aligned push tx error handling and naming with lokid.
- Reduced sync test block generation to keep harness lean.

## [0.16.3-beta]

- Neutrino now derives filter header checkpoints from `chaincfg.Params` instead of hardcoded maps.
- This change affects downstreams that depend on Neutrino: `flnd` and `twallet`.
- Behavior: If no checkpoint exists at a queried height in params, no filter header check is enforced; legacy map remains as a fallback.

## [0.16.2-beta]

- Add AuxPoW-compatible header I/O using SerializeHeader/DeserializeHeader while keeping 80-byte on-disk format.
- Update dependency: go-flokicoin v0.25.7-beta.
- Bump version to 0.16.2-beta (UserAgent 0.16.2).

Compatibility: No on-disk format change; existing headers remain valid.

## [0.16.1-beta.2]

#### Changed
- Updated go.mod to include new/updated module dependencies

## [0.16.1-beta]

#### Changed
- Fixed transaction lookup to match updated API

This is a **pre-release** for testing and feedback.
Developers and early adopters are encouraged to **report issues**.

## [0.16.0-beta]

- This is a **pre-release** for testing and feedback.
- Developers and early adopters are encouraged to **report issues**.
