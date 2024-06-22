package tar

import (
	"archive/tar"
	"compress/bzip2"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ulikunitz/xz"
)

func Extract(src, dst string) error {
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()

	switch GetTarExtension(src) {
	case ".tar":
		extractTarCore(r, dst)
	case ".tar.gz", ".tgz":
		gzr, err := gzip.NewReader(r)
		if err != nil {
			return err
		}
		defer gzr.Close()
		extractTarCore(gzr, dst)
	case ".tar.bz2", ".tbz2":
		bz2r := bzip2.NewReader(r)
		extractTarCore(bz2r, dst)
	case ".tar.xz", ".txz":
		xzr, err := xz.NewReader(r)
		if err != nil {
			return err
		}

		extractTarCore(xzr, dst)
	}
	return nil
}

func ExtractFiles(src string, files []string, dst string) error {

	r, err := os.Open(src)
	if err != nil {
		return err
	}

	defer r.Close()

	switch GetTarExtension(src) {
	case ".tar":
		return extractTarFilesCore(r, files, dst)
	case ".tar.gz", ".tgz":
		gzr, err := gzip.NewReader(r)
		if err != nil {
			return err
		}
		defer gzr.Close()
		return extractTarFilesCore(gzr, files, dst)
	case ".tar.bz2", ".tbz2":
		bz2r := bzip2.NewReader(r)
		return extractTarFilesCore(bz2r, files, dst)
	case ".tar.xz", ".txz":
		xzr, err := xz.NewReader(r)
		if err != nil {
			return err
		}

		return extractTarFilesCore(xzr, files, dst)
	}

	return nil
}

func GetTarExtension(src string) string {
	ext := filepath.Ext(src)

	switch ext {
	case ".tar", ".tgz", ".tbz2", ".txz":
		return ext

	default:
		name := filepath.Base(src)
		index := strings.Index(name, ".tar")
		if index == -1 {
			return ""
		}

		ext := name[index:]
		switch ext {
		case ".tar", ".tar.gz", ".tar.bz2", ".tar.xz":
			return ext

		default:
			return ""
		}
	}
}

func extractTarFilesCore(r io.Reader, files []string, dst string) error {
	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		for _, file := range files {
			if header.Name == file {

				target := filepath.Join(dst, header.Name)

				switch header.Typeflag {
				case tar.TypeDir:
					if err := os.MkdirAll(target, 0755); err != nil {
						return err
					}
				case tar.TypeReg:
					f, err := os.Create(target)
					if err != nil {
						return err
					}
					defer f.Close()
					if _, err := io.Copy(f, tr); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func extractTarCore(r io.Reader, dst string) error {
	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dst, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			f, err := os.Create(target)
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
		}
	}

	return nil
}
