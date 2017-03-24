package types

import (
	"io"
	"os"

	"github.com/smook1980/medialocker/util"

	filetype "gopkg.in/h2non/filetype.v1"
)

type MediaPath struct {
	Realpath string
	Hash     string
	MimeType string
	Type     MediaType
	os.FileInfo
}

func (mp *MediaPath) Update() error {
	file, err := os.Open(mp.Realpath)
	if err != nil {
		return err
	}
	defer file.Close()

	err1 := mp.updateType(file)
	file.Seek(0, 0)
	err2 := mp.updateHash(file)

	if err1 != nil || err2 != nil {
		return util.MultiError(err1, err2)
	}

	return nil
}

func (mp *MediaPath) updateHash(in io.Reader) error {
	hash, err := HashReader(Blake2bHasher, in)

	mp.Hash = hash.URI()
	return err
}

func (mp *MediaPath) updateType(in io.Reader) error {
	mime, mtype, err := pathMediaType(in)

	if err != nil {
		return err
	}

	mp.MimeType = mime
	mp.Type = mtype

	return nil
}

func pathMediaType(file io.Reader) (string, MediaType, error) {
	buffer := make([]byte, 512)

	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", Unknown, err
	}

	if t, err := filetype.Image(buffer); err == nil {
		return t.MIME.Value, Image, nil
	}

	if t, err := filetype.Video(buffer); err == nil {
		return t.MIME.Value, Video, nil
	}

	if t, err := filetype.Archive(buffer); err == nil {
		return t.MIME.Value, Archive, nil
	}

	return "", Unknown, nil
}
