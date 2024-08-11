package audio

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/nakshatraraghav/transcodex/worker/config"
)

type AudioProcessor struct {
	path string
	data *bytes.Buffer
}

func NewAudioProcessor() *AudioProcessor {
	return &AudioProcessor{}
}

func (ap *AudioProcessor) LoadData() error {
	key := config.GetEnv().OBJECT_KEY
	fpath := filepath.Join("assets", filepath.Base(key))
	ap.path = fpath

	file, err := os.Open(fpath)
	if err != nil {
		return err
	}

	buf, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	ap.data = bytes.NewBuffer(buf)

	return nil

}

func (ap *AudioProcessor) ApplyTransformations(map[string]string) []error {
	ap.LoadData()
	return nil
}
