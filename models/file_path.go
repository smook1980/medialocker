package models

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/smook1980/medialocker/types"
	"gopkg.in/h2non/filetype.v1"
)

// FilePath represents a known file
type FilePath struct {
	Model
	Basename    string `gorm:"size:4096;not null"`
	Dirname     string `gorm:"size:4096;not null"`
	MimeType    string `gorm:"size:15;not null;index"`
	MimeSubType string `gorm:"size:100;not null;index"`
	HashURI     string `gorm:"not null;index"`
	BytesSize   int64
	ModTime     time.Time
}

// NewFilePath creates a new FilePath and sets it's fields for the given path.
// hashURI is optional, an empty string causes the value to be updated.
func NewFilePath(pathname, hashURI string) (fp *FilePath, err error) {
	fp = &FilePath{}

	if err = fp.setPath(pathname); err != nil {
		return
	}

	if err = fp.Stat(); err != nil {
		// Warn
	}

	if hashURI == "" {
		if err = fp.calcHash(); err != nil {
			// Warn
		}
	} else {
		fp.HashURI = hashURI

	}

	err = fp.SetMimeType()

	return
}

func (fp *FilePath) setPath(pathname string) error {
	realpath, err := filepath.Abs(pathname)
	if err != nil {
		return err
	}

	fp.Basename = filepath.Base(realpath)
	fp.Dirname = filepath.Dir(realpath)

	return nil
}

// Realpath returns the absolute path
func (fp *FilePath) Realpath() string {
	return filepath.Join(fp.Dirname, fp.Basename)
}

// SetMimeType sets MimeType and MimeSubType fields from file
func (fp *FilePath) SetMimeType() error {
	// Read a file
	kind, err := filetype.MatchFile(fp.Realpath())

	if err != nil {
		return nil
	}

	fp.MimeType = kind.MIME.Type
	fp.MimeSubType = kind.MIME.Subtype

	return nil
}

func (fp *FilePath) Stat() error {
	fileinfo, err := os.Stat(fp.Realpath())

	if err != nil {
		return err
	}

	fp.Basename = fileinfo.Name()
	fp.Dirname = path.Dir(fp.Realpath())
	fp.BytesSize = fileinfo.Size()
	fp.ModTime = fileinfo.ModTime()

	return nil
}

// Open opens the file for reading.
func (fp *FilePath) Open() (f *os.File, err error) {
	f, err = os.Open(fp.Realpath())

	return
}

func (fp *FilePath) calcHash() error {
	f, err := fp.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	hash, err := types.HashReader(types.Blake2bHasher, f)
	if err != nil {
		return err
	}

	fp.HashURI = hash.URI()

	return nil
}

// Hash for file
func (fp *FilePath) Hash() (*types.Hash, error) {
	h, err := types.HashFromURI(fp.HashURI)
	if err != nil {
		if err := fp.calcHash(); err != nil {
			return &types.Hash{}, err
		}

		return types.HashFromURI(fp.HashURI)
	}

	return h, nil
}

func init() {
	registerSchemaMigrator(func(db *gorm.DB) {
		fmt.Println("FILE PATH SCHEMA SHIT YO")
		m := &FilePath{}

		db.
			AutoMigrate(m).
			Model(m).
			AddUniqueIndex("idx_basename_dirname", "basename", "dirname")
	})
}
