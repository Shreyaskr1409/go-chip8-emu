package main

import (
	"image/color"
	"os"

	"chip8-emulator/chip8"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	DEFAULT_WIDTH  int32 = 64
	DEFAULT_HEIGHT int32 = 32
)

type Game struct {
	cpu             *chip8.Cpu
	backgroundColor color.RGBA
	foregroundColor color.RGBA
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 320
}

func main() {
	if len(os.Args) < 2 {
		panic("Please provide c8 file")
	}
	modifier := 10
	filename := os.Args[2]

	c8 := chip8.InitCpu()
	if loadErr := c8.LoadProgram(filename); loadErr != nil {
		panic(loadErr)
	}

	ebiten.SetWindowSize(int(DEFAULT_WIDTH)*modifier, int(DEFAULT_HEIGHT)*modifier)
	ebiten.SetWindowTitle("CHIP-8 EMULATOR - " + filename)
	game := &Game{}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
