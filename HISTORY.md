# History: pigeon

Append-only: add entries at the end, never edit or remove one (CI enforces it). One line each:
`- <YYYY-MM-DD> <who: looper|agent|jody> <kind: issue|decision|release|deploy> <one line> (<links>)`

- 2026-10-30 looper issue Gates fail on main: stratum (via coop) needs tacplus 2.5; move the floor from 2.4 to 2.5 (#5, 164e73e)
- 2026-10-30 looper issue Vulnerable golang.org/x/net v0.61.0 (GO-2026-6866, GO-2026-6860): bump to v0.62.0 (#6, 21ba35f)
- 2026-11-03 looper issue Host.OnMessage drops the inbound trace_id on the results and failures it publishes (#7, c593490)
- 2026-11-06 looper issue Declare the tacplus floor: 2.4 (#1)
- 2026-11-08 looper issue tacplus floor 2.5 is behind the latest released line (2.6): move the floor to 2.6 and commit tacplus.lock (#8, 7ddb081)
- 2026-11-08 looper issue Pin the build: floor 2.6 and a committed tacplus.lock, so the build pipeline and its security checks run on this repository (#9, 7604a0d)
- 2026-11-14 looper release v0.1.3 (https://github.com/tacplus-dev/pigeon/releases/tag/v0.1.3)
- 2026-11-14 looper issue STATE.md's Blocked on names a closed issue (tacplus-dev/looper#12): rewrite it to name what blocks #2 now (tacplus-dev/looper#115) (#11, 93f3f92)
- 2026-11-15 looper issue Release: tag vMAJOR.MINOR.PATCH (start at v0.1.0) and publish its image (#2, 59b9f0e)
- 2026-11-15 looper release v0.1.4 (https://github.com/tacplus-dev/pigeon/releases/tag/v0.1.4)
- 2026-11-15 looper issue Pinned build image build:2.6.1 ships go1.27.3, vulnerable to GO-2026-6859 (net/http): pin tacplus 2.6.8 (go1.27.4) (#10, 2fb48ce)
- 2026-11-17 looper issue New LOW finding in the nightly security report: trivy DS-0026 (no HEALTHCHECK in Dockerfile) (#14)
- 2026-11-17 looper issue Release image keeps the base's vulnerable zlib 1.3.2-r1 (CVE-2026-4738): the Dockerfile never runs apk upgrade (#13)
- 2026-11-17 looper release v0.1.6 (https://github.com/tacplus-dev/pigeon/releases/tag/v0.1.6)
- 2026-11-17 looper issue Release image's launcher is built with go1.27.3 (GO-2026-6859): LAUNCHER_COMMIT c54ae30 predates launcher's go1.27.4 toolchain; move it to launcher v0.1.0 (#12)
- 2026-11-20 looper release v0.1.7 (https://github.com/tacplus-dev/pigeon/releases/tag/v0.1.7)
- 2026-11-20 looper issue Release image v0.1.6 ships openssl 3.5.9-r0 (CVE-2026-107496, -88822, -88823, -95691, -95692): Alpine 3.24 has 3.5.9-r1; cut v0.1.7 to rebuild (#15)
