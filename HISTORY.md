# History: pigeon

Append-only: add entries at the end, never edit or remove one (CI enforces it). One line each:
`- <YYYY-MM-DD> <who: looper|agent|jody> <kind: issue|decision|release|deploy> <one line> (<links>)`

- 2026-10-30 looper issue Gates fail on main: stratum (via coop) needs tacplus 2.5; move the floor from 2.4 to 2.5 (#5, 164e73e)
- 2026-10-30 looper issue Vulnerable golang.org/x/net v0.61.0 (GO-2026-6866, GO-2026-6860): bump to v0.62.0 (#6, 21ba35f)
- 2026-11-03 looper issue Host.OnMessage drops the inbound trace_id on the results and failures it publishes (#7, c593490)
- 2026-11-06 looper issue Declare the tacplus floor: 2.4 (#1)
- 2026-11-08 looper issue tacplus floor 2.5 is behind the latest released line (2.6): move the floor to 2.6 and commit tacplus.lock (#8, 7ddb081)
