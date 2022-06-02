// Copyright (C) 2022 James Grant
//
// This is part of munch as 6502 emulator
//
// Munch is free software: you can redistribute it and/or modify it under the terms of the GNU
// General Public License as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// Munch is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even
// the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with Munch. If not, see
// <https://www.gnu.org/licenses/>.

package munch

func (cpu *Cpu6502) brk() {
	cpu.PC += 1
	cpu.stackPushWord(cpu.PC)
	cpu.php()
	cpu.sei()
	cpu.PC = readWord(cpu.bus, 0xfffe)
}

func (cpu *Cpu6502) and(addr uint16) {
	cpu.A &= cpu.bus.Read(addr)

	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu6502) ora(addr uint16) {
	cpu.A |= cpu.bus.Read(addr)

	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu6502) eor(addr uint16) {
	cpu.A ^= cpu.bus.Read(addr)

	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu6502) dec(addr uint16) {
	v := cpu.bus.Read(addr)
	v--
	cpu.bus.Write(addr, v)

	cpu.testAndSetNegative(v)
	cpu.testAndSetZero(v)
}

func (cpu *Cpu6502) inc(addr uint16) {
	v := cpu.bus.Read(addr)
	v++
	cpu.bus.Write(addr, v)

	cpu.testAndSetNegative(v)
	cpu.testAndSetZero(v)
}

func (cpu *Cpu6502) inx() {
	cpu.X += 1
	cpu.testAndSetNegative(cpu.X)
	cpu.testAndSetZero(cpu.X)
}

func (cpu *Cpu6502) iny() {
	cpu.Y += 1
	cpu.testAndSetNegative(cpu.Y)
	cpu.testAndSetZero(cpu.Y)
}

func (cpu *Cpu6502) dex() {
	cpu.X -= 1
	cpu.testAndSetNegative(cpu.X)
	cpu.testAndSetZero(cpu.X)
}

func (cpu *Cpu6502) dey() {
	cpu.Y -= 1
	cpu.testAndSetNegative(cpu.Y)
	cpu.testAndSetZero(cpu.Y)
}

func (cpu *Cpu6502) compare(a, b uint8) {
	v := a - b

	cpu.testAndSetZero(v)
	cpu.testAndSetNegative(v)
	if a >= v {
		cpu.SetFlag(P_CARRY)
	} else {
		cpu.ClearFlag(P_CARRY)
	}
}

func (cpu *Cpu6502) cmp(addr uint16) {
	v := cpu.bus.Read(addr)
	cpu.compare(cpu.A, v)
}

func (cpu *Cpu6502) cpx(addr uint16) {
	v := cpu.bus.Read(addr)
	cpu.compare(cpu.X, v)
}

func (cpu *Cpu6502) cpy(addr uint16) {
	v := cpu.bus.Read(addr)
	cpu.compare(cpu.Y, v)
}

func (cpu *Cpu6502) adc(addr uint16) {
	v := cpu.bus.Read(addr)
	c := uint8(0)
	if cpu.FlagSet(P_CARRY) {
		c = 1
	}

	if cpu.FlagSet(P_DECIMAL_MODE) {
		l := uint16(cpu.A&0x0f) + uint16(v&0x0f) + uint16(c)
		h := uint16(cpu.A&0xf0) + uint16(v&0xf0)

		cpu.ClearFlag(P_OVERFLOW)
		cpu.ClearFlag(P_CARRY)
		cpu.ClearFlag(P_NEGATIVE)
		cpu.ClearFlag(P_ZERO)

		if (l+h)&0xff == 0 {
			cpu.SetFlag(P_ZERO)
		}
		if l > 9 {
			h += 0x10
			l += 0x06
		}
		if h&0x80 != 0 {
			cpu.SetFlag(P_NEGATIVE)
		}
		if ^(uint16(cpu.A)^uint16(v))&(uint16(cpu.A)^h)&0x80 != 0 {
			cpu.SetFlag(P_OVERFLOW)
		}
		if h > 0x90 {
			h += 0x60
		}
		if h&0xff00 != 0 {
			cpu.SetFlag(P_CARRY)
		}

		cpu.A = uint8(l&0xf) | uint8(h&0xf0)
	} else {
		t := uint16(cpu.A) + uint16(v) + uint16(c)

		cpu.ClearFlag(P_OVERFLOW)
		cpu.ClearFlag(P_CARRY)

		if ^(cpu.A^v)&(cpu.A^uint8(t))&0x80 != 0 {
			cpu.SetFlag(P_OVERFLOW)
		}
		if t&0xff00 != 0 {
			cpu.SetFlag(P_CARRY)
		}

		cpu.A = uint8(t)
		cpu.testAndSetNegative(cpu.A)
		cpu.testAndSetZero(cpu.A)
	}
}

func (cpu *Cpu6502) sbc(addr uint16) {
	v := cpu.bus.Read(addr)
	c := uint16(1)
	if cpu.FlagSet(P_CARRY) {
		c = 0
	}
	t := uint16(cpu.A) - uint16(v) - c

	if cpu.FlagSet(P_DECIMAL_MODE) {
		l := uint16(cpu.A&0x0f) - uint16(v&0x0f) - uint16(c)
		h := uint16(cpu.A&0xf0) - uint16(v&0xf0)

		cpu.ClearFlag(P_OVERFLOW)
		cpu.ClearFlag(P_CARRY)
		cpu.ClearFlag(P_NEGATIVE)
		cpu.ClearFlag(P_ZERO)

		if l&0x10 != 0 {
			l -= 6
			h -= 1
		}
		if (uint16(cpu.A)^uint16(v))&(uint16(cpu.A)^t)&0x80 != 0 {
			cpu.SetFlag(P_OVERFLOW)
		}
		if t&0xff00 == 0 {
			cpu.SetFlag(P_CARRY)
		}
		if t&0x00ff == 0 {
			cpu.SetFlag(P_ZERO)
		}
		if t&0x80 != 0 {
			cpu.SetFlag(P_NEGATIVE)
		}
		if h&0x0100 != 0 {
			h -= 0x60
		}

		cpu.A = uint8(l&0x0f) | uint8(h&0xf0)
	} else {
		cpu.ClearFlag(P_OVERFLOW)
		cpu.ClearFlag(P_CARRY)

		if (cpu.A^v)&(cpu.A^uint8(t))&0x80 != 0 {
			cpu.SetFlag(P_OVERFLOW)
		}
		if t&0xff00 == 0 {
			cpu.SetFlag(P_CARRY)
		}

		cpu.A = uint8(t)
		cpu.testAndSetNegative(cpu.A)
		cpu.testAndSetZero(cpu.A)
	}
}

func (cpu *Cpu6502) clc() {
	cpu.ClearFlag(P_CARRY)
}

func (cpu *Cpu6502) sec() {
	cpu.SetFlag(P_CARRY)
}

func (cpu *Cpu6502) cli() {
	cpu.ClearFlag(P_DISABLE_IRQ)
}

func (cpu *Cpu6502) sei() {
	cpu.SetFlag(P_DISABLE_IRQ)
}

func (cpu *Cpu6502) clv() {
	cpu.ClearFlag(P_OVERFLOW)
}

func (cpu *Cpu6502) cld() {
	cpu.ClearFlag(P_DECIMAL_MODE)
}

func (cpu *Cpu6502) sed() {
	cpu.SetFlag(P_DECIMAL_MODE)
}

func (cpu *Cpu6502) bpl(addr uint16) {
	if !cpu.FlagSet(P_NEGATIVE) {
		if (cpu.PC-1)&0xff00 != addr&0xff00 {
			cpu.waitCycles += 1
		}
		cpu.PC = addr
	}
}

func (cpu *Cpu6502) bmi(addr uint16) {
	if cpu.FlagSet(P_NEGATIVE) {
		if (cpu.PC-1)&0xff00 != addr&0xff00 {
			cpu.waitCycles += 1
		}
		cpu.PC = addr
	}
}

func (cpu *Cpu6502) bvc(addr uint16) {
	if !cpu.FlagSet(P_OVERFLOW) {
		if (cpu.PC-1)&0xff00 != addr&0xff00 {
			cpu.waitCycles += 1
		}
		cpu.PC = addr
	}
}

func (cpu *Cpu6502) bvs(addr uint16) {
	if cpu.FlagSet(P_OVERFLOW) {
		if (cpu.PC-1)&0xff00 != addr&0xff00 {
			cpu.waitCycles += 1
		}
		cpu.PC = addr
	}
}

func (cpu *Cpu6502) bcc(addr uint16) {
	if !cpu.FlagSet(P_CARRY) {
		if (cpu.PC-1)&0xff00 != addr&0xff00 {
			cpu.waitCycles += 1
		}
		cpu.PC = addr
	}
}

func (cpu *Cpu6502) bcs(addr uint16) {
	if cpu.FlagSet(P_CARRY) {
		if (cpu.PC-1)&0xff00 != addr&0xff00 {
			cpu.waitCycles += 1
		}
		cpu.PC = addr
	}
}

func (cpu *Cpu6502) bne(addr uint16) {
	if cpu.P&0x02 == 0x00 {
		if (cpu.PC-1)&0xff00 != addr&0xff00 {
			cpu.waitCycles += 1
		}
		cpu.PC = addr
	}
}

func (cpu *Cpu6502) beq(addr uint16) {
	if cpu.FlagSet(P_ZERO) {
		if (cpu.PC-1)&0xff00 != addr&0xff00 {
			cpu.waitCycles += 1
		}
		cpu.PC = addr
	}
}

func (cpu *Cpu6502) jmp(addr uint16) {
	cpu.PC = addr
}

func (cpu *Cpu6502) lda(addr uint16) {
	cpu.A = cpu.bus.Read(addr)
	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu6502) ldx(addr uint16) {
	cpu.X = cpu.bus.Read(addr)
	cpu.testAndSetNegative(cpu.X)
	cpu.testAndSetZero(cpu.X)
}

func (cpu *Cpu6502) ldy(addr uint16) {
	cpu.Y = cpu.bus.Read(addr)
	cpu.testAndSetNegative(cpu.Y)
	cpu.testAndSetZero(cpu.Y)
}

func (cpu *Cpu6502) sta(addr uint16) {
	cpu.bus.Write(addr, cpu.A)
}

func (cpu *Cpu6502) stx(addr uint16) {
	cpu.bus.Write(addr, cpu.X)
}

func (cpu *Cpu6502) sty(addr uint16) {
	cpu.bus.Write(addr, cpu.Y)
}

func (cpu *Cpu6502) tax() {
	cpu.X = cpu.A
	cpu.testAndSetNegative(cpu.X)
	cpu.testAndSetZero(cpu.X)
}

func (cpu *Cpu6502) txa() {
	cpu.A = cpu.X
	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu6502) txs() {
	cpu.SP = cpu.X
}

func (cpu *Cpu6502) tay() {
	cpu.Y = cpu.A
	cpu.testAndSetNegative(cpu.Y)
	cpu.testAndSetZero(cpu.Y)
}

func (cpu *Cpu6502) tya() {
	cpu.A = cpu.Y
	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu6502) tsx() {
	cpu.X = cpu.SP
	cpu.testAndSetZero(cpu.X)
	cpu.testAndSetNegative(cpu.X)
}

func (cpu *Cpu6502) pha() {
	cpu.stackPush(cpu.A)
}

func (cpu *Cpu6502) pla() {
	cpu.A = cpu.stackPop()
	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu6502) php() {
	// BRK and PHP push P OR #$10, so that the IRQ handler can tell
	// whether the entry was from a BRK or from an /IRQ.
	cpu.stackPush(cpu.P | uint8(P_BRK_COMMAND) | uint8(P_UNUSED))
}

func (cpu *Cpu6502) plp() {
	cpu.P = cpu.stackPop() | uint8(P_BRK_COMMAND) | uint8(P_UNUSED)
}

func (cpu *Cpu6502) jsr(addr uint16) {
	cpu.stackPushWord(cpu.PC - 1)
	cpu.PC = addr
}

func (cpu *Cpu6502) rti() {
	cpu.plp()

	low := cpu.stackPop()
	high := cpu.stackPop()

	cpu.PC = uint16(high)<<8 + uint16(low)
}

func (cpu *Cpu6502) rts() {
	low := cpu.stackPop()
	high := cpu.stackPop()

	cpu.PC = uint16(high)<<8 + uint16(low) + 1
}

func (cpu *Cpu6502) lsr(addr uint16) {
	v := cpu.bus.Read(addr)

	if v&0x01 == 0x01 {
		cpu.SetFlag(P_CARRY)
	} else {
		cpu.ClearFlag(P_CARRY)
	}

	cpu.bus.Write(addr, v>>1)

	v = cpu.bus.Read(addr)

	cpu.testAndSetNegative(v)
	cpu.testAndSetZero(v)
}
func (cpu *Cpu6502) lsrAcc() {
	if cpu.A&0x01 == 0x01 {
		cpu.SetFlag(P_CARRY)
	} else {
		cpu.ClearFlag(P_CARRY)
	}

	cpu.A >>= 1

	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu6502) asl(addr uint16) {
	v := cpu.bus.Read(addr)

	if v&0x80 == 0x80 {
		cpu.SetFlag(P_CARRY)
	} else {
		cpu.ClearFlag(P_CARRY)
	}

	cpu.bus.Write(addr, v<<1)

	v = cpu.bus.Read(addr)

	cpu.testAndSetNegative(v)
	cpu.testAndSetZero(v)
}
func (cpu *Cpu6502) aslAcc() {
	if cpu.A&0x80 == 0x80 {
		cpu.SetFlag(P_CARRY)
	} else {
		cpu.ClearFlag(P_CARRY)
	}

	cpu.A <<= 1

	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu6502) rol(addr uint16) {
	v := cpu.bus.Read(addr)

	carry := v&0x80 == 0x80

	v <<= 1

	if cpu.FlagSet(P_CARRY) {
		v += 1
	}

	if carry {
		cpu.SetFlag(P_CARRY)
	} else {
		cpu.ClearFlag(P_CARRY)
	}

	cpu.bus.Write(addr, v)
	v = cpu.bus.Read(addr)

	cpu.testAndSetNegative(v)
	cpu.testAndSetZero(v)
}
func (cpu *Cpu6502) rolAcc() {
	carry := cpu.A&0x80 == 0x80

	cpu.A <<= 1

	if cpu.FlagSet(P_CARRY) {
		cpu.A += 1
	}

	if carry {
		cpu.SetFlag(P_CARRY)
	} else {
		cpu.ClearFlag(P_CARRY)
	}

	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu6502) ror(addr uint16) {
	v := cpu.bus.Read(addr)

	carry := v&0x01 == 0x01

	v >>= 1

	if cpu.FlagSet(P_CARRY) {
		v += 0x80
	}

	if carry {
		cpu.SetFlag(P_CARRY)
	} else {
		cpu.ClearFlag(P_CARRY)
	}

	cpu.bus.Write(addr, v)
	v = cpu.bus.Read(addr)

	cpu.testAndSetNegative(v)
	cpu.testAndSetZero(v)
}
func (cpu *Cpu6502) rorAcc() {
	carry := cpu.A&0x01 == 0x01

	cpu.A >>= 1

	if cpu.FlagSet(P_CARRY) {
		cpu.A += 0x80
	}

	if carry {
		cpu.SetFlag(P_CARRY)
	} else {
		cpu.ClearFlag(P_CARRY)
	}

	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu6502) bit(addr uint16) {
	v := cpu.bus.Read(addr)

	cpu.testAndSetZero(v & cpu.A)
	cpu.testAndSetNegative(v)

	if v&0x40 == 0x40 {
		cpu.SetFlag(P_OVERFLOW)
	} else {
		cpu.ClearFlag(P_OVERFLOW)
	}
}

func (cpu *Cpu6502) stackPush(a uint8) {
	cpu.bus.Write(uint16(0x100)+uint16(cpu.SP), a)
	cpu.SP -= 1
}

func (cpu *Cpu6502) stackPushWord(a uint16) {
	cpu.stackPush(uint8(a >> 8))
	cpu.stackPush(uint8(a & 0xff))
}

func (cpu *Cpu6502) stackPop() uint8 {
	cpu.SP += 1
	return cpu.bus.Read(uint16(0x100) + uint16(cpu.SP))
}
