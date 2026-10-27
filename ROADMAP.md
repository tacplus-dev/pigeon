# Roadmap: pigeon

Owned by Jody; changed only by an agent at Jody's request. Loopers propose in PROPOSALS.md.

## Now
The generic pigeon host (`pigeonhost`): one run per Kubernetes Job created by the coop, speaking to its caller on the run's coop thread; built-in kind `echo.ping`.
Go. Unreleased; tacplus floor 2.4 (declared, #1 still open).

## Next
- Nothing approved and unblocked; every open approved issue is in Later.

## Later
- Declare the tacplus floor 2.4 and pass the gates on it (#1), blocked on tacplus-dev/tacplus#189
- Release v0.1.0 and publish its image (#2), blocked on tacplus-dev/looper#12 and #1
- golang.org/x/net to v0.60.0+ for five HTTP/2 advisories (#3), blocked on #1
- Base image alpine:3.22 to alpine:3.24 (#4), blocked on #1

## Not doing
- A Rust rewrite: Rust is the tacplus-dev default and pigeon should have been Rust, but Jody accepts Go here for now
