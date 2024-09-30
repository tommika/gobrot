package mandelbrot

import (
	"fmt"
	"log"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type ColorRGBA struct {
	r, g, b, a uint8
}

type Palette []ColorRGBA

type Image struct {
	surface *sdl.Surface
	W, H    int32
}

var BLACK = ColorRGBA{0, 0, 0, 0}

func CreatePalette(count int, r, g, b float64) Palette {
	pal := make(Palette, count, count)
	for i := 0; i < count; i++ {
		v := 0.1 + 0.9*(float64(i)/float64(count))
		r8 := uint8(math.Min(255.0, v*r*255.0))
		g8 := uint8(math.Min(255.0, v*g*255.0))
		b8 := uint8(math.Min(255.0, v*b*255.0))
		a8 := uint8(0)
		pal[i] = ColorRGBA{r8, g8, b8, a8}
	}
	return pal
}

func CreateSurface(w, h int32) *sdl.Surface {
	surface, err := sdl.CreateRGBSurface(0, w, h, 32, 0, 0, 0, 0)
	if err != nil || surface == nil {
		log.Printf("Failed to allocate surface for image: %v\n", err)
		return nil
	}
	log.Printf("Created surface for image: w=%d, h=%d, %s\n", w, h, DescribeSurfaceFormat(surface))
	return surface
}

func SetSurfacePixel(surface *sdl.Surface, x, y int, rgba ColorRGBA) {
	i := int32(y)*surface.Pitch + int32(x)*int32(surface.Format.BytesPerPixel)
	pix := surface.Pixels()
	// Assumed PIXELFORMAT_RGB888
	pix[i] = rgba.b
	pix[i+1] = rgba.g
	pix[i+2] = rgba.r
}

func DescribeSurfaceFormat(surface *sdl.Surface) string {
	pf := surface.Format
	name := sdl.GetPixelFormatName(uint(pf.Format))
	return fmt.Sprintf("Format=%s(%d), BytesPerPixel=%d BitsPerPixel=%d\n",
		name, pf.Format, pf.BytesPerPixel, pf.BitsPerPixel)
}
