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

func (cpu *Cpu6502) nop(addr uint16) {}

func (c *Cpu6502) slo(a uint16) { c.asl(a); c.ora(a) }
func (c *Cpu6502) rla(a uint16) { c.rol(a); c.and(a) }
func (c *Cpu6502) sre(a uint16) { c.lsr(a); c.eor(a) }
func (c *Cpu6502) rra(a uint16) { c.ror(a); c.adc(a) }
func (c *Cpu6502) sax(a uint16) { c.bus.Write(a, c.A&c.X) }
func (c *Cpu6502) lax(a uint16) { c.lda(a); c.ldx(a) }
func (c *Cpu6502) dcp(a uint16) { c.dec(a); c.cmp(a) }
func (c *Cpu6502) isc(a uint16) { c.inc(a); c.sbc(a) }

func (cpu *Cpu6502) brk(a uint16) {
	cpu.PC += 1
	cpu.stackPushWord(cpu.PC)
	cpu.php(a)
	cpu.sei(a)
	cpu.PC = readWord(cpu.bus, 0xfffe)
}

func (cpu *Cpu6502) and(addr uint16) {
	cpu.A &= cpu.bus.Read(addr)

	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu6502) anc(addr uint16) {
	cpu.and(addr)
	cpu.SetFlagValue(P_CARRY, cpu.FlagSet(P_NEGATIVE))
}

func (cpu *Cpu6502) alr(addr uint16) {
	cpu.and(addr)
	cpu.lsrAcc(addr)
}

func (cpu *Cpu6502) arr(addr uint16) {
	// ARR is weird see http://www.cs.cmu.edu/~dsladic/vice/doc/64doc.txt
	if cpu.FlagSet(P_DECIMAL_MODE) && !cpu.DisableDecimal {
		s := cpu.bus.Read(addr)
		t := cpu.A & s /* Perform the AND. */

		AH := t >> 4 /* Separate the high */
		AL := t & 15 /* and low nybbles. */

		cpu.A = t >> 1
		if cpu.FlagSet(P_CARRY) {
			cpu.A |= 0x80
		}

		cpu.SetFlagValue(P_NEGATIVE, cpu.FlagSet(P_CARRY)) /* Set the N and */
		cpu.SetFlagValue(P_ZERO, cpu.A == 0)               /* Z flags traditionally */
		cpu.SetFlagValue(P_OVERFLOW, (t^cpu.A)&64 != 0)    /* and V flag in a weird way. */

		if AL+(AL&1) > 5 { /* BCD "fixup" for low nybble. */
			cpu.A = (cpu.A & 0xF0) | ((cpu.A + 6) & 0xF)
		}
		cpu.SetFlagValue(P_CARRY, AH+(AH&1) > 5)
		if cpu.FlagSet(P_CARRY) { /* Set the Carry flag. */
			cpu.A = (cpu.A + 0x60) & 0xFF
		} /* BCD "fixup" for high nybble. */
	} else {
		// In Binary mode (D flag clear), the instruction effectively does an AND
		// between the accumulator and the immediate parameter, and then shifts
		// the accumulator to the right, copying the C flag to the 8th bit. It
		// sets the Negative and Zero flags just like the ROR would. The ADC code
		// shows up in the Carry and oVerflow flags. The C flag will be copied
		// from the bit 6 of the result (which doesn't seem too logical), and the
		// V flag is the result of an Exclusive OR operation between the bit 6
		// and the bit 5 of the result.  This makes sense, since the V flag will
		// be normally set by an Exclusive OR, too.
		cpu.and(addr)
		cpu.A >>= 1
		if cpu.FlagSet(P_CARRY) {
			cpu.A |= 0b1000_0000
		}
		cpu.testAndSetNZ(cpu.A)
		cpu.SetFlagValue(P_CARRY, cpu.A&0b0100_0000 != 0)
		cpu.SetFlagValue(P_OVERFLOW, (cpu.A&0x40 != 0) != (cpu.A&0x20 != 0))
	}
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

func (cpu *Cpu6502) inx(addr uint16) {
	cpu.X += 1
	cpu.testAndSetNegative(cpu.X)
	cpu.testAndSetZero(cpu.X)
}

func (cpu *Cpu6502) iny(addr uint16) {
	cpu.Y += 1
	cpu.testAndSetNegative(cpu.Y)
	cpu.testAndSetZero(cpu.Y)
}

func (cpu *Cpu6502) dex(addr uint16) {
	cpu.X -= 1
	cpu.testAndSetNegative(cpu.X)
	cpu.testAndSetZero(cpu.X)
}

func (cpu *Cpu6502) dey(addr uint16) {
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

func (cpu *Cpu6502) axs(addr uint16) {
	data := cpu.bus.Read(addr)
	lhs := cpu.A & cpu.X
	cpu.X = lhs - data
	cpu.compare(lhs, data)
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

	if cpu.FlagSet(P_DECIMAL_MODE) && !cpu.DisableDecimal {
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

	if cpu.FlagSet(P_DECIMAL_MODE) && !cpu.DisableDecimal {
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

func (cpu *Cpu6502) clc(addr uint16) {
	cpu.ClearFlag(P_CARRY)
}

func (cpu *Cpu6502) sec(addr uint16) {
	cpu.SetFlag(P_CARRY)
}

func (cpu *Cpu6502) cli(addr uint16) {
	cpu.ClearFlag(P_DISABLE_IRQ)
}

func (cpu *Cpu6502) sei(addr uint16) {
	cpu.SetFlag(P_DISABLE_IRQ)
}

func (cpu *Cpu6502) clv(addr uint16) {
	cpu.ClearFlag(P_OVERFLOW)
}

func (cpu *Cpu6502) cld(addr uint16) {
	cpu.ClearFlag(P_DECIMAL_MODE)
}

func (cpu *Cpu6502) sed(addr uint16) {
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

func (cpu *Cpu6502) tax(addr uint16) {
	cpu.X = cpu.A
	cpu.testAndSetNegative(cpu.X)
	cpu.testAndSetZero(cpu.X)
}

func (cpu *Cpu6502) txa(addr uint16) {
	cpu.A = cpu.X
	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu6502) txs(addr uint16) {
	cpu.SP = cpu.X
}

func (cpu *Cpu6502) tay(addr uint16) {
	cpu.Y = cpu.A
	cpu.testAndSetNegative(cpu.Y)
	cpu.testAndSetZero(cpu.Y)
}

func (cpu *Cpu6502) tya(addr uint16) {
	cpu.A = cpu.Y
	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu6502) tsx(addr uint16) {
	cpu.X = cpu.SP
	cpu.testAndSetZero(cpu.X)
	cpu.testAndSetNegative(cpu.X)
}

func (cpu *Cpu6502) pha(addr uint16) {
	cpu.stackPush(cpu.A)
}

func (cpu *Cpu6502) pla(addr uint16) {
	cpu.A = cpu.stackPop()
	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu6502) php(a uint16) {
	// BRK and PHP push P OR #$10, so that the IRQ handler can tell
	// whether the entry was from a BRK or from an /IRQ.
	cpu.stackPush(cpu.P | uint8(P_BRK_COMMAND) | uint8(P_UNUSED))
}

func (cpu *Cpu6502) plp(a uint16) {
	cpu.P = cpu.stackPop() | uint8(P_BRK_COMMAND) | uint8(P_UNUSED)
}

func (cpu *Cpu6502) jsr(addr uint16) {
	cpu.stackPushWord(cpu.PC - 1)
	cpu.PC = addr
}

func (cpu *Cpu6502) rti(addr uint16) {
	cpu.plp(addr)

	low := cpu.stackPop()
	high := cpu.stackPop()

	cpu.PC = uint16(high)<<8 + uint16(low)
}

func (cpu *Cpu6502) rts(addr uint16) {
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
func (cpu *Cpu6502) lsrAcc(a uint16) {
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
func (cpu *Cpu6502) aslAcc(a uint16) {
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
func (cpu *Cpu6502) rolAcc(a uint16) {
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
func (cpu *Cpu6502) rorAcc(a uint16) {
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

func (c *Cpu6502) shy(a uint16) {
	high := (a & 0xff00) >> 8
	a &= 0x00ff
	if a+uint16(c.X) > 0x00ff {
		a += (((high & uint16(c.Y)) << 8) + uint16(c.X))
	} else {
		a += ((high << 8) + uint16(c.X))
	}
	c.bus.Write(a, c.Y&(uint8(high)+1))
}

func (c *Cpu6502) shx(a uint16) {
	high := (a & 0xff00) >> 8
	a &= 0x00ff
	if a+uint16(c.Y) > 0x00ff {
		a += (((high & uint16(c.X)) << 8) + uint16(c.Y))
	} else {
		a += ((high << 8) + uint16(c.Y))
	}
	c.bus.Write(a, c.X&(uint8(high)+1))
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
