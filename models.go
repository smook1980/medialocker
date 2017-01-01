package medialocker

import (
	"os"
	"time"
)

const (
	VIDEO = "VIDEO"
	IMAGE = "IMAGE"
)

type MediaBlob struct {
}

// MediaPath
type MediaPath struct {
	Abs      string
	Ext      string
	MimeType string
	Size     uint32
	ModTime  time.Time
	Sha256   string
}

// Open the associated file
func (fp *MediaPath) Open() (*os.File, error) {
	return os.Open(fp.Abs)
}

func (fp *MediaPath) isVideo() bool {
	return false
}

func (fp *MediaPath) isImage() bool {
	return false
}

type Video struct {
}
