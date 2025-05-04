package main

import (
	"image/color"
	"log"

	"chip8-emulator/beeper"
	"chip8-emulator/chip8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func NewGame(romPath string) (*Game, error) {
	cpu := chip8.InitCpu()
	if loadErr := cpu.LoadProgram(romPath); loadErr != nil {
		log.Println("FAILED TO LOAD ROM")
		return nil, loadErr
	}
	log.Println("Game loaded successfully")

	b, err := beeper.Init()
	if err != nil {
		log.Println("Error initializing beeper")
	}
	cpu.AddBeep(b.Play)
	log.Println("Beeper initialized successfully")

	game := &Game{
		cpu:             cpu,
		backgroundColor: color.RGBA{0, 0, 0, 255},
		foregroundColor: color.RGBA{255, 255, 255, 255},
		scale:           DEFAULT_SCALE,
	}

	return game, nil
}

func (g *Game) handleInput() {
	keyMap := map[ebiten.Key]uint{
		ebiten.Key1: 0x1, ebiten.Key2: 0x2, ebiten.Key3: 0x3, ebiten.Key4: 0xC,
		ebiten.KeyQ: 0x4, ebiten.KeyW: 0x5, ebiten.KeyE: 0x6, ebiten.KeyR: 0xD,
		ebiten.KeyA: 0x7, ebiten.KeyS: 0x8, ebiten.KeyD: 0x9, ebiten.KeyF: 0xE,
		ebiten.KeyZ: 0xA, ebiten.KeyX: 0x0, ebiten.KeyC: 0xB, ebiten.KeyV: 0xF,
	}

	for key, value := range keyMap {
		if inpututil.IsKeyJustPressed(key) {
			g.cpu.Key(uint8(value), true)
		} else if inpututil.IsKeyJustReleased(key) {
			g.cpu.Key(uint8(value), false)
		}
	}
}
