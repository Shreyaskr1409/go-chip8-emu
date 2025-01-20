package chip8

import (
	"fmt"
)

// 1. create the chip 8 struct
type cpu struct {
	opcode uint16
	memory [4098]uint8
	v      [16]uint8
	i      uint16
	pc     uint16
	stack  [16]uint16
	stkp   uint16

	key      [16]uint8
	display  [32][16]uint8
	drawFlag bool

	delayTimer uint8
	soundTimer uint8
}

// 2. create a function to initialize the chip8 struct (initialize registers)
func (c *cpu) initialize() {
	// TODO
	c.pc = 0x200 // chip-8 program counter starts at 0x200
	c.opcode = 0
	c.i = 0
	c.stkp = 0

	// clear display
	// clear stack
	// clear c registers
	// clear memory

	// load fontset
	for i := 0; i < 80; i += 1 {
		// c.memory[i] = chip8_fontset[i] // i will prepare this fontset later
	}

	// reset timers
}

// 3. emulate one cycle
// fetch opcode
// decode opcode
// execute opcode
// update timers
func (c *cpu) emulateCycle() {
	// in chip-8 addresses are stored in an array containing 1 byte as seem in c.memory
	// as an opcode is 2 bytes long, i will fetch 2 bytes from the array
	c.opcode = uint16(c.memory[c.pc]<<8) | uint16(c.memory[c.pc+1]) // ignore the warning getting shown
	// remember to increment the program counter by 2, not 1

	c.executeOp()
}

func (c *cpu) executeOp() {
	// TODO
	switch c.opcode & 0xF000 {
	case 0x0000:
		// only 0x0NNN, 0x00E0, 0x00EE

		// wikipedia says 0x0NNN is not necessary for most ROMs

		switch c.opcode & 0x000F {
		// only 0x00E0, 0x00EE

		case 0x0000: // clear screen
			for i := 0; i < len(c.display); i++ {
				for j := 0; j < len(c.display[0]); j++ {
					c.display[i][j] = 0x0
				}
			}
			c.drawFlag = true
			c.pc += 2
		case 0x000E: // return from a subroutine
			// we stored the addresses we have been on so far in the stack
			// what we are doing here is returning from an address
			// to the previous address on the stack
			c.stkp -= 1 // remove the top element from the stack pointer
			// since we are returning from it
			c.pc = c.stack[c.stkp]
			c.pc += 2
		}

	case 0x1000: // only 1NNN, jumps to address 0x0NNN
		c.pc = c.opcode & 0x0FFF
	case 0x2000: // only 2NNN, jump to subroutine 0x0NNN
		// before jumping, we will save the address to the stack
		c.stack[c.stkp] = c.pc
		c.stkp += 1
		c.pc = c.opcode & 0x0FFF
	case 0x3000: // only 3XNN, skips next instruction if V[x] == NN
		if uint16(c.v[c.opcode&0xF000]>>8) == c.opcode&0x00FF {
			c.pc = c.pc + 4 // 2 bytes for execution and and 2 for skipping next instruction
		} else {
			c.pc = c.pc + 2
		}
	case 0x4000: // TODO

	default:
		fmt.Println("Invalid opcode ", c.opcode)
	}
}

func updateTimers() {
	// TODO
}
