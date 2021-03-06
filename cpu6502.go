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

	Debug          bool
	DisableDecimal bool

	waitCycles int

	opCodes [0x100]*opcode

	bus *Bus

	pendingIrq bool
	pendingNmi bool
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
	cpu.P = 0b00110100
	cpu.PC = readWord(cpu.bus, 0xfffc)

	cpu.waitCycles = 0
}

func (cpu *Cpu6502) Irq() {
	if !cpu.FlagSet(P_DISABLE_IRQ) {
		cpu.pendingIrq = true
	}
}

func (cpu *Cpu6502) Nmi() {
	cpu.pendingNmi = true
}

func (cpu *Cpu6502) Tick() error {
	if cpu.waitCycles > 0 {
		cpu.waitCycles--
		return nil
	}

	if cpu.pendingIrq || cpu.pendingNmi {
		cpu.stackPushWord(cpu.PC + 1)
		cpu.stackPush(cpu.P &^ uint8(P_BRK_COMMAND))
		cpu.SetFlag(P_DISABLE_IRQ)
		if cpu.pendingNmi {
			cpu.PC = readWord(cpu.bus, 0xfffa)
		} else {
			cpu.PC = readWord(cpu.bus, 0xfffe)
		}
		cpu.waitCycles = 7
		cpu.pendingIrq = false
		cpu.pendingNmi = false
		return nil
	}

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
	if op.addrMode.args == 1 {
		arg = uint16(cpu.bus.Read(argAddr))
	} else if op.addrMode.args == 2 {
		arg = readWord(cpu.bus, argAddr)
	}
	cpu.PC += uint16(op.addrMode.args)
	addr := op.addrMode.addr(cpu, argAddr, arg)
	op.exec(addr)
	cpu.waitCycles += op.wait

	if cpu.Debug {
		fmt.Printf("%s%s\n", debugStr, cpu.StatusString())
	}
	return nil
}

func (cpu *Cpu6502) SetFlag(flag Flag)      { cpu.P = cpu.P | uint8(flag) }
func (cpu *Cpu6502) ClearFlag(flag Flag)    { cpu.P = cpu.P &^ uint8(flag) }
func (cpu *Cpu6502) FlagSet(flag Flag) bool { return cpu.P&uint8(flag) != 0x00 }
func (cpu *Cpu6502) SetFlagValue(flag Flag, v bool) {
	if v {
		cpu.SetFlag(flag)
	} else {
		cpu.ClearFlag(flag)
	}
}

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
	if op.addrMode.args == 1 {
		arg = uint16(cpu.bus.Read(addr + 1))
	} else if op.addrMode.args == 2 {
		arg = readWord(cpu.bus, addr+1)
	}

	return strings.TrimSpace(op.mne + " " + op.addrMode.fmt(addr+1, arg)), uint16(op.addrMode.args + 1)
}

func (cpu *Cpu6502) StatusString() string {
	low := "nv-bdizc"
	var statusReg string
	for i := 7; i >= 0; i-- {
		if (1<<i)&cpu.P != 0 {
			statusReg += strings.ToUpper(string(low[7-i]))
		} else {
			statusReg += strings.ToLower(string(low[7-i]))
		}
	}
	return fmt.Sprintf(
		"PC:%04x A:%02x X:%02x Y:%02x SP:%02x %s",
		cpu.PC,
		cpu.A,
		cpu.X,
		cpu.Y,
		cpu.SP, statusReg,
	)

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

func (cpu *Cpu6502) stackPopWord() uint16 {
	low := cpu.stackPop()
	high := cpu.stackPop()

	return uint16(high)<<8 + uint16(low)
}

func readWord(dev Addressable, a uint16) uint16 {
	return uint16(uint16(dev.Read(a))) + (uint16(dev.Read(a+1)) << 8)
}
