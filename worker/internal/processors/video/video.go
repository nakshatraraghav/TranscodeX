package video

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nakshatraraghav/transcodex/worker/config"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type fn func(string) error

type VideoProcessor struct {
	path string
}

func NewVideoProcessor() *VideoProcessor {
	key := config.GetEnv().OBJECT_KEY
	fpath := filepath.Join("assets", filepath.Base(key))

	return &VideoProcessor{path: fpath}
}

func (vp *VideoProcessor) ApplyTransformations(tmap map[string]string) []error {
	var e []error

	fmap := map[string]fn{
		"TRANSCODE":            vp.TranscodeToMultipleResolutions,
		"TRANSCODE-RESOLUTION": vp.TranscodeToResolution,
	}

	for key, value := range tmap {
		if action, exists := fmap[key]; exists {
			err := action(value)
			if err != nil {
				e = append(e, err)
			}
		} else {
			e = append(e, fmt.Errorf("%v transformation isnt supported, and not applied", key))
		}
	}

	return e
}

func Transcode(path, resolution string) error {
	opath := ReturnNewPath(path, resolution)

	err := ffmpeg.Input(path).
		Filter("scale", ffmpeg.Args{resolution}).
		Output(opath,
			ffmpeg.KwArgs{
				"c:v":    "libx264",
				"preset": "fast",
				"crf":    "23",
			},
			ffmpeg.KwArgs{
				"c:a": "aac",
				"b:a": "128k",
			}).
		OverWriteOutput().
		Run()

	if err != nil {
		return fmt.Errorf("error transcoding video: %w", err)
	}

	return nil
}

func (vp *VideoProcessor) TranscodeToResolution(resolution string) error {
	return Transcode(vp.path, resolution)
}

func (vp *VideoProcessor) TranscodeToMultipleResolutions(parameter string) error {
	fmt.Println("PATH:", vp.path)

	var wg sync.WaitGroup

	rmap := map[string]string{
		"360p":  "640x360",
		"720p":  "1280x720",
		"1080p": "1920x1080",
		"1440p": "2560x1440",
	}

	var errs []error

	for name, resolution := range rmap {
		wg.Add(1)
		go func(name, resolution string) {
			defer wg.Done()

			err := Transcode(vp.path, resolution)
			if err != nil {
				errs = append(errs, fmt.Errorf("resolution %s: %w", name, err))
			}
		}(name, resolution)
	}

	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("transcoding errors: %v", errs)
	}

	return nil
}

func ReturnNewPath(original, resolution string) string {
	f := strings.Split(filepath.Base(original), ".")
	name := f[0]
	extension := f[1]

	nname := name + "_" + resolution + "." + extension

	opath := filepath.Join("assets", nname)

	return opath
}
