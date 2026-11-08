package cmd

import (
	"bytes"
	"crypto/rand"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func run(t *testing.T, args ...string) {
	t.Helper()
	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("fastzip-cli %v: %v", args, err)
	}
}

func readTree(t *testing.T, root string) map[string][]byte {
	t.Helper()
	files := make(map[string][]byte)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if d.IsDir() {
			files[rel+"/"] = nil
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		files[rel] = data
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}

func compareTrees(t *testing.T, want, got string) {
	t.Helper()
	wantFiles := readTree(t, want)
	gotFiles := readTree(t, got)
	if len(wantFiles) != len(gotFiles) {
		t.Fatalf("extracted %d entries, want %d", len(gotFiles), len(wantFiles))
	}
	for name, data := range wantFiles {
		gotData, ok := gotFiles[name]
		if !ok {
			t.Errorf("%s is missing after extraction", name)
			continue
		}
		if !bytes.Equal(data, gotData) {
			t.Errorf("%s differs after extraction", name)
		}
	}
}

func writeFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestRoundTripDirectory(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "src")

	random := make([]byte, 1<<20)
	if _, err := rand.Read(random); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(src, "hello.txt"), []byte("hello\n"))
	writeFile(t, filepath.Join(src, "empty.txt"), nil)
	writeFile(t, filepath.Join(src, "sub", "random.bin"), random)
	writeFile(t, filepath.Join(src, "sub", "deeper", "repeated.txt"), bytes.Repeat([]byte("fastzip "), 10000))
	if err := os.MkdirAll(filepath.Join(src, "empty-dir"), 0o755); err != nil {
		t.Fatal(err)
	}

	archive := filepath.Join(tmp, "archive.zip")
	out := filepath.Join(tmp, "out")
	run(t, "compress", src, archive)
	run(t, "extract", archive, out)

	compareTrees(t, src, filepath.Join(out, "src"))
}

func TestRoundTripFile(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "single.txt")
	writeFile(t, src, []byte("just one file\n"))

	archive := filepath.Join(tmp, "archive.zip")
	out := filepath.Join(tmp, "out")
	run(t, "compress", src, archive)
	run(t, "extract", archive, out)

	got, err := os.ReadFile(filepath.Join(out, "single.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "just one file\n" {
		t.Errorf("single.txt = %q after extraction", got)
	}
}
