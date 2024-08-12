package image

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/h2non/bimg"
	"github.com/nakshatraraghav/transcodex/worker/config"
)

type fn func(string) error

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
	defer file.Close() // Ensure file is closed before attempting to delete

	buf, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	ip.data = bytes.NewBuffer(buf)

	// Close the file before attempting to delete it
	file.Close()

	// Attempt to delete the file after closing it
	err = os.Remove(fpath)
	if err != nil {
		return err
	}

	return nil
}

func (ip *ImageProcessor) ApplyTransformations(tmap map[string]string) []error {

	ip.LoadData()

	var e []error = []error{}

	fmap := map[string]fn{
		"RESIZE":             ip.Resize,
		"FORCE-RESIZE":       ip.ForceResize,
		"ROTATE":             ip.Rotate,
		"CONVERT-FORMAT":     ip.ConvertFormat,
		"WATERMARK":          ip.Watermark,
		"GENERATE-THUMBNAIL": ip.GenerateThumbnail,
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

	ip.SaveChanges()

	return e
}

func (ip *ImageProcessor) Resize(parameter string) error {
	parts := strings.Split(parameter, "x")
	if len(parts) != 2 {
		return errors.New("invalid image resize parameters")
	}

	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}
	y, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	img, err := bimg.NewImage(ip.data.Bytes()).Resize(x, y)
	if err != nil {
		return err
	}

	ip.data = bytes.NewBuffer(img)
	return nil
}

func (ip *ImageProcessor) ForceResize(parameter string) error {
	parts := strings.Split(parameter, "x")
	if len(parts) != 2 {
		return errors.New("invalid image resize parameters")
	}

	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}
	y, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	img, err := bimg.NewImage(ip.data.Bytes()).ForceResize(x, y)
	if err != nil {
		return err
	}

	ip.data = bytes.NewBuffer(img)
	return nil
}

func (ip *ImageProcessor) Rotate(parameter string) error {
	angle, err := strconv.Atoi(parameter)
	if err != nil {
		return err
	}

	img, err := bimg.NewImage(ip.data.Bytes()).Rotate(bimg.Angle(angle))
	if err != nil {
		return err
	}

	ip.data = bytes.NewBuffer(img)

	return nil
}

func (ip *ImageProcessor) ConvertFormat(parameter string) error {
	if parameter != "PNG" && parameter != "JPEG" && parameter != "SVG" && parameter != "WEBP" {
		return errors.New("invalid file conversion format")
	}

	tmap := map[string]bimg.ImageType{
		"PNG":  bimg.PNG,
		"JPEG": bimg.JPEG,
		"SVG":  bimg.SVG,
		"WEBP": bimg.WEBP,
	}

	itype, ok := tmap[parameter]
	if !ok {
		return errors.New("unknown image type")
	}

	img, err := bimg.NewImage(ip.data.Bytes()).Convert(itype)
	if err != nil {
		return fmt.Errorf("image conversion failed: %w", err)
	}

	ip.data = bytes.NewBuffer(img)

	fmt.Println(bimg.DetermineImageType(img) == bimg.WEBP)

	return nil
}

func (ip *ImageProcessor) Watermark(parameter string) error {

	w := bimg.Watermark{
		Text:       parameter,
		Opacity:    0.75,
		Width:      300,
		DPI:        300,
		Margin:     150,
		Font:       "sans bold 12",
		Background: bimg.Color{R: 255, G: 255, B: 255},
	}

	img, err := bimg.NewImage(ip.data.Bytes()).Watermark(w)
	if err != nil {
		return err
	}

	ip.data = bytes.NewBuffer(img)

	return nil
}

func (ip *ImageProcessor) GenerateThumbnail(parameter string) error {

	pixels, err := strconv.Atoi(parameter)
	if err != nil {
		return err
	}

	img, err := bimg.NewImage(ip.data.Bytes()).Thumbnail(pixels)
	if err != nil {
		return err
	}

	ip.data = bytes.NewBuffer(img)

	return nil
}

func (ip *ImageProcessor) SaveChanges() error {

	outputDir := filepath.Join("assets", "output")
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}

	_, file := filepath.Split(ip.path)
	opath := filepath.Join(outputDir, file)

	return bimg.Write(opath, ip.data.Bytes())
}
