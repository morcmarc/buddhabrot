package main

import (
	"fmt"
	"image"
	"runtime"

	"github.com/banthar/Go-SDL/sdl"
	// "github.com/banthar/Go-SDL/ttf"
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

	if sdl.SetVideoMode(s.Width, s.Height, 32, sdl.OPENGL) == nil {
		panic("Could not start SDL")
	}

	if gl.Init() != 0 {
		panic("Could not init OpenGL")
	}

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
	// font := ttf.OpenFont("scpr.ttf", 12)
	// if font == nil {
	// 	panic("Could not open font file")
	// }
	// defer font.Close()

	var (
		// text       *sdl.Surface
		// textColor  sdl.Color = sdl.Color{255, 255, 255, 255}
		running    bool = true
		generating bool = false
		cpus       int  = runtime.NumCPU()
		count      int  = 0
	)

	for running {
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch e.(type) {
			case *sdl.QuitEvent:
				fmt.Printf("Total Samples: %s\n", humanize.Comma(s.Buddhabrot.Samples))
				s.Buddhabrot.saveImage()
				running = false
			}
		}

		if !generating {
			for i := 0; i < cpus; i++ {
				go s.Buddhabrot.StartWorker()
			}
			generating = true
		}

		count++
		if count%10000 == 0 {
			s.RenderImageOntoScreen(s.Buddhabrot.Img)
			// text = ttf.RenderUTF8_Blended(font, fmt.Sprintf("Total Samples: %s\n", humanize.Comma(s.Buddhabrot.Samples)), textColor)
			// if text == nil {
			// 	panic("Could not render font")
			// }
			count = 0
		}
	}
}

func (s *SdlHandler) RenderImageOntoScreen(img *image.RGBA) {
	var tex gl.Texture = s.getTexture(img)
	var tc TexCoords = TexCoords{0, 0, 1, 1}

	// Erase screen
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// Draw image as texture
	tex.Bind(gl.TEXTURE_2D)
	drawQuad(0, 0, s.Width, s.Height, tc.TX, tc.TY, tc.TX2, tc.TY2)
	tex.Unbind(gl.TEXTURE_2D)

	// Render
	sdl.GL_SwapBuffers()
}

func (s *SdlHandler) getTexture(img *image.RGBA) gl.Texture {
	id := gl.GenTexture()
	id.Bind(gl.TEXTURE_2D)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, s.Width, s.Height, 0, gl.RGBA, gl.UNSIGNED_BYTE, []byte(img.Pix))

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
