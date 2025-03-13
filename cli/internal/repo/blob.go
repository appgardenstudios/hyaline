package repo

import (
	"io"

	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetBlobBytes(blob object.Blob) (bytes []byte, err error) {
	r, err := blob.Reader()
	if err != nil {
		return
	}
	bytes, err = io.ReadAll(r)
	return
}
