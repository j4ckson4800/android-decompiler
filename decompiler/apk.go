package decompiler

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/j4ckson4800/android-decompiler/decompiler/smali"
	"github.com/j4ckson4800/android-decompiler/decompiler/smali/resource"
)

var ErrApkNotFoundInXapk = errors.New("apk not found in xapk")

type Apk struct {
	ManifestXML string
	Dexes       []smali.Dex
	Resources   resource.Table

	cfg smali.Config
}

func NewApkFromZip(r *zip.Reader, opts ...Option) (*Apk, error) {
	cfg := ParseConfig{}

	for _, opt := range opts {
		opt(&cfg)
	}

	var err error
	if !hasDexAndManifest(r) {
		r, err = extractApkFromXapk(r)
		if err != nil {
			return nil, fmt.Errorf("extract apk from xapk: %w", err)
		}
	}

	apk := &Apk{
		cfg: smali.Config{
			SanitizeAnnotations: cfg.SanitizeAnnotations,
		},
	}
	for _, file := range r.File {
		if strings.HasSuffix(file.Name, "AndroidManifest.xml") {
			if err := apk.readManifest(file); err != nil {
				return nil, fmt.Errorf("read manifest: %w", err)
			}
		}
		if strings.HasSuffix(file.Name, ".dex") {
			if err := apk.readDex(file); err != nil {
				if !cfg.FailOnInvalidDex {
					continue
				}
				return nil, fmt.Errorf("read dex: %w", err)
			}
		}
		if strings.HasSuffix(file.Name, ".arsc") {
			if err := apk.readResourceFile(file); err != nil {
				if !cfg.FailOnInvalidResource {
					continue
				}
				return nil, fmt.Errorf("read resource file: %w", err)
			}
		}
	}

	return apk, nil
}

func NewApk(reader io.ReaderAt, size int64, opts ...Option) (*Apk, error) {
	cfg := ParseConfig{}

	for _, opt := range opts {
		opt(&cfg)
	}

	r, err := zip.NewReader(reader, size)
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}

	apk, err := NewApkFromZip(r, opts...)
	if err != nil {
		return nil, fmt.Errorf("new apk from zip: %w", err)
	}

	return apk, nil
}

func hasDexAndManifest(r *zip.Reader) bool {
	hasDex := false
	hasManifest := false
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".dex") {
			hasDex = true
		}
		if strings.HasSuffix(f.Name, "AndroidManifest.xml") {
			hasManifest = true
		}

		if hasDex && hasManifest {
			return true
		}
	}
	return false
}

func extractApkFromXapk(r *zip.Reader) (*zip.Reader, error) {
	for _, file := range r.File {

		if !strings.HasSuffix(file.Name, ".apk") {
			continue
		}

		baseName := filepath.Base(file.Name)
		if strings.HasPrefix(baseName, "config.") { // skip config apks
			continue
		}

		if strings.Count(baseName, ".") == 1 {
			continue
		}

		reader, err := func() (*zip.Reader, error) {
			rc, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("open: %w", err)
			}
			defer rc.Close()

			buf := new(bytes.Buffer)
			if _, err := buf.ReadFrom(rc); err != nil {
				return nil, fmt.Errorf("read from: %w", err)
			}

			zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
			if err != nil {
				return nil, fmt.Errorf("open zip: %w", err)
			}
			return zr, nil
		}()
		if err != nil {
			return nil, fmt.Errorf("open zip reader: %w", err)
		}
		return reader, nil
	}
	return nil, ErrApkNotFoundInXapk
}

func (a *Apk) readDex(file *zip.File) error {
	rc, err := file.Open()
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer rc.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(rc); err != nil {
		return fmt.Errorf("read from: %w", err)
	}

	dex, err := smali.NewDex(bytes.NewReader(buf.Bytes()), a.cfg)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}
	dex.Filename = file.Name

	a.Dexes = append(a.Dexes, dex)
	return nil
}

func (a *Apk) readManifest(file *zip.File) error {
	rc, err := file.Open()
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer rc.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(rc); err != nil {
		return fmt.Errorf("read from: %w", err)
	}

	// NOTE: it's in binary xml format, decode it later
	// https://github.com/google/agi/tree/main/core/os/android/binaryxml
	a.ManifestXML = buf.String()
	return nil
}

func (a *Apk) readResourceFile(file *zip.File) error {
	rc, err := file.Open()
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer rc.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(rc); err != nil {
		return fmt.Errorf("read from: %w", err)
	}

	parser := smali.NewParser(bytes.NewReader(buf.Bytes()))
	table, err := resource.NewTable(parser)
	if err != nil {
		return fmt.Errorf("new table: %w", err)
	}

	a.Resources = table
	return nil
}
