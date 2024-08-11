package image

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/nakshatraraghav/transcodex/worker/config"
)

type ImageProcessor struct {
	path string
	data *bytes.Buffer
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

func (ip *ImageProcessor) LoadData() error {
	key := config.GetEnv().OBJECT_KEY
	fpath := filepath.Join("assets", filepath.Base(key))
	ip.path = fpath

	file, err := os.Open(fpath)
	if err != nil {
		return err
	}

	buf, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	ip.data = bytes.NewBuffer(buf)

	return nil

}

func (ip *ImageProcessor) ApplyTransformations(map[string]string) ([]byte, error) {
	ip.LoadData()
	return nil, nil
}
