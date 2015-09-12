package resize

import (
	"image"
	"image/color"
	"math"

	"github.com/nfnt/resize"
)

// Color local type color.Color
type Color color.Color

// Options for resizing an image
type Options struct {
	// MaintainRatio sets whether to maintain aspect ratio when resizing
	MaintainRatio bool
	// Crop enables cropping of the image beyond the new dimensions
	Crop bool
	// Envelope enables the use of a border to fill the non-image space.
	// Use only when Crop is false.
	Envelope bool
	// Color of the Envelope when used
	Color Color
}

// Resizer interface defines a resizeable image
type Resizer interface {
	Resize(uint, uint, Options) image.Image
}

// SubImager interface defines images that support SubImage-ing
type SubImager interface {
	SubImage(image.Rectangle) image.Image
}

// Image is the internal resizeable image type
type Image struct {
	image.Image
}

// New creates a new resizeable image
func New(i image.Image) *Image {
	return &Image{Image: i}
}

// Resize the Image to to the given target width (tw) and height (th) with the provided options.
func (i *Image) Resize(tw, th uint, o Options) image.Image {
	w, h := tw, th
	img := i.Image

	if o.MaintainRatio && tw > 0 && th > 0 {
		cw, ch := float64(i.Image.Bounds().Dx()), float64(i.Image.Bounds().Dy())
		if o.Crop {
			sx, sy := float64(w)/cw, float64(h)/ch
			s := math.Max(sx, sy)
			cropw, croph := round(float64(tw)/s), round(float64(th)/s)
			offsetw, offseth := round((cw-cropw)/2), round((ch-croph)/2)
			img = crop(img.(SubImager), uint(cropw), uint(croph), uint(offsetw), uint(offseth))
		} else if o.Envelope {
			// TODO implement envelope feature
		}
	}

	return resize.Resize(w, h, img, resize.NearestNeighbor)
}

func crop(i SubImager, tw, th, ow, oh uint) image.Image {
	return i.SubImage(image.Rect(int(ow), int(oh), int(tw+ow), int(th+oh)))
}

func envelope(i image.Image, tw, th, ow, oh uint, c Color) image.Image {
	return i
}

func round(v float64) float64 {
	if v < 0 {
		return math.Ceil(v - 0.5)
	}
	return math.Floor(v + 0.5)
}
