package img

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"

	"github.com/codeblanche/golibs/img/internal/resize"
	"github.com/codeblanche/golibs/img/internal/watermark"
)

// WatermarkOptions options for watermarking the image
type WatermarkOptions watermark.Options

// ResizeOptions options for resizing the image
type ResizeOptions resize.Options

// I image type
type I struct {
	image.Image
	path     string
	encoding string
}

// New creates a new image
func New() *I {
	return &I{}
}

// Load loads an image from the given io.Reader
func (i *I) Load(r io.Reader) (err error) {
	i.Image, i.encoding, err = image.Decode(r)
	return err
}

// LoadFile loads an image from file
func (i *I) LoadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	i.Image, i.encoding, err = image.Decode(f)
	return nil
}

// Resize resizes the image to the given w(idth) and h(eight)
func (i *I) Resize(w, h uint, o ResizeOptions) *I {
	r := resize.New(i.Image)
	return &I{
		Image:    r.Resize(w, h, resize.Options(o)),
		path:     i.path,
		encoding: i.encoding,
	}
}

// Watermark adds a watermark (w) to the image
func (i *I) Watermark(w *I, o *WatermarkOptions) *I {
	return &I{
		Image:    i.Image,
		path:     i.path,
		encoding: i.encoding,
	}
}

// Save the image using the same encoding as the original file
func (i *I) Save(path string) error {
	var err error
	switch i.encoding {
	case "jpeg":
		err = i.SaveJPG(path, -1)
	case "png":
		err = i.SavePNG(path)
	case "gif":
		err = i.SaveGIF(path, 256)
	}
	if err != nil {
		return err
	}
	return nil
}

// SaveJPG saves the image as a JPEG
func (i *I) SaveJPG(path string, quality int) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	i.path = path
	i.encoding = "jpeg"
	var opts *jpeg.Options
	if quality > 0 {
		opts = &jpeg.Options{Quality: quality}
	}
	return jpeg.Encode(f, i.Image, opts)
}

// SavePNG saves the image as a PNG
func (i *I) SavePNG(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	i.path = path
	i.encoding = "png"
	return png.Encode(f, i.Image)
}

// SaveGIF saves the image as a GIF
func (i *I) SaveGIF(path string, colors int) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	i.path = path
	i.encoding = "gif"
	var opts *gif.Options
	if colors > 0 {
		opts = &gif.Options{NumColors: colors}
	}
	return gif.Encode(f, i.Image, opts)
}
