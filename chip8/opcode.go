package chip8

import (
	"log"
	"math/rand"
)

// No need to call directly from main loop, this will be called from emulateCycle()
func (c *Cpu) executeOp() {
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
		if uint16(c.v[c.opcode&0x0F00>>8]) == c.opcode&0x00FF {
			c.pc += 4 // 2 bytes for execution and and 2 for skipping next instruction
		} else {
			c.pc += 2
		}
	case 0x4000: // only 4XNN, skips next instruction if V[x] != NN
		if uint16(c.v[c.opcode&0x0F00>>8]) != c.opcode&0x00FF {
			c.pc += 4
		} else {
			c.pc += 2
		}
	case 0x5000: // only 5XY0
		if c.v[c.opcode&0x0F00>>8] == c.v[c.opcode&0x00F0>>4] {
			c.pc += 4
		} else {
			c.pc += 2
		}
	case 0x6000: // only 6XNN, Sets V[x] to NN
		c.v[c.opcode&0x0F00>>8] = uint8(c.opcode & 0x00FF)
		c.pc += 2
	case 0x7000: // only 7XNN, adds NN to V[x] without changing carry flag
		c.v[c.opcode&0x0F00>>8] += uint8(c.opcode & 0x00FF)
		c.pc += 2
	case 0x8000: // 8XY0/1/2/3/4/5/6/E
		switch c.opcode & 0x000F {
		case 0x0000: // 8XY0, sets V[x] to value of V[y]
			c.v[c.opcode&0x0F00>>8] = c.v[c.opcode&0x00F0>>4]
			c.pc += 2
		case 0x0001: // 8XY1, sets V[x] to V[x]|V[y]
			c.v[c.opcode&0x0F00>>8] |= c.v[c.opcode&0x00F0>>4]
			c.pc += 2
		case 0x0002: // 8XY2, sets V[X] to V[x]&V[y]
			c.v[c.opcode&0x0F00>>8] &= c.v[c.opcode&0x00F0>>4]
			c.pc += 2
		case 0x0003: // 8XY3, sets V[x] to V[x]^V[y]
			c.v[c.opcode&0x0F00>>8] ^= c.v[c.opcode&0x00F0>>4]
			c.pc += 2
		case 0x0004: // 8XY4, sets V[x] to V[x]+V[y] with overflow flag V[0xF] set if required
			if c.v[c.opcode&0x0F00>>8] > 0xFF-c.v[c.opcode&0x00F0>>4] {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 0
			}
			c.v[c.opcode&0x0F00>>8] += c.v[c.opcode&0x00F0>>4]
			c.pc += 2
		case 0x0005: // 8XY5, sets V[x] to V[x]+V[y] with underflow flag V[0xF] unset if required
			if c.v[c.opcode&0x00F0>>4] > c.v[c.opcode&0x0F00>>8] {
				c.v[0xF] = 0
			} else {
				c.v[0xF] = 1
			}
			c.v[c.opcode&0x0F00>>8] = c.v[c.opcode&0x0F00>>8] - c.v[c.opcode&0x00F0>>4]
			c.pc += 2
		case 0x0006: // 8XY6, sets V[0xF] to l.s.b. of X and shifts X to right by 1
			c.v[0xF] = c.v[c.opcode&0x0F00>>8] & 0x1
			c.v[c.opcode&0x0F00>>8] = c.v[c.opcode&0x0F00>>8] >> 1
			c.pc += 2
		case 0x0007: // 8XY7, sets V[x] to V[y] - V[x] where V[0xF] is 0 when underflow
			if c.v[c.opcode&0x0F00>>8] > c.v[c.opcode&0x00F0>>4] {
				c.v[0xF] = 0
			} else {
				c.v[0xF] = 1
			}
			c.v[c.opcode&0x0F00>>8] = c.v[c.opcode&0x00F0>>4] - c.v[c.opcode&0x0F00>>8]
			c.pc += 2
		case 0x000E: // 8XYE, shifts V[x] <<= 1 then V[0xF] is 1 if MSB earlier was 1 else 0
			c.v[0xF] = c.v[c.opcode&0x0F00>>8] >> 7
			c.v[c.opcode&0x0F00>>8] = c.v[c.opcode&0x0F00>>8] << 1
			c.pc += 2
		default:
			log.Printf("Invalid opcode %X\n", c.opcode)
		}
	case 0x9000: // 9XY0, skips next instruction if V[x] != V[y]
		if c.v[c.opcode&0x0F00>>8] != c.v[c.opcode&0x00F0>>4] {
			c.pc += 4
		} else {
			c.pc += 2
		}
	case 0xA000: // ANNN, sets I to the address of NNN
		c.i = c.opcode & 0x0FFF
		c.pc += 2
	case 0xB000: // BNNN, jumps to address NNN + v[0x0]
		c.pc = c.opcode&0x0FFF + uint16(c.v[0x0])
	case 0xC000: // CXNN, V[x] = NN & random no b'w 0 and 255
		c.v[c.opcode&0x0F00>>8] = uint8(rand.Intn(256)) & uint8(c.opcode&0x00FF)
		c.pc += 2
	case 0xD000: // DXYN, draws a sprite at coordinate (V[x], V[y]) of height N

		// ---------WIKIPEDIA EXPLANATION----------
		// Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels
		// and a height of N pixels. Each row of 8 pixels is read as bit-coded
		// starting from memory location I; I value does not change after the
		// execution of this instruction. As described above, VF is set to 1 if
		// any screen pixels are flipped from set to unset when the sprite is
		// drawn, and to 0 if that does not happen.

		x := c.v[c.opcode&0x0F00>>8]
		y := c.v[c.opcode&0x00F0>>4]
		h := c.opcode & 0x000F
		w := uint16(8)

		c.v[0xF] = 0

		for j := uint16(0); j < h; j++ {
			pixel := c.memory[c.i+j]

			for i := uint16(0); i < w; i++ {
				if (pixel & (0x80 >> i)) != 0 {
					if c.display[y+uint8(j)][x+uint8(i)] == 1 {
						c.v[0xF] = 1 // Set flag
					}
					c.display[y+uint8(j)][x+uint8(i)] ^= 1 // draw
				}
			}
		}

		c.drawFlag = true
		c.pc += 2
	case 0xE000:
		switch c.opcode & 0x00FF {
		case 0x009E: // EX9E, skips next instruction if key stored in V[x] is pressed by player
			if c.key[c.v[c.opcode&0x0F00>>8]] == 1 {
				c.pc += 4
			} else {
				c.pc += 2
			}
		case 0x00A1: // EXA1, skips next instruction if key stored in V[x] is not pressed by player
			if c.key[c.v[c.opcode&0x0F00>>8]] == 0 {
				c.pc += 4
			} else {
				c.pc += 2
			}
		default:
			log.Printf("Invalid opcode %X\n", c.opcode)
		}
	case 0xF000:
		switch c.opcode & 0x00FF {
		case 0x0007: // FX07, V[x] = delayTimer
			c.v[c.opcode&0x0F00>>8] = c.delayTimer
			c.pc += 2
		case 0x000A: // FX0A, wait for key press (hault everything except
			// delay and sound should still continue to work) then V[x] = key
			pressed := false
			for i := 0; i < len(c.key); i++ {
				if c.key[i] != 0 {
					c.v[c.opcode&0x0F00>>8] = uint8(i)
					pressed = true
				}
			}
			if !pressed {
				return
			}
			c.pc += 2
		case 0x0015: // FX15, delayTimer = V[x]
			c.delayTimer = c.v[c.opcode&0x0F00>>8]
			c.pc += 2
		case 0x0018: // FX18, soundTimer = V[x]
			c.soundTimer = c.v[c.opcode&0x0F00>>8]
			c.pc += 2
		case 0x001E: // FX1E, adds V[x] to I, V[0xF] is not affected
			c.i += uint16(c.v[c.opcode&0x0F00>>8])
			c.pc += 2
		case 0x0029: // FX29, refer to the wikipedia explanation

			// ---------WIKIPEDIA EXPLANATION----------
			// Sets I to the location of the sprite for the character in
			// VX(only consider the lowest nibble). Characters 0-F (in
			// hexadecimal) are represented by a 4x5 font.[24]

			lowest_nibble := c.v[c.opcode&0x0F00>>8] & 0x0F
			c.i = uint16(lowest_nibble * 5)

			// In CHIP-8, font data is typically stored at the beginning of memory
			// Each character sprite is 5 bytes long (4x5 pixels)

			c.pc += 2
		case 0x0033: // FX33, stores BCD value of V[x]
			// MSB in I, then I+1, then LSP in I+2
			log.Println("Executing: ", c.opcode)
			c.memory[c.i] = c.v[c.opcode&0x0F00>>8] / 100
			c.memory[c.i+1] = c.v[c.opcode&0x0F00>>8] % 10
			c.memory[c.i+2] = c.v[c.opcode&0x0F00>>8] / 10
			c.pc += 2
		case 0x0055: // FX55, stores V[0] t0 V[x] (inclusive) with values from memory,
			// starting at address I till I + x
			for i := 0; i < int(c.opcode&0x0F00>>8)+1; i++ {
				c.memory[i+int(c.i)] = c.v[i]
			}
			c.pc += 2
		case 0x0065:
			for i := 0; i < int(c.opcode&0x0F00>>8)+1; i++ {
				c.v[i] = c.memory[c.i+uint16(i)]
			}
			c.pc += 2
		default:
			log.Printf("Invalid opcode %X\n", c.opcode)
		}

	default:
		log.Printf("Invalid opcode %X\n", c.opcode)
	}
}
