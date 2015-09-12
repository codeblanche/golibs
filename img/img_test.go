package img

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/color"
	"io"
	"os"
	"path"
	"testing"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	tmp string
)

func testPhoto() io.Reader {
	return base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(TestPhoto))
}

func testWatermark() io.Reader {
	return base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(TestWatermarkImage))
}

func TestMain(m *testing.M) {
	// Set up
	tmp = path.Join(os.TempDir(), "go/tests/golibs/img", uuid.New())
	os.MkdirAll(tmp, 0755)

	fmt.Println(tmp)

	// Run tests
	result := m.Run()

	// Tear down
	// os.RemoveAll(tmp)

	// Exit
	os.Exit(result)
}

func TestLoad(t *testing.T) {
	assert := assert.New(t)
	i := New()
	assert.Nil(i.Image)
	err := i.Load(testPhoto())
	assert.NoError(err)
	assert.NotNil(i.Image)
}

func TestLoadFile(t *testing.T) {
	assert := assert.New(t)
	path := path.Join(tmp, "test-load-file.png")
	f, err := os.Create(path)
	assert.NoError(err)
	io.Copy(f, testPhoto())
	f.Close()
	i := New()
	err = i.LoadFile(path)
	assert.NoError(err)
}

func TestSave(t *testing.T) {
	assert := assert.New(t)
	i := New()
	err := i.Load(testPhoto())
	assert.NoError(err)
	path := path.Join(tmp, "test-save.png")
	err = i.Save(path)
	assert.NoError(err)
	_, err = os.Stat(path)
	assert.NoError(err)
}

func TestSaveJPG(t *testing.T) {
	assert := assert.New(t)
	i := New()
	err := i.Load(testPhoto())
	assert.NoError(err)
	path := path.Join(tmp, "test-save.jpg")
	err = i.SaveJPG(path, 20)
	assert.NoError(err)
	_, err = os.Stat(path)
	assert.NoError(err)
}

func TestSavePNG(t *testing.T) {
	assert := assert.New(t)
	i := New()
	err := i.Load(testPhoto())
	assert.NoError(err)
	path := path.Join(tmp, "test-save.png")
	err = i.SavePNG(path)
	assert.NoError(err)
	_, err = os.Stat(path)
	assert.NoError(err)
}

func TestSaveGIF(t *testing.T) {
	assert := assert.New(t)
	i := New()
	err := i.Load(testPhoto())
	assert.NoError(err)
	path := path.Join(tmp, "test-save.gif")
	err = i.SaveGIF(path, 20)
	assert.NoError(err)
	_, err = os.Stat(path)
	assert.NoError(err)
}

func TestResize(t *testing.T) {
	assert := assert.New(t)
	o := New()
	err := o.Load(testPhoto())
	assert.NoError(err)
	options := ResizeOptions{
		MaintainRatio: true,
		Crop:          true,
		Envelope:      true,
		Color:         color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
	}
	i := o.Resize(300, 300, options)
	assert.Equal(300, i.Image.Bounds().Dx())
	assert.Equal(300, i.Image.Bounds().Dy())
	path := path.Join(tmp, "test-resize-300.jpg")
	err = i.SaveJPG(path, 100)
	assert.NoError(err)
}

func TestWatermark(t *testing.T) {
	assert := assert.New(t)

	_ = assert
}
