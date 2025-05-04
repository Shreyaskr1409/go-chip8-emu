package main

import (
	"image/color"
	"log"
	"os"

	"chip8-emulator/chip8"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	DEFAULT_WIDTH  int32 = 64
	DEFAULT_HEIGHT int32 = 32
	DEFAULT_SCALE  int   = 10
)

type Game struct {
	cpu             *chip8.Cpu
	backgroundColor color.RGBA
	foregroundColor color.RGBA
	scale           int
}

func (g *Game) Update() error {
	g.handleInput()
	g.cpu.EmulateCycle()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.backgroundColor)
	if g.cpu.Draw() {
		buffer := g.cpu.Buffer()

		for y := 0; y < int(DEFAULT_HEIGHT); y++ {
			for x := 0; x < int(DEFAULT_WIDTH); x++ {
				if buffer[y][x] != 0 {
					for dy := 0; dy < g.scale; dy++ {
						for dx := 0; dx < g.scale; dx++ {
							screen.Set(x*g.scale+dx, y*g.scale+dy, g.foregroundColor)
						}
					}
				}
			}
		}
	}
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth, screenHeight int) {
	return int(DEFAULT_WIDTH) * g.scale, int(DEFAULT_HEIGHT) * g.scale
}

func main() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer file.Close()
	log.SetOutput(file)
	log.Println("Application starts")

	if len(os.Args) < 2 {
		log.Panic("Please provide c8 file")
	}
	filename := os.Args[1]

	ebiten.SetWindowSize(int(DEFAULT_WIDTH)*DEFAULT_SCALE, int(DEFAULT_HEIGHT)*DEFAULT_SCALE)
	ebiten.SetWindowTitle("CHIP-8 EMULATOR - " + filename)

	game, err := NewGame(filename)
	if err != nil {
		log.Println("Error encountered while creating new game:\n", err)
	}
	log.Println("Game initialized")

	if err := ebiten.RunGame(game); err != nil {
		log.Println("Error:\n", err)
		panic(err)
	}
}
