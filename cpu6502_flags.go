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

func (cpu *Cpu6502) testAndSetNegative(b uint8) {
	if b&0x80 == 0x80 {
		cpu.SetFlag(P_NEGATIVE)
	} else {
		cpu.ClearFlag(P_NEGATIVE)
	}
}

func (cpu *Cpu6502) testAndSetZero(b uint8) {
	if b == 0x00 {
		cpu.SetFlag(P_ZERO)
	} else {
		cpu.ClearFlag(P_ZERO)
	}
}

func (cpu *Cpu6502) testAndSetNZ(b uint8) {
	cpu.SetFlagValue(P_NEGATIVE, b&0x80 != 0)
	cpu.SetFlagValue(P_ZERO, b == 0)
}
