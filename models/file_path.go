package models

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"gopkg.in/h2non/filetype.v1"
	"github.com/jinzhu/gorm"
)

// FilePath represents a known file
type FilePath struct {
	gorm.Model
	Realpath    string `gorm:"size:4096;not null;unique;unique_index"`
	Basename    string `gorm:"size:4096"`
	Dirname     string `gorm:"size:4096;not null;index"`
	MimeType    string `gorm:"size:15;not null;index"`
	MimeSubType string `gorm:"size:100;not null;index"`
	Sha256      string `gorm:"not null;index"`
	BytesSize   int64
	ModTime     time.Time
}

// NewFilePath FilePath from a given path
func NewFilePath(pathname string) (fp *FilePath, err error) {
	realpath, err := filepath.Abs(pathname)
	if err != nil {
		return
	}

	fp = &FilePath{Realpath: realpath}

	if err = fp.Stat(); err != nil {
		return
	}

	if err = fp.SetSha256(); err != nil {
		return
	}

	err = fp.SetMimeType()

	return
}

// SetMimeType sets MimeType and MimeSubType fields from file
func (fp *FilePath) SetMimeType() error {
	// Read a file
	kind, err := filetype.MatchFile(fp.Realpath)

	if err != nil {
		return nil
	}

	fp.MimeType = kind.MIME.Type
	fp.MimeSubType = kind.MIME.Subtype

	return nil
}

func (fp *FilePath) Stat() error {
	fileinfo, err := os.Stat(fp.Realpath)

	if err != nil {
		return err
	}

	fp.Basename = fileinfo.Name()
	fp.Dirname = path.Dir(fp.Realpath)
	fp.BytesSize = fileinfo.Size()
	fp.ModTime = fileinfo.ModTime()

	return nil
}

func (fp *FilePath) Open() (f *os.File, err error) {
	f, err = os.Open(fp.Realpath)

	return
}

func (fp *FilePath) SetSha256() error {
	hasher := sha256.New()
	f, err := fp.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := io.Copy(hasher, f); err != nil {
		return err
	}

	fp.Sha256 = hex.EncodeToString(hasher.Sum(nil))
	return nil
}
