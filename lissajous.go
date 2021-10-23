package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"time"
)

type Buffer struct {
	content *[]byte
}

func (b Buffer) Read(p []byte) (n int, err error) {
	n = copy(p, *b.content)
	*b.content = (*b.content)[n:]
	if n == 0 {
		return 0, io.EOF
	}
	return n, nil
}

func (b Buffer) Write(p []byte) (n int, err error) {
	*b.content = append(*b.content, p...)
	return len(p), nil
}

func init() {
	rand.Seed(time.Now().Unix())
}

func RandomColor() color.Color {
	n := func() uint8 {
		return uint8(rand.Intn(255))
	}
	return color.RGBA{n(), n(), n(), 1}
}

func Lissajous(cycles float64, palette []color.Color) io.Reader {
	index := func() uint8 {
		x := rand.Intn(len(palette)-1) + 1
		return uint8(x)
	}
	const (
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), index())
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	out := Buffer{content: &[]byte{}}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
	return out
}
