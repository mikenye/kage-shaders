package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed gradient_border.kage
var shaderProgram []byte

//go:embed gopher.png
var pngGopher []byte

// Struct implementing the ebiten.Game interface.
type Game struct {
	borderEnable int
	t            int
	shader       *ebiten.Shader
	sprite       *ebiten.Image
}

// Assume a fixed layout.
func (g *Game) Layout(_, _ int) (int, int) {
	return 512, 512
}

// Update game state
func (g *Game) Update() error {

	// time variable
	g.t++
	if g.t >= 120 {
		g.t = 0
	}

	// toggle border on and off
	if g.t < 60 {
		g.borderEnable = 1
	} else {
		g.borderEnable = 0
	}
	return nil
}

// Core drawing function from where we call DrawRectShader.
func (g *Game) Draw(screen *ebiten.Image) {

	// create draw options
	opts := &ebiten.DrawRectShaderOptions{}
	opts.GeoM.Translate(0, 0) // you could adjust the drawing position here

	// prep uniforms to pass to shader
	opts.Uniforms = make(map[string]interface{})

	// set border colour
	opts.Uniforms["BorderColor"] = []float32{0.87, 0.68, 0.15, .5}

	// toggle border on and off
	opts.Uniforms["BorderEnable"] = g.borderEnable

	// send image to shader
	opts.Images[0] = g.sprite

	// draw shader
	screen.DrawRectShader(512, 512, g.shader, opts)
}

func main() {

	// load image
	img, _, err := image.Decode(bytes.NewReader(pngGopher))
	if err != nil {
		log.Fatal(err)
	}

	// create game struct
	game := &Game{}

	// compile the shader
	game.shader, err = ebiten.NewShader(shaderProgram)
	if err != nil {
		log.Fatal(err)
	}

	// draw the gopher image to the centre of the sprite
	game.sprite = ebiten.NewImage(512, 512)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(256-img.Bounds().Dx()/2), float64(256-img.Bounds().Dy()/2))
	game.sprite.DrawImage(ebiten.NewImageFromImage(img), op)

	// configure window and run game
	ebiten.SetWindowTitle("Kage Shader: Gradient Border")
	ebiten.SetWindowSize(512, 512)
	err = ebiten.RunGame(game)
	if err != nil {
		log.Fatal(err)
	}
}
