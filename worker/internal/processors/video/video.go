package video

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/nakshatraraghav/transcodex/worker/config"
)

type VideoProcessor struct {
	data *bytes.Buffer
}

func NewVideoProcessor() *VideoProcessor {
	return &VideoProcessor{}
}

func (vp *VideoProcessor) LoadData() error {
	key := config.GetEnv().OBJECT_KEY

	fpath := filepath.Join("assets", filepath.Base(key))

	file, err := os.Open(fpath)
	if err != nil {
		return err
	}

	buf, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	vp.data = bytes.NewBuffer(buf)

	return nil

}

func (vp *VideoProcessor) ApplyTransformations(map[string]string) ([]byte, error) {
	return nil, nil
}
