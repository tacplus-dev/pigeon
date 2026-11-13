// Copyright (c) 2026 Jody Scott
// SPDX-License-Identifier: Apache-2.0

// Package ci_test tests the release scripts in ci/ (.github/workflows/release.yml) on throwaway
// git repositories: which image tags a release gets, which GitHub release is Latest, its notes,
// and which build image builds a tree.
package ci_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// repo makes a git repository with one commit per tag, tagged in the order given, and returns its path.
func repo(t *testing.T, tags ...string) string {
	t.Helper()
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	for _, tag := range tags {
		git(t, dir, "commit", "-q", "--allow-empty", "-m", "change for "+tag)
		git(t, dir, "tag", "-a", "-m", tag, tag)
	}
	return dir
}

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com",
		"GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_NOSYSTEM=1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out))
}

// script runs ci/<name> in dir and returns its standard output, and whether it succeeded.
func script(t *testing.T, dir, name string, args ...string) (string, bool) {
	t.Helper()
	path, err := filepath.Abs(name)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(path, args...)
	cmd.Dir = dir
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err == nil
}

func TestReleaseTags(t *testing.T) {
	for _, tc := range []struct {
		name    string
		tags    []string
		release string
		want    string
	}{
		{"the first release starts its line", []string{"v0.1.0"}, "v0.1.0", "0.1.0 0.1"},
		{"a patch moves its line", []string{"v0.1.0", "v0.1.1"}, "v0.1.1", "0.1.1 0.1"},
		{"an older patch never moves its line back", []string{"v0.1.0", "v0.1.1", "v0.1.2"}, "v0.1.0", "0.1.0"},
		{"a patch to an older line moves that line, not the newer one", []string{"v0.1.0", "v0.2.0", "v0.1.1"}, "v0.1.1", "0.1.1 0.1"},
		{"a newer line does not stop its own line moving", []string{"v0.1.0", "v0.1.1", "v0.2.0"}, "v0.2.0", "0.2.0 0.2"},
		{"patches compare by number, not text", []string{"v0.1.9", "v0.1.10"}, "v0.1.10", "0.1.10 0.1"},
		{"patches compare by number, not text (the lower)", []string{"v0.1.9", "v0.1.10"}, "v0.1.9", "0.1.9"},
		{"0.10 is another line than 0.1", []string{"v0.1.0", "v0.10.0", "v0.1.1"}, "v0.1.1", "0.1.1 0.1"},
		{"a pre-release gets only its own tag", []string{"v0.1.0", "v0.2.0-rc.1"}, "v0.2.0-rc.1", "0.2.0-rc.1"},
		{"a pre-release does not stop its line moving", []string{"v0.2.0-rc.1", "v0.2.0"}, "v0.2.0", "0.2.0 0.2"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := script(t, repo(t, tc.tags...), "release-tags.sh", tc.release)
			if !ok || got != tc.want {
				t.Errorf("release-tags.sh %s with %v = %q (ok %v), want %q", tc.release, tc.tags, got, ok, tc.want)
			}
		})
	}
	for _, bad := range []string{"0.1.0", "v0.1", "v1.2.3.4", "release"} {
		if got, ok := script(t, repo(t, "v0.1.0"), "release-tags.sh", bad); ok {
			t.Errorf("release-tags.sh %s = %q, want it refused", bad, got)
		}
	}
}

func TestReleaseLatest(t *testing.T) {
	for _, tc := range []struct {
		name    string
		tags    []string
		release string
		want    string
	}{
		{"the first release is Latest", []string{"v0.1.0"}, "v0.1.0", "--latest"},
		{"the highest release is Latest", []string{"v0.1.0", "v0.1.1", "v0.2.0"}, "v0.2.0", "--latest"},
		{"an older line's patch is never Latest", []string{"v0.1.0", "v0.2.0", "v0.1.1"}, "v0.1.1", "--latest=false"},
		{"an older release published later is not Latest", []string{"v0.1.0", "v0.1.1", "v0.1.2"}, "v0.1.0", "--latest=false"},
		{"by number, not text", []string{"v0.9.0", "v0.10.0"}, "v0.10.0", "--latest"},
		{"by number, not text (the lower)", []string{"v0.9.0", "v0.10.0"}, "v0.9.0", "--latest=false"},
		{"a pre-release is a pre-release", []string{"v0.1.0", "v0.2.0-rc.1"}, "v0.2.0-rc.1", "--prerelease"},
		{"a pre-release does not hide the highest release", []string{"v0.1.0", "v0.2.0-rc.1"}, "v0.1.0", "--latest"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := script(t, repo(t, tc.tags...), "release-latest.sh", tc.release)
			if !ok || got != tc.want {
				t.Errorf("release-latest.sh %s with %v = %q (ok %v), want %q", tc.release, tc.tags, got, ok, tc.want)
			}
		})
	}
}

func TestReleaseNotes(t *testing.T) {
	dir := repo(t, "v0.1.0", "v0.1.1", "v0.2.0-rc.1", "v0.2.0")
	git(t, dir, "checkout", "-q", "-b", "release/0.1", "v0.1.1")
	git(t, dir, "commit", "-q", "--allow-empty", "-m", "a fix on the older line")
	git(t, dir, "tag", "-a", "-m", "v0.1.2", "v0.1.2")

	for _, tc := range []struct {
		name    string
		release string
		want    []string
		notWant []string
	}{
		{"the first release lists every commit", "v0.1.0", []string{"Changes:", "change for v0.1.0"}, nil},
		{"a release lists the commits since the previous one", "v0.1.1",
			[]string{"Changes since v0.1.0:", "change for v0.1.1"}, []string{"change for v0.1.0"}},
		{"a pre-release is not a previous release", "v0.2.0",
			[]string{"Changes since v0.1.1:", "change for v0.2.0", "change for v0.2.0-rc.1"}, []string{"change for v0.1.1"}},
		{"a patch on an older line lists that line's changes", "v0.1.2",
			[]string{"Changes since v0.1.1:", "a fix on the older line"}, []string{"change for v0.2.0"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := script(t, dir, "release-notes.sh", tc.release)
			if !ok {
				t.Fatalf("release-notes.sh %s failed", tc.release)
			}
			for _, w := range tc.want {
				if !strings.Contains(got, w) {
					t.Errorf("release-notes.sh %s = %q, want it to contain %q", tc.release, got, w)
				}
			}
			for _, w := range tc.notWant {
				if strings.Contains(got, w) {
					t.Errorf("release-notes.sh %s = %q, want it not to contain %q", tc.release, got, w)
				}
			}
		})
	}
	if got, ok := script(t, dir, "release-notes.sh", "v9.9.9"); ok {
		t.Errorf("release-notes.sh for a missing tag = %q, want it refused", got)
	}
}

func TestBuildImage(t *testing.T) {
	const pinned = "ghcr.io/tacplus-dev/tacplus/build:2.6.1@sha256:91f42995e0e3d6f959a7f75177caf94b2aded90fc36c8492dbcd320ca0d187f5"
	for _, tc := range []struct {
		name  string
		files map[string]string
		want  string // "" when it must refuse
	}{
		{"the image tacplus.lock pins", map[string]string{
			"tacplus.yaml": "kind: service\ntacplus: \"2.6\"\n",
			"tacplus.lock": "tacplus: 2.6.1\nbuild_image: " + pinned + "\n",
		}, pinned},
		{"before tacplus.lock: the floor's build image", map[string]string{
			"tacplus.yaml": "kind: service\n# The compiler floor.\ntacplus: \"2.5\"\nlanguage: go\n",
		}, "ghcr.io/tacplus-dev/tacplus/build:2.5"},
		{"a lock without a digest is refused", map[string]string{
			"tacplus.yaml": "tacplus: \"2.6\"\n",
			"tacplus.lock": "tacplus: 2.6.1\nbuild_image: ghcr.io/tacplus-dev/tacplus/build:2.6.1\n",
		}, ""},
		{"no lock and no floor is refused", map[string]string{"tacplus.yaml": "kind: service\n"}, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			for name, content := range tc.files {
				if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			got, ok := script(t, ".", "build-image.sh", dir)
			switch {
			case tc.want == "" && ok:
				t.Errorf("build-image.sh = %q, want it refused", got)
			case tc.want != "" && (!ok || got != tc.want):
				t.Errorf("build-image.sh = %q (ok %v), want %q", got, ok, tc.want)
			}
		})
	}
}
