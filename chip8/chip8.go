package chip8

import (
	"fmt"
	"log"
	"os"
)

type Cpu struct {
	opcode uint16
	memory [4098]uint8
	v      [16]uint8
	i      uint16
	pc     uint16
	stack  [16]uint16
	stkp   uint16

	key      [16]uint8
	display  [32][64]uint8
	drawFlag bool

	delayTimer uint8
	soundTimer uint8
	beeper     func()
}

var fontset = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

func InitCpu() *Cpu {
	c := Cpu{
		drawFlag: true,
		pc:       0x200, // chip-8 program counter starts at 0x200
		beeper:   func() {},
	}

	// load fontset
	for i := 0; i < 80; i += 1 {
		c.memory[i] = fontset[i] // i will prepare this fontset later
	}

	return &c
}

func (c *Cpu) EmulateCycle() {
	// in chip-8 addresses are stored in an array containing 1 byte in c.memory
	// as an opcode is 2 bytes long, i will fetch 2 bytes from the array
	c.opcode = uint16(c.memory[c.pc]<<8) | uint16(c.memory[c.pc+1])
	// ignore the warning getting shown

	// remember to increment the program counter by 2, not 1

	c.executeOp()
	log.Println("Opcode executing: ", c.opcode)
	c.updateTimers()
}

func (c *Cpu) updateTimers() {
	if c.delayTimer > 0 {
		c.delayTimer -= 1
	}
	if c.soundTimer > 0 {
		if c.soundTimer == 1 {
			c.beeper()
		}
		c.soundTimer = c.soundTimer - 1
	}
}

func (c *Cpu) Buffer() [32][64]uint8 {
	return c.display
}

func (c *Cpu) Draw() bool {
	drawFlag := c.drawFlag
	c.drawFlag = false
	return drawFlag
}

func (c *Cpu) AddBeep(fn func()) {
	c.beeper = fn
}

func (c *Cpu) Key(num uint8, down bool) {
	if down {
		c.key[num] = 1
	} else {
		c.key[num] = 0
	}
}

func (c *Cpu) LoadProgram(filename string) error {
	file, fileErr := os.OpenFile(filename, os.O_RDONLY, 0o777)
	if fileErr != nil {
		return fileErr
	}
	defer file.Close()

	fStat, fStatErr := file.Stat()
	if fStatErr != nil {
		return fStatErr
	}
	if int64(len(c.memory)-512) < fStat.Size() { // 512 = 0x200, where program is loaded
		return fmt.Errorf("program size bigger than memory %X", fStat.Size())
	}

	buffer := make([]byte, fStat.Size())
	if _, readErr := file.Read(buffer); readErr != nil {
		return readErr
	}

	// load buffer into memory
	for i := 0; i < len(buffer); i++ {
		c.memory[i+512] = buffer[i]
	}

	return nil
}
