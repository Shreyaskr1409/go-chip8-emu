package main

import (
	"image/color"
	"log"

	"chip8-emulator/beeper"
	"chip8-emulator/chip8"

	"github.com/hajimehoshi/ebiten/v2"
)

func NewGame(romPath string, scale int) (*Game, error) {
	cpu := chip8.InitCpu()
	if loadErr := cpu.LoadProgram(romPath); loadErr != nil {
		log.Println("FAILED TO LOAD ROM")
		return nil, loadErr
	}

	b, err := beeper.Init()
	if err != nil {
		log.Println("Error initializing beeper")
	}
	cpu.AddBeep(b.Play)

	game := &Game{
		cpu:             cpu,
		backgroundColor: color.RGBA{0, 0, 0, 255},
		foregroundColor: color.RGBA{255, 255, 255, 255},
	}

	return game, nil
}

func (g *Game) handleInput() {
	keyMap := map[ebiten.Key]uint{
		ebiten.Key1: 0x1, ebiten.Key2: 0x2, ebiten.Key3: 0x3, ebiten.Key4: 0xC,
		ebiten.KeyQ: 0x4, ebiten.KeyW: 0x5, ebiten.KeyE: 0x6, ebiten.KeyR: 0xD,
		// Work in progress
	}
}
