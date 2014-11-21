package main

import (
	"fmt"
	"image"
	"image/color"
	"runtime"

	"github.com/banthar/Go-SDL/sdl"
	"github.com/banthar/Go-SDL/ttf"
	"github.com/dustin/go-humanize"
	"github.com/go-gl/gl"
)

type SdlHandler struct {
	Width      int
	Height     int
	Buddhabrot *Buddha
}

type TexCoords struct {
	TX, TY, TX2, TY2 float32
}

func NewSdlHandler(w, h int, b *Buddha) *SdlHandler {
	s := &SdlHandler{
		Width:      w,
		Height:     h,
		Buddhabrot: b,
	}
	return s
}

func (s *SdlHandler) Start() {
	runtime.LockOSThread()

	// Init SDL
	sdl.Init(sdl.INIT_VIDEO)
	defer sdl.Quit()

	// Prevent tearing
	sdl.GL_SetAttribute(sdl.GL_SWAP_CONTROL, 1)

	screen := sdl.SetVideoMode(s.Width, s.Height, 32, sdl.OPENGL)
	if screen == nil {
		panic(sdl.GetError())
	}

	if gl.Init() != 0 {
		panic("Could not init OpenGL")
	}

	if ttf.Init() != 0 {
		panic("Could not init TTF")
	}
	defer ttf.Quit()

	// Set window title
	sdl.WM_SetCaption("Go Buddhabrot", "Go Buddhabrot")

	// Set up the main view
	gl.Enable(gl.TEXTURE_2D)
	gl.Viewport(0, 0, s.Width, s.Height)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(s.Width), float64(s.Height), 0, -1, 1)

	// Clear screen
	gl.ClearColor(0, 0, 0, 0)

	// Set up fonts
	font := ttf.OpenFont("SourceSansPro-Black.otf", 12)
	if font == nil {
		panic("Could not open font file")
	}
	defer font.Close()

	var (
		text      *sdl.Surface
		textColor sdl.Color = sdl.Color{255, 255, 255, 0}
		running   bool      = true
		cpus      int       = runtime.NumCPU()
		count     int       = 0
	)

	for i := 0; i < cpus; i++ {
		go s.Buddhabrot.StartWorker()
	}

	for running {
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch e.(type) {
			case *sdl.QuitEvent:
				fmt.Printf("Total Samples: %s\n", humanize.Comma(s.Buddhabrot.Samples))
				s.Buddhabrot.saveImage()
				running = false
			}
		}

		if count%1000 == 0 {
			// Erase screen
			gl.Clear(gl.COLOR_BUFFER_BIT)

			// Render image
			s.RenderImageOntoScreen(s.Buddhabrot.Img)

			// Render sample counter
			// DOES NOT WORK AT THE MOMENT
			text = ttf.RenderUTF8_Blended(font, humanize.Comma(s.Buddhabrot.Samples), textColor)
			if text == nil {
				panic("Could not render font")
			}
			s.RenderSampleCounter(text)

			// Render
			sdl.GL_SwapBuffers()

			// Prevent overflow
			count = 0
		}

		count++
	}
}

func (s *SdlHandler) RenderImageOntoScreen(img *image.RGBA) {
	var tex gl.Texture = s.getTexture(s.Width, s.Height, []byte(img.Pix))
	var tc TexCoords = TexCoords{0, 0, 1, 1}

	// Draw image as texture
	tex.Bind(gl.TEXTURE_2D)
	drawQuad(0, 0, s.Width, s.Height, tc.TX, tc.TY, tc.TX2, tc.TY2)
	tex.Unbind(gl.TEXTURE_2D)
}

func (s *SdlHandler) RenderSampleCounter(text *sdl.Surface) {
	rect := image.Rect(0, 0, int(text.W), int(text.H))
	data := image.NewRGBA(rect)

	for j := 0; j < int(text.H); j++ {
		for i := 0; i < int(text.W); i++ {
			r, g, b, a := text.At(i, j).RGBA()
			c := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			data.SetRGBA(i, j, c)
		}
	}

	var tex gl.Texture = s.getTexture(int(text.W), int(text.H), data.Pix)
	var tc TexCoords = TexCoords{0, 0, 1, 1}

	// Draw image as texture
	tex.Bind(gl.TEXTURE_2D)
	drawQuad(0, 0, int(text.W), int(text.H), tc.TX, tc.TY, tc.TX2, tc.TY2)
	tex.Unbind(gl.TEXTURE_2D)
}

func (s *SdlHandler) getTexture(w, h int, data []byte) gl.Texture {
	id := gl.GenTexture()
	id.Bind(gl.TEXTURE_2D)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, w, h, 0, gl.RGBA, gl.UNSIGNED_BYTE, data)

	if gl.GetError() != gl.NO_ERROR {
		id.Delete()
		panic("Failed to load a texture")
		return 0
	}

	return id
}

func drawQuad(x, y, w, h int, u, v, u2, v2 float32) {
	gl.Begin(gl.QUADS)

	gl.TexCoord2f(float32(u), float32(v))
	gl.Vertex2i(int(x), int(y))

	gl.TexCoord2f(float32(u2), float32(v))
	gl.Vertex2i(int(x+w), int(y))

	gl.TexCoord2f(float32(u2), float32(v2))
	gl.Vertex2i(int(x+w), int(y+h))

	gl.TexCoord2f(float32(u), float32(v2))
	gl.Vertex2i(int(x), int(y+h))

	gl.End()
}
