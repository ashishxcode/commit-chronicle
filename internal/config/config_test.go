package config

import (
	"os"
	"path/filepath"
	"testing"
)

// mkRepo creates dir/.git so isGitRepo treats dir as a repo.
func mkRepo(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
}

func TestCountRepos(t *testing.T) {
	root := t.TempDir()
	mkRepo(t, filepath.Join(root, "a"))
	mkRepo(t, filepath.Join(root, "b"))
	if err := os.MkdirAll(filepath.Join(root, "notarepo"), 0o755); err != nil {
		t.Fatal(err)
	}
	if got := CountRepos(root); got != 2 {
		t.Errorf("CountRepos = %d, want 2", got)
	}
	if got := CountRepos(filepath.Join(root, "does-not-exist")); got != 0 {
		t.Errorf("CountRepos(missing) = %d, want 0", got)
	}
}

func TestSaveRootDedupes(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	root := filepath.Join(dir, "work")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}

	p1, err := SaveRoot(root)
	if err != nil {
		t.Fatalf("SaveRoot: %v", err)
	}
	if _, err := SaveRoot(root); err != nil { // second save should be a no-op
		t.Fatalf("SaveRoot (dup): %v", err)
	}

	lines, ok := linesFromFile(p1)
	if !ok || len(lines) != 1 {
		t.Fatalf("roots file = %v (ok=%v), want exactly one entry", lines, ok)
	}
	if filepath.Clean(lines[0]) != filepath.Clean(root) {
		t.Errorf("saved %q, want %q", lines[0], root)
	}
}
