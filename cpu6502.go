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

import (
	"fmt"
	"log"
	"strings"
)

type Cpu6502 struct {
	A  uint8  // Accumulator
	X  uint8  // X index register
	Y  uint8  // Y index register
	SP uint8  // Stack pointer
	P  uint8  // Status register
	PC uint16 // PC register

	Debug bool

	waitCycles int

	// opCodes [0xff]func()
	opCodes [0xff]*opcode

	bus *Bus
}

type Flag uint8

const (
	P_CARRY Flag = 1 << iota
	P_ZERO
	P_DISABLE_IRQ
	P_DECIMAL_MODE
	P_BRK_COMMAND
	P_UNUSED
	P_OVERFLOW
	P_NEGATIVE
)

func NewCpu6502(bus *Bus) *Cpu6502 {
	cpu := &Cpu6502{
		bus: bus,
	}
	cpu.initOpcodes()

	bus.Ticker(cpu)

	cpu.Reset()

	return cpu
}

func (cpu *Cpu6502) Reset() {
	cpu.A = 0xaa
	cpu.X = 0
	cpu.Y = 0
	cpu.SP = 0xfd
	cpu.P = 0b00110000
	cpu.PC = readWord(cpu.bus, 0xfffc)

	cpu.waitCycles = 0
}

func (cpu *Cpu6502) Tick() error {
	if cpu.waitCycles > 0 {
		cpu.waitCycles--
		return nil
	}

	// TODO: Check interrupts

	var debugStr string
	if cpu.Debug {
		asm, _ := cpu.Disassemble(cpu.PC)
		debugStr = fmt.Sprintf("%04x    %-26s", cpu.PC, asm)
	}

	opcode := cpu.bus.Read(cpu.PC)
	cpu.PC++

	op := cpu.opCodes[opcode]
	if op == nil {
		log.Fatalf("Invalid opcode 0x%02x at 0x%04x", opcode, cpu.PC-1)
		return nil
	}
	argAddr := cpu.PC
	var arg uint16
	if op.args == 1 {
		arg = uint16(cpu.bus.Read(argAddr))
	} else if op.args == 2 {
		arg = readWord(cpu.bus, argAddr)
	}
	cpu.PC += uint16(op.args)
	addr := op.addr(cpu, argAddr, arg)
	op.exec(addr)
	cpu.waitCycles += op.wait

	if cpu.Debug {
		// PC:0014 A:39 X:14 Y:ec SP:fb nv‑BdIzc
		var statusReg string
		if cpu.FlagSet(P_NEGATIVE) {
			statusReg += "N"
		} else {
			statusReg += "n"
		}
		if cpu.FlagSet(P_OVERFLOW) {
			statusReg += "V-B"
		} else {
			statusReg += "v-B"
		}
		if cpu.FlagSet(P_DECIMAL_MODE) {
			statusReg += "D"
		} else {
			statusReg += "d"
		}
		if cpu.FlagSet(P_DISABLE_IRQ) {
			statusReg += "I"
		} else {
			statusReg += "i"
		}
		if cpu.FlagSet(P_ZERO) {
			statusReg += "Z"
		} else {
			statusReg += "z"
		}
		if cpu.FlagSet(P_CARRY) {
			statusReg += "C"
		} else {
			statusReg += "c"
		}
		fmt.Printf(
			"%sPC:%04x A:%02x X:%02x Y:%02x SP:%02x nv‑BdIzc\n",
			debugStr,
			cpu.PC,
			cpu.A,
			cpu.X,
			cpu.Y,
			cpu.SP,
		)

	}
	return nil
}

func (cpu *Cpu6502) SetFlag(flag Flag)      { cpu.P = cpu.P | uint8(flag) }
func (cpu *Cpu6502) ClearFlag(flag Flag)    { cpu.P = cpu.P &^ uint8(flag) }
func (cpu *Cpu6502) FlagSet(flag Flag) bool { return cpu.P&uint8(flag) != 0x00 }

func (cpu *Cpu6502) Waiting() bool {
	return cpu.waitCycles > 0
}

func (cpu *Cpu6502) Disassemble(addr uint16) (string, uint16) {
	opcode := cpu.bus.Read(addr)
	op := cpu.opCodes[opcode]

	if op == nil {
		log.Fatalf("Invalid opcode 0x%02x at 0x%04x", opcode, cpu.PC-1)
		return "", 0
	}

	var arg uint16
	if op.args == 1 {
		arg = uint16(cpu.bus.Read(addr + 1))
	} else if op.args == 2 {
		arg = readWord(cpu.bus, addr+1)
	}

	return strings.TrimSpace(op.mne + " " + op.fmt(addr+1, arg)), uint16(op.args + 1)
}

func readWord(dev Addressable, a uint16) uint16 {
	return uint16(uint16(dev.Read(a))) + (uint16(dev.Read(a+1)) << 8)
}
