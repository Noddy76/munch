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
	"os"
	"testing"
	"time"
)

func TestFunctionalTest(t *testing.T) {
	bus := NewBus()
	bus.Addressable(0x0000, 0xffff, NewRam(0x10000))
	cpu := NewCpu6502(bus)

	dat, err := os.ReadFile("functional_tests/6502_functional_test.bin")
	if err != nil {
		t.Fatal("Unable to open test binary")
	}
	if len(dat) > 0x10000 {
		log.Fatalf("test binary too large (%d bytes)", len(dat))
	}

	for a, b := range dat {
		bus.Write(uint16(a), b)
	}

	cpu.PC = 0x0400

	start := time.Now()

	pc := cpu.PC
	for pc != 0x3469 {
		for cpu.Waiting() {
			bus.Tick()
		}
		bus.Tick()
		if pc == cpu.PC {
			regs := fmt.Sprintf(
				"A: %02x, X: %02x, Y: %02x, SP: %02x, PC: %04x, P: %08b",
				cpu.A,
				cpu.X,
				cpu.Y,
				cpu.SP,
				cpu.PC,
				cpu.P,
			)
			t.Fatalf("Test trapped at $%04x : %s", cpu.PC, regs)
		}
		pc = cpu.PC
	}

	duration := time.Since(start).Seconds()
	cycles := bus.TickCount()

	t.Logf(
		"%d cycles in %f seconds (%0.2f MHz)",
		cycles,
		duration,
		float64(cycles)/(duration*1000000),
	)
}

func TestSbcZeroMinusOne(t *testing.T) {
	bus := NewBus()
	bus.Addressable(0x0600, 0xffff, NewRom([]uint8{0xa9, 0x00, 0xe9, 0x01}))
	cpu := NewCpu6502(bus)
	cpu.PC = 0x0600

	for cpu.PC != 0x0604 {
		bus.Tick()
	}

	if cpu.FlagSet(P_OVERFLOW) {
		t.Fatal("Overflow flag was set")
	}

	if !cpu.FlagSet(P_NEGATIVE) {
		t.Fatal("Negative flag was not set")
	}

	if cpu.FlagSet(P_ZERO) {
		t.Fatal("Zero flag was set")
	}

	if cpu.FlagSet(P_CARRY) {
		t.Fatal("Carry flag was set")
	}

	if cpu.A != 0xfe {
		t.Fatalf("A is $%02x not $fe", cpu.A)
	}
}
