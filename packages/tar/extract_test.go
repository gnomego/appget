package tar_test

import (
	"testing"

	"github.com/gnomego/avm/packages/tar"
)

func TestExt(t *testing.T) {
	tests := []struct {
		src string
		ext string
	}{
		{"foo.tar", ".tar"},
		{"foo.tar.gz", ".tar.gz"},
		{"foo.tgz", ".tgz"},
		{"foo.tar.bz2", ".tar.bz2"},
		{"foo.tbz2", ".tbz2"},
		{"foo.tar.xz", ".tar.xz"},
		{"foo.txz", ".txz"},
		{"foo.tar.zst", ""},
	}

	for _, tt := range tests {
		if got := tar.GetTarExtension(tt.src); got != tt.ext {
			t.Errorf("getTarExtension(%q) = %q; want %q", tt.src, got, tt.ext)
		}
	}
}
