package selfupdate

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompressionNotRequired(t *testing.T) {
	buf := []byte{'a', 'b', 'c'}
	want := bytes.NewReader(buf)
	r, err := uncompress(want, "https://github.com/foo/bar/releases/download/v1.2.3/foo", "foo")
	if err != nil {
		t.Fatal(err)
	}
	have, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	for i, b := range have {
		if buf[i] != b {
			t.Error(i, "th elem is not the same as wanted. want", buf[i], "but got", b)
		}
	}
}

func TestUncompress(t *testing.T) {
	for _, n := range []string{
		"testdata/foo.zip",
		"testdata/single-file.zip",
		"testdata/single-file.gz",
		"testdata/single-file.gzip",
		"testdata/foo.tar.gz",
	} {
		t.Run(n, func(t *testing.T) {
			f, err := os.Open(n)
			if err != nil {
				t.Fatal(err)
			}

			var ext string
			if strings.HasSuffix(n, ".tar.gz") {
				ext = ".tar.gz"
			} else {
				ext = filepath.Ext(n)
			}

			url := "https://github.com/foo/bar/releases/download/v1.2.3/bar" + ext
			r, err := uncompress(f, url, "bar")
			if err != nil {
				t.Fatal(err)
			}

			bytes, err := ioutil.ReadAll(r)
			if err != nil {
				t.Fatal(err)
			}
			s := string(bytes)
			if s != "this is test\n" {
				t.Fatal("Uncompressing zip failed into unexpected content", s)
			}
		})
	}
}

func TestUncompressInvalidArchive(t *testing.T) {
	for _, a := range []struct {
		name string
		msg  string
	}{
		{"testdata/invalid.zip", "not a valid zip file"},
		{"testdata/invalid.gz", "Failed to uncompress gzip file"},
		{"testdata/invalid-tar.tar.gz", "Failed to unarchive .tar file"},
		{"testdata/invalid-gzip.tar.gz", "Failed to uncompress .tar.gz file"},
	} {
		f, err := os.Open(a.name)
		if err != nil {
			t.Fatal(err)
		}

		var ext string
		if strings.HasSuffix(a.name, ".tar.gz") {
			ext = ".tar.gz"
		} else {
			ext = filepath.Ext(a.name)
		}

		url := "https://github.com/foo/bar/releases/download/v1.2.3/bar" + ext
		_, err = uncompress(f, url, "bar")
		if err == nil {
			t.Fatal("Error should be raised")
		}
		if !strings.Contains(err.Error(), a.msg) {
			t.Fatal("Unexpected error:", err)
		}
	}
}

func TestTargetNotFoundInZip(t *testing.T) {
	for _, f := range []string{
		"testdata/empty.zip",
		"testdata/bar-not-found.zip",
	} {
		t.Run(f, func(t *testing.T) {
			f, err := os.Open(f)
			if err != nil {
				t.Fatal(err)
			}

			_, err = uncompress(f, "https://github.com/foo/bar/releases/download/v1.2.3/bar.zip", "bar")
			if err == nil {
				t.Fatal("Error should be raised for")
			}
			if !strings.Contains(err.Error(), "command is not found") {
				t.Fatal("Unexpected error:", err)
			}
		})
	}
}

func TestTargetNotFoundInGZip(t *testing.T) {
	f, err := os.Open("testdata/bar-not-found.gzip")
	if err != nil {
		t.Fatal(err)
	}

	_, err = uncompress(f, "https://github.com/foo/bar/releases/download/v1.2.3/bar.gzip", "bar")
	if err == nil {
		t.Fatal("Error should be raised for")
	}
	if !strings.Contains(err.Error(), "does not match to command") {
		t.Fatal("Unexpected error:", err)
	}
}

func TestTargetNotFoundInTarGz(t *testing.T) {
	for _, f := range []string{
		"testdata/empty.tar.gz",
		"testdata/bar-not-found.tar.gz",
	} {
		t.Run(f, func(t *testing.T) {
			f, err := os.Open(f)
			if err != nil {
				t.Fatal(err)
			}

			_, err = uncompress(f, "https://github.com/foo/bar/releases/download/v1.2.3/bar.tar.gz", "bar")
			if err == nil {
				t.Fatal("Error should be raised for")
			}
			if !strings.Contains(err.Error(), "command is not found") {
				t.Fatal("Unexpected error:", err)
			}
		})
	}
}
