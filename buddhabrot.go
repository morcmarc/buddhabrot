package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/dustin/go-humanize"
)

const (
	REAL_L float64 = -2.0
	REAL_H float64 = 2.0
	IMAG_L float64 = -2.0
	IMAG_H float64 = 2.0
	// REAL_L float64 = -0.0453594
	// REAL_H float64 = -0.0473594
	// IMAG_L float64 = -0.985749
	// IMAG_H float64 = -0.987749
)

type Buddha struct {
	Width   int
	Height  int
	Palette []int
	Img     *image.RGBA
	Samples int64
}

type Trajectories struct {
	Real []float64
	Imag []float64
}

type DensityMatrix struct {
	Red   []int
	Green []int
	Blue  []int
}

func NewBuddha(w, h int, p []int) *Buddha {
	rect := image.Rect(0, 0, w, h)
	return &Buddha{
		Width:   w,
		Height:  h,
		Palette: p,
		Samples: 0,
		Img:     image.NewRGBA(rect),
	}
}

func NewTrajectories(s int) *Trajectories {
	return &Trajectories{
		Real: make([]float64, s),
		Imag: make([]float64, s),
	}
}

func NewDensityMatrix(w, h int) *DensityMatrix {
	s := w * h
	return &DensityMatrix{
		Red:   make([]int, s),
		Green: make([]int, s),
		Blue:  make([]int, s),
	}
}

func (b *Buddha) Run(out *os.File) {
	e := MaxInt(b.Palette)
	densityMtx := NewDensityMatrix(b.Width, b.Height)

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	fmt.Println("Press control-c to terminate execution and render image.")

	cpus := runtime.NumCPU()
	for i := 0; i < cpus; i++ {
		go b.mainLoop(densityMtx, e)
	}

	sig := <-signalChannel
	switch sig {
	default:
		fmt.Printf("\nTotal Samples: %s\n", humanize.Comma(b.Samples))
		b.writeImage(out, densityMtx)
	}

}

func (b *Buddha) mainLoop(densityMtx *DensityMatrix, e int) {
	for {
		Cr, Ci := GetRandomComplex(REAL_L, REAL_H, IMAG_L, IMAG_H)

		if isInMSet(Cr, Ci, 0, e) {
			continue
		}

		n, t := b.getTrajectories(Cr, Ci, e)
		b.evaluateTrajectories(n, t, densityMtx)

		b.Samples++
	}
}

func (b *Buddha) getTrajectories(Cr, Ci float64, e int) (int, *Trajectories) {
	var ZrPrev float64 = 0.0
	var ZiPrev float64 = 0.0
	var t *Trajectories = NewTrajectories(e)
	var n int = 0

	for n = 0; n < e; n++ {
		t.Real[n] = ZrPrev*ZrPrev - ZiPrev*ZiPrev + Cr
		t.Imag[n] = 2.0*ZrPrev*ZiPrev + Ci
		ZrPrev = t.Real[n]
		ZiPrev = t.Imag[n]
		if ZrPrev*ZrPrev+ZiPrev*ZiPrev > 4.0 {
			break
		}
	}

	return n, t
}

func (b *Buddha) evaluateTrajectories(n int, t *Trajectories, d *DensityMatrix) {
	for i := 0; i < n; i++ {
		// Translate R and I coordinates into pixel coords.
		tr := ((t.Real[i] - REAL_L) / math.Abs(REAL_H-REAL_L))
		ti := ((t.Imag[i] - IMAG_L) / math.Abs(IMAG_H-IMAG_L))
		coordX := int(math.Floor(tr * float64(b.Width)))
		coordY := int(math.Floor(ti * float64(b.Height)))

		if coordX < 0 || coordX >= b.Width || coordY < 0 || coordY >= b.Height {
			continue
		}

		if n < b.Palette[0] {
			d.Red[coordY*b.Width+coordX]++
		}

		if n < b.Palette[1] {
			d.Green[coordY*b.Width+coordX]++
		}

		if n < b.Palette[2] {
			d.Blue[coordY*b.Width+coordX]++
		}
	}
}

func (b *Buddha) writeImage(out *os.File, d *DensityMatrix) {
	fmt.Println("Writing image file")
	var maxR int = MaxInt(d.Red)
	var maxG int = MaxInt(d.Green)
	var maxB int = MaxInt(d.Blue)
	fmt.Printf("Hits: R=%d G=%d B=%d\n", maxR, maxG, maxB)
	var fR float64 = float64(256) / float64(maxR)
	var fG float64 = float64(256) / float64(maxG)
	var fB float64 = float64(256) / float64(maxB)

	for i := 0; i < b.Width; i++ {
		for j := 0; j < b.Height; j++ {
			cR := uint8(float64(d.Red[j*b.Width+i]) * fR)
			cG := uint8(float64(d.Green[j*b.Width+i]) * fG)
			cB := uint8(float64(d.Blue[j*b.Width+i]) * fB)
			c := color.RGBA{cR, cG, cB, 255}
			// Set color AND rotate image (j <=> i)
			b.Img.SetRGBA(j, i, c)
		}
	}

	err := png.Encode(out, b.Img)
	if err != nil {
		log.Fatalf("Could not write image: %s", err)
	}
}

func isInMSet(Cr, Ci float64, minIter, maxIter int) bool {
	// Main cardioid
	if !(((Cr-0.25)*(Cr-0.25)+(Ci*Ci))*(((Cr-0.25)*(Cr-0.25)+(Ci*Ci))+(Cr-0.25)) < 0.25*Ci*Ci) {
		// 2nd order period bulb
		if !((Cr+1.0)*(Cr+1.0)+(Ci*Ci) < 0.0625) {
			// smaller bulb left of the period-2 bulb
			if !((((Cr + 1.309) * (Cr + 1.309)) + Ci*Ci) < 0.00345) {
				// smaller bulb bottom of the main cardioid
				if !((((Cr + 0.125) * (Cr + 0.125)) + (Ci-0.744)*(Ci-0.744)) < 0.0088) {
					// smaller bulb top of the main cardioid
					if !((((Cr + 0.125) * (Cr + 0.125)) + (Ci+0.744)*(Ci+0.744)) < 0.0088) {
						return false
					}
				}
			}
		}
	}

	return true
}
