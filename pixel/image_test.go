package pixel_test

import (
	goimage "image"
	"image/color"
	"math/rand"
	"testing"

	"tinygo.org/x/drivers/pixel"
)

func TestImageRGB565BE(t *testing.T) {
	image := pixel.NewImage[pixel.RGB565BE](5, 3)
	if width, height := image.Size(); width != 5 && height != 3 {
		t.Errorf("image.Size(): expected 5, 3 but got %d, %d", width, height)
	}
	for _, c := range []color.RGBA{
		{R: 0xff, A: 0xff},
		{G: 0xff, A: 0xff},
		{B: 0xff, A: 0xff},
		{R: 0x10, A: 0xff},
		{G: 0x10, A: 0xff},
		{B: 0x10, A: 0xff},
	} {
		image.Set(4, 2, pixel.NewColor[pixel.RGB565BE](c.R, c.G, c.B))
		c2 := image.Get(4, 2).RGBA()
		if c2 != c {
			t.Errorf("failed to roundtrip color: expected %v but got %v", c, c2)
		}
	}
}

func TestImageRGB444BE(t *testing.T) {
	image := pixel.NewImage[pixel.RGB444BE](5, 3)
	if width, height := image.Size(); width != 5 && height != 3 {
		t.Errorf("image.Size(): expected 5, 3 but got %d, %d", width, height)
	}
	for _, c := range []color.RGBA{
		{R: 0xff, A: 0xff},
		{G: 0xff, A: 0xff},
		{B: 0xff, A: 0xff},
		{R: 0x11, A: 0xff},
		{G: 0x11, A: 0xff},
		{B: 0x11, A: 0xff},
	} {
		encoded := pixel.NewColor[pixel.RGB444BE](c.R, c.G, c.B)
		image.Set(0, 0, encoded)
		image.Set(0, 1, encoded)
		encoded2 := image.Get(0, 0)
		encoded3 := image.Get(0, 1)
		if encoded != encoded2 {
			t.Errorf("failed to roundtrip color %v: expected %d but got %d", c, encoded, encoded2)
		}
		if encoded != encoded3 {
			t.Errorf("failed to roundtrip color %v: expected %d but got %d", c, encoded, encoded3)
		}
		c2 := encoded2.RGBA()
		if c2 != c {
			t.Errorf("failed to roundtrip color: expected %v but got %v", c, c2)
		}
		c3 := encoded3.RGBA()
		if c3 != c {
			t.Errorf("failed to roundtrip color: expected %v but got %v", c, c3)
		}
	}
}

func TestImageMonochrome(t *testing.T) {
	image := pixel.NewImage[pixel.Monochrome](128, 64)
	if width, height := image.Size(); width != 128 && height != 64 {
		t.Errorf("image.Size(): expected 128, 64 but got %d, %d", width, height)
	}
	for _, expected := range []color.RGBA{
		{R: 0xff, G: 0xff, B: 0xff},
		{G: 0xff},
		{R: 0xff, G: 0xff},
		{G: 0xff, B: 0xff},
		{R: 0x00},
		{G: 0x00, A: 0xff},
		{B: 0x00, A: 0xff},
	} {
		encoded := pixel.NewColor[pixel.Monochrome](expected.R, expected.G, expected.B)
		image.Set(5, 3, encoded)
		actual := image.Get(5, 3).RGBA()
		switch {
		case int(expected.R)+int(expected.G)+int(expected.B) > 128*3:
			// should be true eg white
			if actual.R == 0 || actual.G == 0 || actual.B == 0 {
				t.Errorf("failed to roundtrip color: expected %v but got %v", expected, actual)
			}
		default:
			// should be false eg black
			if actual.R != 0 || actual.G != 0 || actual.B != 0 {
				t.Errorf("failed to roundtrip color: expected %v but got %v", expected, actual)
			}
		}
	}
}

// 128x128
var rprofile = []byte{
	0x00, 0x00, 0x11, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00,
	0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0x80, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x80, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x80, 0x00,
	0x00, 0x07, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xC0, 0x03, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xC0, 0x00, 0x00, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF8, 0x0F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFC, 0x00,
	0x00, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xE8, 0x17, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF8, 0x00, 0x07, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xC0,
	0x3F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFC, 0x3F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE,
	0x3F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFC, 0x3F, 0xFF, 0xFF, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xFF, 0xFF, 0xFE,
	0x3F, 0xFF, 0xFE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3F, 0xFF, 0xFE, 0x3F, 0xFF, 0xFC, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3F, 0xFF, 0xFC,
	0x3F, 0xFF, 0xFC, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3F, 0xFF, 0xFE, 0x3F, 0xFF, 0xFC, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3F, 0xFF, 0xFE,
	0x3F, 0xFF, 0xFC, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1F, 0xFF, 0xFE, 0x3F, 0xFF, 0xFC, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3F, 0xFF, 0xFC,
	0x3F, 0xFF, 0xFC, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1F, 0xFF, 0xFE, 0x3F, 0xFF, 0xFC, 0x00, 0x00, 0x00, 0x00, 0xBF, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x3F, 0xFF, 0xFE,
	0x3F, 0xFF, 0xFC, 0x01, 0xF8, 0x00, 0x5F, 0xFF, 0xFF, 0xFC, 0x00, 0x02, 0x80, 0x1F, 0xFF, 0xFE, 0x3F, 0xFF, 0xFC, 0x03, 0xFE, 0x03, 0xFF, 0xFD, 0xBF, 0xFF, 0x80, 0x1F, 0xE0, 0x3F, 0xFF, 0xFC,
	0x3F, 0xFF, 0xFC, 0x03, 0xFE, 0x01, 0xFF, 0xF7, 0x6B, 0xFF, 0x80, 0x1F, 0xC0, 0x1F, 0xFF, 0xFE, 0x02, 0xFF, 0xFC, 0x03, 0xDF, 0x17, 0xFA, 0x00, 0x00, 0x37, 0xF0, 0x3F, 0xE0, 0x1F, 0xFF, 0xD0,
	0x00, 0x07, 0xFC, 0x07, 0x07, 0xBF, 0x00, 0x00, 0x00, 0x01, 0xFE, 0xF8, 0x78, 0x3F, 0xF8, 0x00, 0x00, 0x01, 0xFC, 0x06, 0x1B, 0xFC, 0x00, 0x00, 0x00, 0x00, 0x3F, 0xEA, 0x78, 0x1F, 0x80, 0x00,
	0x00, 0x01, 0xFC, 0x03, 0x1B, 0xFE, 0x00, 0x00, 0x00, 0x00, 0x3F, 0xE2, 0x68, 0x1F, 0xC0, 0x00, 0x00, 0x01, 0xFC, 0x07, 0x5B, 0xE8, 0x00, 0x00, 0x00, 0x00, 0x07, 0xC4, 0x38, 0x1F, 0x80, 0x00,
	0x00, 0x01, 0xF8, 0x03, 0x3F, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0xF7, 0x78, 0x3F, 0xC0, 0x00, 0x00, 0x01, 0xFC, 0x03, 0x1F, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x7C, 0x68, 0x1F, 0x80, 0x00,
	0x00, 0x01, 0xFC, 0x07, 0x9F, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFC, 0x70, 0x1F, 0x80, 0x00, 0x00, 0x01, 0xFC, 0x03, 0x9E, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x3F, 0x3E, 0xE8, 0x1F, 0xC0, 0x00,
	0x00, 0x01, 0xFC, 0x01, 0xFF, 0xFF, 0xE0, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xE0, 0x1F, 0x80, 0x00, 0x00, 0x01, 0xF8, 0x01, 0xFF, 0x55, 0xF8, 0x00, 0x00, 0x03, 0xEA, 0xFF, 0xC0, 0x1F, 0xC0, 0x00,
	0x00, 0x01, 0xFC, 0x01, 0xFE, 0xAF, 0xF0, 0x00, 0x00, 0x03, 0xFD, 0xBF, 0xE0, 0x1F, 0x80, 0x00, 0x00, 0x01, 0xFC, 0x01, 0xF8, 0x00, 0x7C, 0x00, 0x00, 0x07, 0x80, 0x03, 0xE0, 0x1F, 0xC0, 0x00,
	0x00, 0x01, 0xFC, 0x03, 0xC1, 0xE0, 0x3E, 0x00, 0x00, 0x1F, 0x03, 0xC0, 0xF0, 0x1F, 0x80, 0x00, 0x00, 0x00, 0xF8, 0x07, 0x82, 0xF8, 0x0F, 0x00, 0x00, 0x3E, 0x05, 0xE0, 0x7C, 0x1F, 0x80, 0x00,
	0x00, 0x01, 0xFC, 0x07, 0x01, 0xF8, 0x17, 0x00, 0x00, 0x3C, 0x09, 0xE0, 0x78, 0x1F, 0xC0, 0x00, 0x00, 0x03, 0xFC, 0x0F, 0x06, 0xFC, 0x07, 0x80, 0x00, 0x7C, 0x1D, 0xE0, 0x3C, 0x1F, 0xC0, 0x00,
	0x00, 0xFF, 0xFC, 0x1E, 0x03, 0xFC, 0x03, 0x80, 0x00, 0x78, 0x1F, 0xE0, 0x1E, 0x1F, 0xFE, 0x80, 0x2F, 0xFF, 0xFC, 0x1C, 0x07, 0xF8, 0x01, 0x80, 0x00, 0x60, 0x1F, 0xE0, 0x16, 0x1F, 0xFF, 0xF8,
	0x1F, 0xFF, 0xF8, 0x3E, 0x07, 0xFC, 0x01, 0xC0, 0x00, 0xE8, 0x1F, 0xE0, 0x1E, 0x1F, 0xFF, 0xEC, 0x3F, 0xFF, 0xFC, 0x3C, 0x03, 0xF8, 0x01, 0xC0, 0x00, 0xE0, 0x17, 0xE0, 0x07, 0x1F, 0xFF, 0xFE,
	0x3F, 0xFF, 0xFC, 0x38, 0x03, 0xE8, 0x00, 0xC0, 0x00, 0xC0, 0x07, 0xC0, 0x03, 0x1F, 0xFF, 0xFC, 0x3F, 0xFF, 0xFC, 0x38, 0x00, 0xC0, 0x00, 0xE0, 0x00, 0xC0, 0x00, 0x00, 0x03, 0x9F, 0xFF, 0xFE,
	0x3F, 0xFF, 0xFC, 0x78, 0x01, 0xE0, 0x00, 0xC0, 0x00, 0xC0, 0x00, 0x00, 0x03, 0x1F, 0xFF, 0xFC, 0x3F, 0xFF, 0xFC, 0x38, 0x00, 0x00, 0x00, 0xE0, 0x01, 0xC0, 0x00, 0x00, 0x03, 0x9F, 0xFF, 0xFE,
	0x3F, 0xFF, 0xF8, 0x78, 0x00, 0x00, 0x00, 0xE0, 0x01, 0xC0, 0x00, 0x00, 0x03, 0x9F, 0xFF, 0xFE, 0x3F, 0xFF, 0xFC, 0x68, 0x00, 0x00, 0x00, 0xE0, 0x01, 0xC0, 0x00, 0x00, 0x03, 0x9F, 0xFF, 0xFC,
	0x3F, 0xFF, 0xFC, 0x70, 0x00, 0x00, 0x00, 0xE0, 0x01, 0xC0, 0x00, 0x00, 0x01, 0x9F, 0xFF, 0xFE, 0x3F, 0xFF, 0xFC, 0x78, 0x00, 0x00, 0x00, 0xE0, 0x01, 0xC0, 0x00, 0x00, 0x03, 0x9F, 0xFF, 0xFE,
	0x3F, 0xFF, 0xFC, 0x68, 0x00, 0x00, 0x00, 0xE0, 0x01, 0xC0, 0x00, 0x00, 0x03, 0x9F, 0xFF, 0xFE, 0x3F, 0xFF, 0xF8, 0x70, 0x00, 0x00, 0x00, 0xC0, 0x01, 0xC0, 0x00, 0x00, 0x01, 0x9F, 0xFF, 0xFC,
	0x3F, 0xFF, 0xFC, 0x68, 0x00, 0x00, 0x00, 0xE0, 0x01, 0xC0, 0x00, 0x00, 0x03, 0x9F, 0xFF, 0xFE, 0x3F, 0xFF, 0xFC, 0x38, 0x00, 0x00, 0x00, 0xC0, 0x00, 0xC0, 0x00, 0x00, 0x03, 0x9F, 0xFF, 0xFE,
	0x07, 0xFF, 0xFC, 0x78, 0x00, 0x00, 0x00, 0xC0, 0x80, 0xC0, 0x00, 0x00, 0x03, 0x1F, 0xFF, 0xF8, 0x00, 0x37, 0xF8, 0x38, 0x00, 0x00, 0x01, 0xC7, 0xF0, 0xE0, 0x00, 0x00, 0x03, 0x9F, 0xFE, 0x00,
	0x00, 0x5F, 0xFC, 0x38, 0x00, 0x00, 0x01, 0xC3, 0xE8, 0xE0, 0x00, 0x00, 0x03, 0x1F, 0xFB, 0x00, 0x00, 0x01, 0xFC, 0x3C, 0x00, 0x00, 0x01, 0xC7, 0xF8, 0xE0, 0x00, 0x00, 0x07, 0x1F, 0xC0, 0x00,
	0x00, 0x01, 0xFC, 0x3C, 0x00, 0x00, 0x03, 0x9F, 0xFC, 0x70, 0x00, 0x00, 0x07, 0x1F, 0x80, 0x00, 0x00, 0x01, 0xF8, 0x1E, 0x00, 0x00, 0x03, 0xBF, 0xFE, 0x78, 0x00, 0x00, 0x1E, 0x1F, 0xC0, 0x00,
	0x00, 0x01, 0xFC, 0x1E, 0x00, 0x00, 0x07, 0x9F, 0xFF, 0x78, 0x00, 0x00, 0x16, 0x1F, 0x80, 0x00, 0x00, 0x01, 0xFC, 0x1E, 0x00, 0x00, 0x07, 0x73, 0xC3, 0x3C, 0x00, 0x00, 0x1E, 0x1F, 0xC0, 0x00,
	0x00, 0x01, 0xFC, 0x1F, 0x00, 0x00, 0x1E, 0x60, 0x01, 0x9E, 0x00, 0x00, 0x3C, 0x1F, 0x80, 0x00, 0x00, 0x01, 0xF8, 0x1F, 0x80, 0x00, 0x7E, 0x40, 0x01, 0x17, 0x80, 0x00, 0x7E, 0x1F, 0xC0, 0x00,
	0x00, 0x01, 0xFC, 0x1F, 0x80, 0x00, 0x3C, 0x40, 0x01, 0x9F, 0x00, 0x00, 0x7C, 0x1F, 0x80, 0x00, 0x00, 0x01, 0xFC, 0x1F, 0xC0, 0x00, 0xFC, 0x40, 0x01, 0x07, 0xC0, 0x00, 0xFC, 0x1F, 0xC0, 0x00,
	0x00, 0x01, 0xFC, 0x1F, 0xF8, 0x0B, 0xE8, 0x7B, 0xFF, 0x81, 0xF8, 0x03, 0xFC, 0x1F, 0x80, 0x00, 0x00, 0x03, 0xF8, 0x1C, 0xFF, 0xFF, 0xC0, 0x3F, 0xFE, 0x00, 0xFF, 0xFF, 0xDE, 0x1F, 0xC0, 0x00,
	0x00, 0x05, 0xFC, 0x1C, 0xFF, 0xFF, 0x80, 0x7F, 0xDE, 0x00, 0xFF, 0x7F, 0xDC, 0x1F, 0x80, 0x00, 0x00, 0xFF, 0xFC, 0x1C, 0x3F, 0xFE, 0x00, 0x0E, 0xB8, 0x00, 0x3F, 0xFF, 0x1C, 0x1F, 0xFC, 0x00,
	0x2F, 0xFF, 0xF8, 0x1C, 0x00, 0x00, 0x00, 0x04, 0x98, 0x00, 0x01, 0x00, 0x1E, 0x1F, 0xFF, 0xF4, 0x3F, 0xFF, 0xFC, 0x3C, 0x00, 0x00, 0x00, 0x06, 0x98, 0x00, 0x00, 0x00, 0x1C, 0x1F, 0xFF, 0xFE,
	0x3F, 0xFF, 0xFC, 0x1C, 0x00, 0x00, 0x00, 0x06, 0x98, 0x00, 0x00, 0x00, 0x1E, 0x1F, 0xFF, 0xFE, 0x7F, 0xFF, 0xF8, 0x3C, 0x00, 0x00, 0x00, 0x02, 0xB8, 0x00, 0x00, 0x00, 0x1C, 0x1F, 0xFF, 0xFE,
	0x3F, 0xFF, 0xFC, 0x1C, 0x00, 0x00, 0x00, 0x03, 0xE8, 0x00, 0x00, 0x00, 0x1C, 0x1F, 0xFF, 0xFE, 0x7F, 0xFF, 0xFC, 0x3C, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1E, 0x1F, 0xFF, 0xFE,
	0x3F, 0xFF, 0xFC, 0x3C, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x1C, 0x1F, 0xFF, 0xFE, 0x7F, 0xFF, 0xF8, 0x3C, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1E, 0x1F, 0xFF, 0xFE,
	0x3F, 0xFF, 0xFC, 0x3C, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1C, 0x1F, 0xFF, 0xFE, 0x7F, 0xFF, 0xFC, 0x3C, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1E, 0x1F, 0xFF, 0xFE,
	0x3F, 0xFF, 0xFC, 0x3C, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1C, 0x1F, 0xFF, 0xFE, 0x7F, 0xFF, 0xF8, 0x3C, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1C, 0x1F, 0xFF, 0xFE,
	0x7F, 0xFF, 0xFC, 0x3C, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1E, 0x1F, 0xFF, 0xFE, 0x3F, 0xFF, 0xFC, 0x3C, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1C, 0x1F, 0xFF, 0xFE,
	0x7F, 0xFF, 0xFC, 0x38, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x16, 0x1F, 0xFF, 0xFE, 0x0B, 0xFF, 0xF8, 0x3C, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1C, 0x1F, 0xFF, 0xF8,
	0x00, 0x7F, 0xFC, 0x3C, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1E, 0x1F, 0xFE, 0x00, 0x00, 0x01, 0xFC, 0x3C, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1C, 0x1F, 0xC0, 0x00,
	0x00, 0x03, 0xFC, 0x3C, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x16, 0x1F, 0x80, 0x00, 0x00, 0x01, 0xF8, 0x3C, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1C, 0x1F, 0xC0, 0x00,
	0x00, 0x01, 0xFC, 0x3C, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1E, 0x1F, 0x80, 0x00, 0x00, 0x01, 0xFC, 0x38, 0x15, 0xB7, 0x80, 0x00, 0x00, 0x3E, 0x00, 0x00, 0x1C, 0x1F, 0xC0, 0x00,
	0x00, 0x01, 0xF8, 0x3C, 0x16, 0xAB, 0x00, 0x00, 0x00, 0x52, 0x00, 0x00, 0x16, 0x1F, 0x80, 0x00, 0x00, 0x01, 0xFC, 0x3C, 0x3F, 0xFB, 0x80, 0x00, 0x00, 0xFF, 0x80, 0x00, 0x1C, 0x1F, 0xC0, 0x00,
	0x00, 0x01, 0xFC, 0x3C, 0x1F, 0xE9, 0x00, 0x00, 0x01, 0xF7, 0x80, 0x00, 0x1E, 0x1F, 0x80, 0x00, 0x00, 0x01, 0xF8, 0x3C, 0x03, 0x82, 0x1D, 0xC6, 0x39, 0xC0, 0x07, 0x00, 0x14, 0x1F, 0xC0, 0x00,
	0x00, 0x01, 0xFC, 0x3C, 0x03, 0x87, 0x14, 0x85, 0x15, 0xC1, 0x03, 0x80, 0x1E, 0x1F, 0x80, 0x00, 0x00, 0x01, 0xFC, 0x3C, 0x03, 0x87, 0x9F, 0xE6, 0x3D, 0xC0, 0x1F, 0xC0, 0x1C, 0x1F, 0xC0, 0x00,
	0x00, 0xBF, 0xF8, 0x3C, 0x03, 0x83, 0x9F, 0xE7, 0x3F, 0x8D, 0x1E, 0xC0, 0x16, 0x1F, 0xD4, 0x00, 0x17, 0xFF, 0xFC, 0x3C, 0x03, 0x83, 0x9E, 0xE7, 0x79, 0xD7, 0xBC, 0xE0, 0x1C, 0x1F, 0xFF, 0xE8,
	0x1F, 0xFF, 0xFC, 0x3C, 0x03, 0x83, 0x9C, 0xE3, 0x3F, 0x9F, 0xBC, 0xE0, 0x1E, 0x1F, 0xFF, 0xD0, 0x3F, 0xFF, 0xF8, 0x3C, 0x03, 0x83, 0x9E, 0xE7, 0x79, 0xC7, 0xBC, 0xE0, 0x16, 0x1F, 0xFF, 0xFE,
	0x3F, 0xFF, 0xFC, 0x3C, 0x03, 0x83, 0xBC, 0xE3, 0xF9, 0xC3, 0xBC, 0xE0, 0x1C, 0x1F, 0xFF, 0xFC, 0x3F, 0xFF, 0xFC, 0x3C, 0x03, 0x83, 0x9E, 0xEB, 0xE9, 0xE3, 0xBD, 0xE0, 0x1E, 0x1F, 0xFF, 0xFE,
	0x3F, 0xFF, 0xF8, 0x3C, 0x03, 0x83, 0x9C, 0xE3, 0xF1, 0xC3, 0x9C, 0xC0, 0x1C, 0x1F, 0xFF, 0xFC, 0x3F, 0xFF, 0xFC, 0x3C, 0x03, 0x83, 0x9E, 0xE9, 0xE0, 0xFF, 0x9F, 0xC0, 0x16, 0x1F, 0xFF, 0xFE,
	0x3F, 0xFF, 0xFC, 0x3C, 0x03, 0x83, 0xBC, 0x61, 0xE0, 0xFF, 0x07, 0xC0, 0x1E, 0x1F, 0xFF, 0xFE, 0x3F, 0xFF, 0xFC, 0x3C, 0x00, 0x00, 0x00, 0x00, 0xC0, 0x00, 0x00, 0x00, 0x1C, 0x3F, 0xFF, 0xFC,
	0x3F, 0xFF, 0xF8, 0x3C, 0x00, 0x00, 0x00, 0x00, 0xE0, 0x00, 0x00, 0x00, 0x1E, 0x1F, 0xFF, 0xFE, 0x3F, 0xFF, 0xFC, 0x3C, 0x00, 0x00, 0x00, 0x07, 0xC0, 0x00, 0x00, 0x00, 0x1C, 0x1F, 0xFF, 0xFC,
	0x3F, 0xFF, 0xFC, 0x3C, 0x00, 0x00, 0x00, 0x17, 0xC0, 0x00, 0x00, 0x00, 0x16, 0x3F, 0xFF, 0xFE, 0x3F, 0xFF, 0xFC, 0x3C, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x1E, 0x3F, 0xFF, 0xFE,
	0x3F, 0xFF, 0xFC, 0x3C, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x1E, 0x3F, 0xFF, 0xFC, 0x3F, 0xFF, 0xFE, 0x3C, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1C, 0x7F, 0xFF, 0xFE,
	0x2F, 0xFF, 0xFF, 0xBC, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x88, 0x80, 0x5F, 0xFF, 0xFF, 0xF8, 0x00, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00,
	0x01, 0x3F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xA0, 0x00, 0x01, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF0, 0x00,
	0x00, 0x01, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x80, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x80, 0x00,
	0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x80, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}

func TestImageFromBytesMonochrome(t *testing.T) {
	image := pixel.NewImageFromBytes[pixel.Monochrome](128, 128, rprofile)
	if width, height := image.Size(); width != 128 && height != 128 {
		t.Errorf("image.Size(): expected 128, 128 but got %d, %d", width, height)
	}

	raw := image.RawBuffer()
	for i, b := range raw {
		if b != rprofile[i] {
			t.Fatalf("failed to roundtrip image. expected %v but got %v", rprofile[i], b)
		}
	}
}

// Test pixel formats by filling them with noise and checking whether they
// contain the same data afterwards.
func TestImageNoise(t *testing.T) {
	t.Run("RGB888", func(t *testing.T) {
		testImageNoiseN[pixel.RGB888](t)
	})
	t.Run("RGB565BE", func(t *testing.T) {
		testImageNoiseN[pixel.RGB565BE](t)
	})
	t.Run("RGB555", func(t *testing.T) {
		testImageNoiseN[pixel.RGB555](t)
	})
	t.Run("RGB444BE", func(t *testing.T) {
		testImageNoiseN[pixel.RGB444BE](t)
	})
	t.Run("Monochrome", func(t *testing.T) {
		testImageNoiseN[pixel.Monochrome](t)
	})
}

// Run the testImageNoise multiple times, because a single test might not catch
// all bugs (since the test uses random data).
func testImageNoiseN[T pixel.Color](t *testing.T) {
	for i := 0; i < 10; i++ {
		testImageNoise[T](t)
	}
}

func testImageNoise[T pixel.Color](t *testing.T) {
	// Create an image of a random width/height for extra testing.
	width := rand.Int()%500 + 10
	height := rand.Int()%500 + 10
	t.Log("image size:", width, height)

	// Create two images: the to-be-tested image object and a reference image.
	img := pixel.NewImage[T](width, height)
	ref := goimage.NewRGBA(goimage.Rect(0, 0, width, height))

	// Fill the two images with noise.
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Set a random color in both images.
			c := pixel.NewColor[T](uint8(rand.Uint32()), uint8(rand.Uint32()), uint8(rand.Uint32()))
			img.Set(x, y, c)
			ref.Set(x, y, c.RGBA())
		}
	}

	// Compare the two images. They should match.
	mismatch := 0
	firstX := 0
	firstY := 0
	var firstExpected, firstActual color.RGBA
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.Get(x, y).RGBA()
			r2, g2, b2, _ := ref.At(x, y).RGBA()
			c2 := color.RGBA{R: uint8(r2 >> 8), G: uint8(g2 >> 8), B: uint8(b2 >> 8), A: 255}
			if c != c2 {
				mismatch++
				if mismatch == 1 {
					firstX = x
					firstY = y
					firstExpected = c
					firstActual = c2
				}
			}
		}
	}
	if mismatch != 0 {
		t.Errorf("mismatch found: %d pixels are different (first diff at (%d, %d), expected %v, actual %v)", mismatch, firstX, firstY, firstExpected, firstActual)
	}
}
