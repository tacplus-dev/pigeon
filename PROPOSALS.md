# Proposals: pigeon

Written by loopers. Each proposal: what, why, links. Jody promotes accepted ones into ROADMAP.md (via an agent);
the looper then marks it `promoted` or `declined` here with the date.

## Open
- **Close #1 as superseded by #5; update ROADMAP.md's floor lines** (2026-10-29, #5): #1 asks for the floor 2.4, which can no longer build: stratum (via coop) declares 2.5 (tacplus-dev/stratum@206a973), so #5 moved the floor to 2.5. ROADMAP.md still says "tacplus floor 2.4" under Now and Later, and #3 and #4 are blocked on #1. They should wait on #5 (itself waiting on tacplus-dev/looper#27) instead.
- **Close #3 as superseded by #6** (2026-10-29, #6): #3 asks for golang.org/x/net v0.60.0+; #6 moved it to v0.62.0 and ran what #3 waits for: `tacplus build` (2.5.1, Go on PATH), `go mod tidy` (no diff; the stale v0.57.0 hashes are gone), `go build`/`go test ./...`, and govulncheck v1.9.0: `No vulnerabilities found.`, so none of GO-2026-6603/6610/6611/6612/6617 remain. The run through the looper's wrapper no longer waits: on 2026-10-30, with the fixed wrapper (tacplus-dev/looper@7c1f306, tacplus 2.5.3 in `build:2.5`), the same gates and govulncheck pass on main (863e151).

## Promoted / declined
