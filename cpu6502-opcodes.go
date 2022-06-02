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

import "fmt"

type opcode struct {
	mne  string
	wait int
	args int
	fmt  func(uint16, uint16) string           // func(argAddr, arg) => string
	addr func(*Cpu6502, uint16, uint16) uint16 // func(argAddr, arg) => address
	exec func(uint16)
}

func nilFmt(argAddr, arg uint16) string { return "" }
func immFmt(argAddr, arg uint16) string { return fmt.Sprintf("#$%02x", arg) }
func zpFmt(argAddr, arg uint16) string  { return fmt.Sprintf("$%02x", arg) }
func zpxFmt(argAddr, arg uint16) string { return fmt.Sprintf("$%02x,X", arg) }
func zpyFmt(argAddr, arg uint16) string { return fmt.Sprintf("$%02x,Y", arg) }
func izxFmt(argAddr, arg uint16) string { return fmt.Sprintf("($%02x, X)", arg) }
func izyFmt(argAddr, arg uint16) string { return fmt.Sprintf("($%02x),Y", arg) }
func absFmt(argAddr, arg uint16) string { return fmt.Sprintf("$%04x", arg) }
func abxFmt(argAddr, arg uint16) string { return fmt.Sprintf("$%04x,X", arg) }
func abyFmt(argAddr, arg uint16) string { return fmt.Sprintf("$%04x,Y", arg) }
func indFmt(argAddr, arg uint16) string { return fmt.Sprintf("($%04x)", arg) }
func relFmt(argAddr, arg uint16) string {
	if arg < 0x80 {
		arg += argAddr
	} else {
		arg += argAddr - 0x100
	}
	return fmt.Sprintf("$%04x", arg+1)
}

func nilAddr(cpu *Cpu6502, argAddr, arg uint16) uint16 { return 0 }
func immAddr(cpu *Cpu6502, argAddr, arg uint16) uint16 { return argAddr }
func zpAddr(cpu *Cpu6502, argAddr, arg uint16) uint16  { return arg }
func zpxAddr(cpu *Cpu6502, argAddr, arg uint16) uint16 { return uint16(uint8(arg) + cpu.X) }
func zpyAddr(cpu *Cpu6502, argAddr, arg uint16) uint16 { return uint16(uint8(arg) + cpu.Y) }
func izxAddr(cpu *Cpu6502, argAddr, arg uint16) uint16 {
	return readWord(cpu.bus, uint16(uint8(arg)+cpu.X))
}
func izyAddr(cpu *Cpu6502, argAddr, arg uint16) uint16 {
	return readWord(cpu.bus, arg) + uint16(cpu.Y)
}
func absAddr(cpu *Cpu6502, argAddr, arg uint16) uint16 { return arg }
func abxAddr(cpu *Cpu6502, argAddr, arg uint16) uint16 {
	addr := arg + uint16(cpu.X)
	if arg&0xff00 != addr&0xff00 {
		cpu.waitCycles += 1
	}
	return addr
}
func abyAddr(cpu *Cpu6502, argAddr, arg uint16) uint16 {
	addr := arg + uint16(cpu.Y)
	if arg&0xff00 != addr&0xff00 {
		cpu.waitCycles += 1
	}
	return addr
}
func indAddr(cpu *Cpu6502, argAddr, arg uint16) uint16 { return readWord(cpu.bus, arg) }
func relAddr(cpu *Cpu6502, argAddr, arg uint16) uint16 {
	if arg < 0x80 {
		arg += argAddr
	} else {
		arg += argAddr - 0x100
	}
	return arg + 1
}

func (c *Cpu6502) initOpcodes() {

	ora := func(a uint16) { c.ora(a) }
	brk := func(a uint16) { c.brk() }
	asl := func(a uint16) { c.asl(a) }
	aslA := func(a uint16) { c.aslAcc() }
	php := func(a uint16) { c.php() }
	bpl := func(a uint16) { c.bpl(a) }
	clc := func(a uint16) { c.clc() }
	nop := func(a uint16) {}
	jsr := func(a uint16) { c.jsr(a) }
	and := func(a uint16) { c.and(a) }
	bit := func(a uint16) { c.bit(a) }
	rol := func(a uint16) { c.rol(a) }
	rolA := func(a uint16) { c.rolAcc() }
	plp := func(a uint16) { c.plp() }
	bmi := func(a uint16) { c.bmi(a) }
	sec := func(a uint16) { c.sec() }
	rti := func(a uint16) { c.rti() }
	eor := func(a uint16) { c.eor(a) }
	lsr := func(a uint16) { c.lsr(a) }
	lsrA := func(a uint16) { c.lsrAcc() }
	pha := func(a uint16) { c.pha() }
	jmp := func(a uint16) { c.jmp(a) }
	bvc := func(a uint16) { c.bvc(a) }
	cli := func(a uint16) { c.cli() }
	rts := func(a uint16) { c.rts() }
	adc := func(a uint16) { c.adc(a) }
	sbc := func(a uint16) { c.sbc(a) }
	ror := func(a uint16) { c.ror(a) }
	rorA := func(a uint16) { c.rorAcc() }
	pla := func(a uint16) { c.pla() }
	bvs := func(a uint16) { c.bvs(a) }
	sei := func(a uint16) { c.sei() }
	sta := func(a uint16) { c.sta(a) }
	stx := func(a uint16) { c.stx(a) }
	sty := func(a uint16) { c.sty(a) }
	dey := func(a uint16) { c.dey() }
	txa := func(a uint16) { c.txa() }
	bcc := func(a uint16) { c.bcc(a) }
	tya := func(a uint16) { c.tya() }
	txs := func(a uint16) { c.txs() }
	lda := func(a uint16) { c.lda(a) }
	ldx := func(a uint16) { c.ldx(a) }
	ldy := func(a uint16) { c.ldy(a) }
	tay := func(a uint16) { c.tay() }
	tax := func(a uint16) { c.tax() }
	bcs := func(a uint16) { c.bcs(a) }
	clv := func(a uint16) { c.clv() }
	tsx := func(a uint16) { c.tsx() }
	cpy := func(a uint16) { c.cpy(a) }
	cmp := func(a uint16) { c.cmp(a) }
	dec := func(a uint16) { c.dec(a) }
	dex := func(a uint16) { c.dex() }
	iny := func(a uint16) { c.iny() }
	bne := func(a uint16) { c.bne(a) }
	inc := func(a uint16) { c.inc(a) }
	cld := func(a uint16) { c.cld() }
	cpx := func(a uint16) { c.cpx(a) }
	inx := func(a uint16) { c.inx() }
	beq := func(a uint16) { c.beq(a) }
	sed := func(a uint16) { c.sed() }

	// &opcode{mne: "", wait: , args: , fmt: Fmt, addr: Addr, exec: }

	c.opCodes[0x00] = &opcode{mne: "BRK", wait: 7, args: 0, fmt: nilFmt, addr: nilAddr, exec: brk}
	c.opCodes[0x01] = &opcode{mne: "ORA", wait: 6, args: 1, fmt: izxFmt, addr: izxAddr, exec: ora}
	c.opCodes[0x05] = &opcode{mne: "ORA", wait: 3, args: 1, fmt: zpFmt, addr: zpAddr, exec: ora}
	c.opCodes[0x06] = &opcode{mne: "ASL", wait: 5, args: 1, fmt: zpFmt, addr: zpAddr, exec: asl}
	c.opCodes[0x08] = &opcode{mne: "PHP", wait: 3, args: 0, fmt: nilFmt, addr: nilAddr, exec: php}
	c.opCodes[0x09] = &opcode{mne: "ORA", wait: 2, args: 1, fmt: immFmt, addr: immAddr, exec: ora}
	c.opCodes[0x0a] = &opcode{mne: "ASL", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: aslA}
	c.opCodes[0x0d] = &opcode{mne: "ORA", wait: 4, args: 2, fmt: absFmt, addr: absAddr, exec: ora}
	c.opCodes[0x0e] = &opcode{mne: "ASL", wait: 6, args: 2, fmt: absFmt, addr: absAddr, exec: asl}

	c.opCodes[0x10] = &opcode{mne: "BPL", wait: 2, args: 1, fmt: relFmt, addr: relAddr, exec: bpl}
	c.opCodes[0x11] = &opcode{mne: "ORA", wait: 5, args: 1, fmt: izyFmt, addr: izyAddr, exec: ora}
	c.opCodes[0x15] = &opcode{mne: "ORA", wait: 4, args: 1, fmt: zpxFmt, addr: zpxAddr, exec: ora}
	c.opCodes[0x16] = &opcode{mne: "ASL", wait: 6, args: 1, fmt: zpxFmt, addr: zpxAddr, exec: asl}
	c.opCodes[0x18] = &opcode{mne: "CLC", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: clc}
	c.opCodes[0x19] = &opcode{mne: "ORA", wait: 4, args: 2, fmt: abyFmt, addr: abyAddr, exec: ora}
	c.opCodes[0x1a] = &opcode{mne: "NOP", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: nop}
	c.opCodes[0x1d] = &opcode{mne: "ORA", wait: 4, args: 2, fmt: abxFmt, addr: abxAddr, exec: ora}
	c.opCodes[0x1e] = &opcode{mne: "ASL", wait: 7, args: 2, fmt: abxFmt, addr: abxAddr, exec: asl}

	c.opCodes[0x20] = &opcode{mne: "JSR", wait: 6, args: 2, fmt: absFmt, addr: absAddr, exec: jsr}
	c.opCodes[0x21] = &opcode{mne: "AND", wait: 6, args: 1, fmt: izxFmt, addr: izxAddr, exec: and}
	c.opCodes[0x24] = &opcode{mne: "BIT", wait: 3, args: 1, fmt: zpFmt, addr: zpAddr, exec: bit}
	c.opCodes[0x25] = &opcode{mne: "AND", wait: 3, args: 1, fmt: zpFmt, addr: zpAddr, exec: and}
	c.opCodes[0x26] = &opcode{mne: "ROL", wait: 5, args: 1, fmt: zpFmt, addr: zpAddr, exec: rol}
	c.opCodes[0x28] = &opcode{mne: "PLP", wait: 4, args: 0, fmt: nilFmt, addr: nilAddr, exec: plp}
	c.opCodes[0x29] = &opcode{mne: "AND", wait: 2, args: 1, fmt: immFmt, addr: immAddr, exec: and}
	c.opCodes[0x2a] = &opcode{mne: "ROL", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: rolA}
	c.opCodes[0x2c] = &opcode{mne: "BIT", wait: 4, args: 2, fmt: absFmt, addr: absAddr, exec: bit}
	c.opCodes[0x2d] = &opcode{mne: "AND", wait: 4, args: 2, fmt: absFmt, addr: absAddr, exec: and}
	c.opCodes[0x2e] = &opcode{mne: "ROL", wait: 6, args: 2, fmt: absFmt, addr: absAddr, exec: rol}

	c.opCodes[0x30] = &opcode{mne: "BMI", wait: 2, args: 1, fmt: relFmt, addr: relAddr, exec: bmi}
	c.opCodes[0x31] = &opcode{mne: "AND", wait: 5, args: 1, fmt: izyFmt, addr: izyAddr, exec: and}
	c.opCodes[0x35] = &opcode{mne: "AND", wait: 4, args: 1, fmt: zpxFmt, addr: zpxAddr, exec: and}
	c.opCodes[0x36] = &opcode{mne: "ROL", wait: 6, args: 1, fmt: zpxFmt, addr: zpxAddr, exec: rol}
	c.opCodes[0x38] = &opcode{mne: "SEC", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: sec}
	c.opCodes[0x39] = &opcode{mne: "AND", wait: 4, args: 2, fmt: abyFmt, addr: abyAddr, exec: and}
	c.opCodes[0x3a] = &opcode{mne: "NOP*", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: nop}
	c.opCodes[0x3d] = &opcode{mne: "AND", wait: 4, args: 2, fmt: abxFmt, addr: abxAddr, exec: and}
	c.opCodes[0x3e] = &opcode{mne: "ROL", wait: 7, args: 2, fmt: abxFmt, addr: abxAddr, exec: rol}

	c.opCodes[0x40] = &opcode{mne: "RTI", wait: 6, args: 0, fmt: nilFmt, addr: nilAddr, exec: rti}
	c.opCodes[0x41] = &opcode{mne: "EOR", wait: 6, args: 1, fmt: izxFmt, addr: izxAddr, exec: eor}
	c.opCodes[0x45] = &opcode{mne: "EOR", wait: 3, args: 1, fmt: zpFmt, addr: zpAddr, exec: eor}
	c.opCodes[0x46] = &opcode{mne: "LSR", wait: 5, args: 1, fmt: zpFmt, addr: zpAddr, exec: lsr}
	c.opCodes[0x48] = &opcode{mne: "PHA", wait: 3, args: 0, fmt: nilFmt, addr: nilAddr, exec: pha}
	c.opCodes[0x49] = &opcode{mne: "EOR", wait: 2, args: 1, fmt: immFmt, addr: immAddr, exec: eor}
	c.opCodes[0x4a] = &opcode{mne: "LSR", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: lsrA}
	c.opCodes[0x4c] = &opcode{mne: "JMP", wait: 3, args: 2, fmt: absFmt, addr: absAddr, exec: jmp}
	c.opCodes[0x4d] = &opcode{mne: "EOR", wait: 4, args: 2, fmt: absFmt, addr: absAddr, exec: eor}
	c.opCodes[0x4e] = &opcode{mne: "LSR", wait: 6, args: 2, fmt: absFmt, addr: absAddr, exec: lsr}

	c.opCodes[0x50] = &opcode{mne: "BVC", wait: 2, args: 1, fmt: relFmt, addr: relAddr, exec: bvc}
	c.opCodes[0x51] = &opcode{mne: "EOR", wait: 5, args: 1, fmt: izyFmt, addr: izyAddr, exec: eor}
	c.opCodes[0x55] = &opcode{mne: "EOR", wait: 4, args: 1, fmt: zpxFmt, addr: zpxAddr, exec: eor}
	c.opCodes[0x56] = &opcode{mne: "LSR", wait: 6, args: 1, fmt: zpxFmt, addr: zpxAddr, exec: lsr}
	c.opCodes[0x58] = &opcode{mne: "CLI", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: cli}
	c.opCodes[0x59] = &opcode{mne: "EOR", wait: 4, args: 2, fmt: abyFmt, addr: abyAddr, exec: eor}
	c.opCodes[0x5a] = &opcode{mne: "NOP*", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: nop}
	c.opCodes[0x5d] = &opcode{mne: "EOR", wait: 4, args: 2, fmt: abxFmt, addr: abxAddr, exec: eor}
	c.opCodes[0x5e] = &opcode{mne: "LSR", wait: 7, args: 2, fmt: abxFmt, addr: abxAddr, exec: lsr}

	c.opCodes[0x60] = &opcode{mne: "RTS", wait: 6, args: 0, fmt: nilFmt, addr: nilAddr, exec: rts}
	c.opCodes[0x61] = &opcode{mne: "ADC", wait: 6, args: 1, fmt: izxFmt, addr: izxAddr, exec: adc}
	c.opCodes[0x65] = &opcode{mne: "ADC", wait: 3, args: 1, fmt: zpFmt, addr: zpAddr, exec: adc}
	c.opCodes[0x66] = &opcode{mne: "ROR", wait: 5, args: 1, fmt: zpFmt, addr: zpAddr, exec: ror}
	c.opCodes[0x68] = &opcode{mne: "PLA", wait: 4, args: 0, fmt: nilFmt, addr: nilAddr, exec: pla}
	c.opCodes[0x69] = &opcode{mne: "ADC", wait: 2, args: 1, fmt: immFmt, addr: immAddr, exec: adc}
	c.opCodes[0x6a] = &opcode{mne: "ROR", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: rorA}
	c.opCodes[0x6c] = &opcode{mne: "JMP", wait: 5, args: 2, fmt: indFmt, addr: indAddr, exec: jmp}
	c.opCodes[0x6d] = &opcode{mne: "ADC", wait: 4, args: 2, fmt: absFmt, addr: absAddr, exec: adc}
	c.opCodes[0x6e] = &opcode{mne: "ROR", wait: 6, args: 2, fmt: absFmt, addr: absAddr, exec: ror}

	c.opCodes[0x70] = &opcode{mne: "BVS", wait: 2, args: 1, fmt: relFmt, addr: relAddr, exec: bvs}
	c.opCodes[0x71] = &opcode{mne: "ADC", wait: 5, args: 1, fmt: izyFmt, addr: izyAddr, exec: adc}
	c.opCodes[0x75] = &opcode{mne: "ADC", wait: 4, args: 1, fmt: zpxFmt, addr: zpxAddr, exec: adc}
	c.opCodes[0x76] = &opcode{mne: "ROR", wait: 6, args: 1, fmt: zpxFmt, addr: zpxAddr, exec: ror}
	c.opCodes[0x78] = &opcode{mne: "SEI", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: sei}
	c.opCodes[0x79] = &opcode{mne: "ADC", wait: 4, args: 2, fmt: abyFmt, addr: abyAddr, exec: adc}
	c.opCodes[0x7a] = &opcode{mne: "NOP*", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: nop}
	c.opCodes[0x7d] = &opcode{mne: "ADC", wait: 4, args: 2, fmt: abxFmt, addr: abxAddr, exec: adc}
	c.opCodes[0x7e] = &opcode{mne: "ROR", wait: 7, args: 2, fmt: abxFmt, addr: abxAddr, exec: ror}

	c.opCodes[0x81] = &opcode{mne: "STA", wait: 6, args: 1, fmt: izxFmt, addr: izxAddr, exec: sta}
	c.opCodes[0x84] = &opcode{mne: "STY", wait: 3, args: 1, fmt: zpFmt, addr: zpAddr, exec: sty}
	c.opCodes[0x85] = &opcode{mne: "STA", wait: 3, args: 1, fmt: zpFmt, addr: zpAddr, exec: sta}
	c.opCodes[0x86] = &opcode{mne: "STX", wait: 3, args: 1, fmt: zpFmt, addr: zpAddr, exec: stx}
	c.opCodes[0x88] = &opcode{mne: "DEY", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: dey}
	c.opCodes[0x8a] = &opcode{mne: "TXA", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: txa}
	c.opCodes[0x8c] = &opcode{mne: "STY", wait: 4, args: 2, fmt: absFmt, addr: absAddr, exec: sty}
	c.opCodes[0x8d] = &opcode{mne: "STA", wait: 4, args: 2, fmt: absFmt, addr: absAddr, exec: sta}
	c.opCodes[0x8e] = &opcode{mne: "STX", wait: 4, args: 2, fmt: absFmt, addr: absAddr, exec: stx}

	c.opCodes[0x90] = &opcode{mne: "BCC", wait: 2, args: 1, fmt: relFmt, addr: relAddr, exec: bcc}
	c.opCodes[0x91] = &opcode{mne: "STA", wait: 6, args: 1, fmt: izyFmt, addr: izyAddr, exec: sta}
	c.opCodes[0x94] = &opcode{mne: "STY", wait: 4, args: 1, fmt: zpxFmt, addr: zpxAddr, exec: sty}
	c.opCodes[0x95] = &opcode{mne: "STA", wait: 4, args: 1, fmt: zpxFmt, addr: zpxAddr, exec: sta}
	c.opCodes[0x96] = &opcode{mne: "STX", wait: 4, args: 1, fmt: zpyFmt, addr: zpyAddr, exec: stx}
	c.opCodes[0x98] = &opcode{mne: "TYA", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: tya}
	c.opCodes[0x99] = &opcode{mne: "STA", wait: 5, args: 2, fmt: abyFmt, addr: abyAddr, exec: sta}
	c.opCodes[0x9a] = &opcode{mne: "TXS", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: txs}
	c.opCodes[0x9d] = &opcode{mne: "STA", wait: 5, args: 2, fmt: abxFmt, addr: abxAddr, exec: sta}

	c.opCodes[0xa0] = &opcode{mne: "LDY", wait: 2, args: 1, fmt: immFmt, addr: immAddr, exec: ldy}
	c.opCodes[0xa1] = &opcode{mne: "LDX", wait: 6, args: 1, fmt: izxFmt, addr: izxAddr, exec: lda}
	c.opCodes[0xa2] = &opcode{mne: "LDX", wait: 2, args: 1, fmt: immFmt, addr: immAddr, exec: ldx}
	c.opCodes[0xa4] = &opcode{mne: "LDY", wait: 3, args: 1, fmt: zpFmt, addr: zpAddr, exec: ldy}
	c.opCodes[0xa5] = &opcode{mne: "LDA", wait: 3, args: 1, fmt: zpFmt, addr: zpAddr, exec: lda}
	c.opCodes[0xa6] = &opcode{mne: "LDX", wait: 3, args: 1, fmt: zpFmt, addr: zpAddr, exec: ldx}
	c.opCodes[0xa8] = &opcode{mne: "TAY", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: tay}
	c.opCodes[0xa9] = &opcode{mne: "LDA", wait: 2, args: 1, fmt: immFmt, addr: immAddr, exec: lda}
	c.opCodes[0xaa] = &opcode{mne: "TAX", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: tax}
	c.opCodes[0xac] = &opcode{mne: "LDY", wait: 4, args: 2, fmt: absFmt, addr: absAddr, exec: ldy}
	c.opCodes[0xad] = &opcode{mne: "LDA", wait: 4, args: 2, fmt: absFmt, addr: absAddr, exec: lda}
	c.opCodes[0xae] = &opcode{mne: "LDX", wait: 4, args: 2, fmt: absFmt, addr: absAddr, exec: ldx}

	c.opCodes[0xb0] = &opcode{mne: "BCS", wait: 2, args: 1, fmt: relFmt, addr: relAddr, exec: bcs}
	c.opCodes[0xb1] = &opcode{mne: "LDA", wait: 5, args: 1, fmt: izyFmt, addr: izyAddr, exec: lda}
	c.opCodes[0xb4] = &opcode{mne: "LDY", wait: 4, args: 1, fmt: zpxFmt, addr: zpxAddr, exec: ldy}
	c.opCodes[0xb5] = &opcode{mne: "LDA", wait: 4, args: 1, fmt: zpxFmt, addr: zpxAddr, exec: lda}
	c.opCodes[0xb6] = &opcode{mne: "LDX", wait: 4, args: 1, fmt: zpyFmt, addr: zpyAddr, exec: ldx}
	c.opCodes[0xb8] = &opcode{mne: "CLV", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: clv}
	c.opCodes[0xb9] = &opcode{mne: "LDA", wait: 4, args: 2, fmt: abyFmt, addr: abyAddr, exec: lda}
	c.opCodes[0xba] = &opcode{mne: "TSX", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: tsx}
	c.opCodes[0xbc] = &opcode{mne: "LDY", wait: 4, args: 2, fmt: abxFmt, addr: abxAddr, exec: ldy}
	c.opCodes[0xbd] = &opcode{mne: "LDA", wait: 4, args: 2, fmt: abxFmt, addr: abxAddr, exec: lda}
	c.opCodes[0xbe] = &opcode{mne: "LDX", wait: 2, args: 2, fmt: abyFmt, addr: abyAddr, exec: ldx}

	c.opCodes[0xc0] = &opcode{mne: "CPY", wait: 2, args: 1, fmt: immFmt, addr: immAddr, exec: cpy}
	c.opCodes[0xc1] = &opcode{mne: "CMP", wait: 6, args: 1, fmt: izxFmt, addr: izxAddr, exec: cmp}
	c.opCodes[0xc4] = &opcode{mne: "CPY", wait: 3, args: 1, fmt: zpFmt, addr: zpAddr, exec: cpy}
	c.opCodes[0xc5] = &opcode{mne: "CMP", wait: 3, args: 1, fmt: zpFmt, addr: zpAddr, exec: cmp}
	c.opCodes[0xc6] = &opcode{mne: "DEC", wait: 5, args: 1, fmt: zpFmt, addr: zpAddr, exec: dec}
	c.opCodes[0xc8] = &opcode{mne: "INY", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: iny}
	c.opCodes[0xc9] = &opcode{mne: "CMP", wait: 2, args: 1, fmt: immFmt, addr: immAddr, exec: cmp}
	c.opCodes[0xca] = &opcode{mne: "DEX", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: dex}
	c.opCodes[0xcc] = &opcode{mne: "CPY", wait: 4, args: 2, fmt: absFmt, addr: absAddr, exec: cpy}
	c.opCodes[0xcd] = &opcode{mne: "CMP", wait: 4, args: 2, fmt: absFmt, addr: absAddr, exec: cmp}
	c.opCodes[0xce] = &opcode{mne: "DEC", wait: 6, args: 2, fmt: absFmt, addr: absAddr, exec: dec}

	c.opCodes[0xd0] = &opcode{mne: "BNE", wait: 2, args: 1, fmt: relFmt, addr: relAddr, exec: bne}
	c.opCodes[0xd1] = &opcode{mne: "CMP", wait: 5, args: 1, fmt: izyFmt, addr: izyAddr, exec: cmp}
	c.opCodes[0xd5] = &opcode{mne: "CMP", wait: 2, args: 1, fmt: zpxFmt, addr: zpxAddr, exec: cmp}
	c.opCodes[0xd6] = &opcode{mne: "DEC", wait: 6, args: 1, fmt: zpxFmt, addr: zpxAddr, exec: dec}
	c.opCodes[0xd8] = &opcode{mne: "CLD", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: cld}
	c.opCodes[0xd9] = &opcode{mne: "CMP", wait: 4, args: 2, fmt: abyFmt, addr: abyAddr, exec: cmp}
	c.opCodes[0xda] = &opcode{mne: "NOP*", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: nop}
	c.opCodes[0xdd] = &opcode{mne: "CMP", wait: 4, args: 2, fmt: abxFmt, addr: abxAddr, exec: cmp}
	c.opCodes[0xde] = &opcode{mne: "DEC", wait: 7, args: 2, fmt: abxFmt, addr: abxAddr, exec: dec}

	c.opCodes[0xe0] = &opcode{mne: "CPX", wait: 2, args: 1, fmt: immFmt, addr: immAddr, exec: cpx}
	c.opCodes[0xe1] = &opcode{mne: "SBC", wait: 6, args: 1, fmt: izxFmt, addr: izxAddr, exec: sbc}
	c.opCodes[0xe4] = &opcode{mne: "CPX", wait: 3, args: 1, fmt: zpFmt, addr: zpAddr, exec: cpx}
	c.opCodes[0xe5] = &opcode{mne: "SBC", wait: 3, args: 1, fmt: zpFmt, addr: zpAddr, exec: sbc}
	c.opCodes[0xe6] = &opcode{mne: "INC", wait: 5, args: 1, fmt: zpFmt, addr: zpAddr, exec: inc}
	c.opCodes[0xe8] = &opcode{mne: "INX", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: inx}
	c.opCodes[0xe9] = &opcode{mne: "SBC", wait: 2, args: 1, fmt: immFmt, addr: immAddr, exec: sbc}
	c.opCodes[0xea] = &opcode{mne: "NOP", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: nop}
	c.opCodes[0xec] = &opcode{mne: "CPX", wait: 4, args: 2, fmt: absFmt, addr: absAddr, exec: cpx}
	c.opCodes[0xed] = &opcode{mne: "SBC", wait: 4, args: 2, fmt: absFmt, addr: absAddr, exec: sbc}
	c.opCodes[0xee] = &opcode{mne: "INC", wait: 6, args: 2, fmt: absFmt, addr: absAddr, exec: inc}

	c.opCodes[0xf0] = &opcode{mne: "BEQ", wait: 2, args: 1, fmt: relFmt, addr: relAddr, exec: beq}
	c.opCodes[0xf1] = &opcode{mne: "SBC", wait: 5, args: 1, fmt: izyFmt, addr: izyAddr, exec: sbc}
	c.opCodes[0xf5] = &opcode{mne: "SBC", wait: 4, args: 1, fmt: zpxFmt, addr: zpxAddr, exec: sbc}
	c.opCodes[0xf6] = &opcode{mne: "INC", wait: 6, args: 1, fmt: zpxFmt, addr: zpxAddr, exec: inc}
	c.opCodes[0xf8] = &opcode{mne: "SED", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: sed}
	c.opCodes[0xf9] = &opcode{mne: "SBC", wait: 4, args: 2, fmt: abyFmt, addr: abyAddr, exec: sbc}
	c.opCodes[0xfa] = &opcode{mne: "NOP*", wait: 2, args: 0, fmt: nilFmt, addr: nilAddr, exec: nop}
	c.opCodes[0xfd] = &opcode{mne: "SBC", wait: 4, args: 2, fmt: abxFmt, addr: abxAddr, exec: sbc}
	c.opCodes[0xfe] = &opcode{mne: "INC", wait: 7, args: 2, fmt: abxFmt, addr: abxAddr, exec: inc}

	// c.opCodes[0x00] = func() { c.waitCycles = 7; c.brk() }                            // BRK
	// c.opCodes[0x01] = func() { c.waitCycles = 6; c.ora(c.indexedIndirectAddress()) }  // ORA ($00, X)
	// c.opCodes[0x05] = func() { c.waitCycles = 3; c.ora(c.zeroPageAddress()) }         // ORA $00
	// c.opCodes[0x06] = func() { c.waitCycles = 5; c.asl(c.zeroPageAddress()) }         // ASL $00
	// c.opCodes[0x08] = func() { c.waitCycles = 3; c.php() }                            // PHP
	// c.opCodes[0x09] = func() { c.waitCycles = 2; c.ora(c.immediateAddress()) }        // ORA #$00
	// c.opCodes[0x0a] = func() { c.waitCycles = 2; c.aslAcc() }                         // ASL A
	// c.opCodes[0x0d] = func() { c.waitCycles = 4; c.ora(c.absoluteAddress()) }         // ORA $0000
	// c.opCodes[0x0e] = func() { c.waitCycles = 6; c.asl(c.absoluteAddress()) }         // ASL $0000
	// c.opCodes[0x10] = func() { c.waitCycles = 2; c.bpl(c.relativeAddress()) }         // BPL $0000 (rel)
	// c.opCodes[0x11] = func() { c.waitCycles = 5; c.ora(c.indirectIndexedAddress()) }  // ORA ($00), Y
	// c.opCodes[0x15] = func() { c.waitCycles = 4; c.ora(c.zeroPageXIndexedAddress()) } // ORA $00, X
	// c.opCodes[0x16] = func() { c.waitCycles = 6; c.asl(c.zeroPageXIndexedAddress()) } // ASL $00, X
	// c.opCodes[0x18] = func() { c.waitCycles = 2; c.clc() }                            // CLC
	// c.opCodes[0x19] = func() { c.waitCycles = 4; c.ora(c.absoluteIndexedYAddress()) } // ORA $0000, Y
	// c.opCodes[0x1a] = func() { c.waitCycles = 2 }                                     // NOP
	// c.opCodes[0x1d] = func() { c.waitCycles = 4; c.ora(c.absoluteIndexedXAddress()) } // ORA $0000, X
	// c.opCodes[0x1e] = func() { c.waitCycles = 7; c.asl(c.absoluteIndexedXAddress()) } // ASL $0000, X
	// c.opCodes[0x20] = func() { c.waitCycles = 6; c.jsr(c.absoluteAddress()) }         // JSR $0000
	// c.opCodes[0x21] = func() { c.waitCycles = 6; c.and(c.indexedIndirectAddress()) }  // AND ($00, X)
	// c.opCodes[0x24] = func() { c.waitCycles = 3; c.bit(c.zeroPageAddress()) }         // BIT $00
	// c.opCodes[0x25] = func() { c.waitCycles = 3; c.and(c.zeroPageAddress()) }         // AND $00
	// c.opCodes[0x26] = func() { c.waitCycles = 5; c.rol(c.zeroPageAddress()) }         // ROL $00
	// c.opCodes[0x28] = func() { c.waitCycles = 4; c.plp() }                            // PLP
	// c.opCodes[0x29] = func() { c.waitCycles = 2; c.and(c.immediateAddress()) }        // AND #$00
	// c.opCodes[0x2a] = func() { c.waitCycles = 2; c.rolAcc() }                         // ROL A
	// c.opCodes[0x2c] = func() { c.waitCycles = 4; c.bit(c.absoluteAddress()) }         // BIT $0000
	// c.opCodes[0x2d] = func() { c.waitCycles = 4; c.and(c.absoluteAddress()) }         // AND $0000
	// c.opCodes[0x2e] = func() { c.waitCycles = 6; c.rol(c.absoluteAddress()) }         // ROL $0000
	// c.opCodes[0x30] = func() { c.waitCycles = 2; c.bmi(c.relativeAddress()) }         // BMI $0000 (rel)
	// c.opCodes[0x31] = func() { c.waitCycles = 5; c.and(c.indirectIndexedAddress()) }  // AND ($00), Y
	// c.opCodes[0x35] = func() { c.waitCycles = 4; c.and(c.zeroPageXIndexedAddress()) } // AND $00, X
	// c.opCodes[0x36] = func() { c.waitCycles = 6; c.rol(c.zeroPageXIndexedAddress()) } // ROL $00, X
	// c.opCodes[0x38] = func() { c.waitCycles = 2; c.sec() }                            // SEC
	// c.opCodes[0x39] = func() { c.waitCycles = 4; c.and(c.absoluteIndexedYAddress()) } // AND $0000, Y
	// c.opCodes[0x3a] = func() { c.waitCycles = 2 }                                     // NOP
	// c.opCodes[0x3d] = func() { c.waitCycles = 4; c.and(c.absoluteIndexedXAddress()) } // AND $0000, X
	// c.opCodes[0x3e] = func() { c.waitCycles = 7; c.rol(c.absoluteIndexedXAddress()) } // ROL $0000, X
	// c.opCodes[0x40] = func() { c.waitCycles = 6; c.rti() }                            // RTI
	// c.opCodes[0x41] = func() { c.waitCycles = 6; c.eor(c.indexedIndirectAddress()) }  // EOR ($00, X)
	// c.opCodes[0x45] = func() { c.waitCycles = 3; c.eor(c.zeroPageAddress()) }         // EOR $00
	// c.opCodes[0x46] = func() { c.waitCycles = 5; c.lsr(c.zeroPageAddress()) }         // LSR $00
	// c.opCodes[0x48] = func() { c.waitCycles = 3; c.pha() }                            // PHA
	// c.opCodes[0x49] = func() { c.waitCycles = 2; c.eor(c.immediateAddress()) }        // EOR #$00
	// c.opCodes[0x4a] = func() { c.waitCycles = 2; c.lsrAcc() }                         // LSR A
	// c.opCodes[0x4c] = func() { c.waitCycles = 3; c.jmp(c.absoluteAddress()) }         // JMP $0000
	// c.opCodes[0x4d] = func() { c.waitCycles = 4; c.eor(c.absoluteAddress()) }         // EOR $0000
	// c.opCodes[0x4e] = func() { c.waitCycles = 6; c.lsr(c.absoluteAddress()) }         // LSR $0000
	// c.opCodes[0x50] = func() { c.waitCycles = 2; c.bvc(c.relativeAddress()) }         // BVC $0000 (rel)
	// c.opCodes[0x51] = func() { c.waitCycles = 5; c.eor(c.indirectIndexedAddress()) }  // EOR ($00), Y
	// c.opCodes[0x55] = func() { c.waitCycles = 4; c.eor(c.zeroPageXIndexedAddress()) } // EOR $00, X
	// c.opCodes[0x56] = func() { c.waitCycles = 6; c.lsr(c.zeroPageXIndexedAddress()) } // LSR $00, X
	// c.opCodes[0x58] = func() { c.waitCycles = 2; c.cli() }                            // CLI
	// c.opCodes[0x59] = func() { c.waitCycles = 4; c.eor(c.absoluteIndexedYAddress()) } // EOR $0000, Y
	// c.opCodes[0x5a] = func() { c.waitCycles = 2 }                                     // NOP
	// c.opCodes[0x5d] = func() { c.waitCycles = 4; c.eor(c.absoluteIndexedXAddress()) } // EOR $0000, X
	// c.opCodes[0x5e] = func() { c.waitCycles = 7; c.lsr(c.absoluteIndexedXAddress()) } // LSR $0000, X
	// c.opCodes[0x60] = func() { c.waitCycles = 6; c.rts() }                            // RTS
	// c.opCodes[0x61] = func() { c.waitCycles = 6; c.adc(c.indexedIndirectAddress()) }  // ADC ($00, X)
	// c.opCodes[0x65] = func() { c.waitCycles = 3; c.adc(c.zeroPageAddress()) }         // ADC $00
	// c.opCodes[0x66] = func() { c.waitCycles = 5; c.ror(c.zeroPageAddress()) }         // ROR $00
	// c.opCodes[0x68] = func() { c.waitCycles = 4; c.pla() }                            // PLA
	// c.opCodes[0x69] = func() { c.waitCycles = 2; c.adc(c.immediateAddress()) }        // ADC #$00
	// c.opCodes[0x6a] = func() { c.waitCycles = 2; c.rorAcc() }                         // ROR A
	// c.opCodes[0x6c] = func() { c.waitCycles = 5; c.jmp(c.indirectAbsoluteAddress()) } // JMP ($0000)
	// c.opCodes[0x6d] = func() { c.waitCycles = 4; c.adc(c.absoluteAddress()) }         // ADC $0000
	// c.opCodes[0x6e] = func() { c.waitCycles = 6; c.ror(c.absoluteAddress()) }         // ROR $0000
	// c.opCodes[0x70] = func() { c.waitCycles = 2; c.bvs(c.relativeAddress()) }         // BVS $0000 (rel)
	// c.opCodes[0x71] = func() { c.waitCycles = 5; c.adc(c.indirectIndexedAddress()) }  // ADC ($00), Y
	// c.opCodes[0x75] = func() { c.waitCycles = 4; c.adc(c.zeroPageXIndexedAddress()) } // ADC $00, X
	// c.opCodes[0x76] = func() { c.waitCycles = 6; c.ror(c.zeroPageXIndexedAddress()) } // ROR $00, X
	// c.opCodes[0x78] = func() { c.waitCycles = 2; c.sei() }                            // SEI
	// c.opCodes[0x79] = func() { c.waitCycles = 4; c.adc(c.absoluteIndexedYAddress()) } // ADC $0000, Y
	// c.opCodes[0x7a] = func() { c.waitCycles = 2 }                                     // NOP
	// c.opCodes[0x7d] = func() { c.waitCycles = 4; c.adc(c.absoluteIndexedXAddress()) } // ADC $0000, X
	// c.opCodes[0x7e] = func() { c.waitCycles = 7; c.ror(c.absoluteIndexedXAddress()) } // ROR $0000, X
	// c.opCodes[0x81] = func() { c.waitCycles = 6; c.sta(c.indexedIndirectAddress()) }  // STA ($00, X)
	// c.opCodes[0x84] = func() { c.waitCycles = 3; c.sty(c.zeroPageAddress()) }         // STY $00
	// c.opCodes[0x85] = func() { c.waitCycles = 3; c.sta(c.zeroPageAddress()) }         // STA $00
	// c.opCodes[0x86] = func() { c.waitCycles = 3; c.stx(c.zeroPageAddress()) }         // STX $00
	// c.opCodes[0x88] = func() { c.waitCycles = 2; c.dey() }                            // DEY
	// c.opCodes[0x8a] = func() { c.waitCycles = 2; c.txa() }                            // TXA
	// c.opCodes[0x8c] = func() { c.waitCycles = 4; c.sty(c.absoluteAddress()) }         // STY $0000
	// c.opCodes[0x8d] = func() { c.waitCycles = 4; c.sta(c.absoluteAddress()) }         // STA $0000
	// c.opCodes[0x8e] = func() { c.waitCycles = 4; c.stx(c.absoluteAddress()) }         // STX $0000
	// c.opCodes[0x90] = func() { c.waitCycles = 2; c.bcc(c.relativeAddress()) }         // BCC $0000 (rel)
	// c.opCodes[0x91] = func() { c.waitCycles = 6; c.sta(c.indirectIndexedAddress()) }  // STA ($00), Y
	// c.opCodes[0x94] = func() { c.waitCycles = 3; c.sty(c.zeroPageXIndexedAddress()) } // STY $00, X
	// c.opCodes[0x95] = func() { c.waitCycles = 4; c.sta(c.zeroPageXIndexedAddress()) } // STA $00, X
	// c.opCodes[0x96] = func() { c.waitCycles = 4; c.stx(c.zeroPageYIndexedAddress()) } // STX $00, Y
	// c.opCodes[0x98] = func() { c.waitCycles = 2; c.tya() }                            // TYA
	// c.opCodes[0x99] = func() { c.waitCycles = 5; c.sta(c.absoluteIndexedYAddress()) } // STA $0000, Y
	// c.opCodes[0x9a] = func() { c.waitCycles = 2; c.txs() }                            // TXS
	// c.opCodes[0x9d] = func() { c.waitCycles = 5; c.sta(c.absoluteIndexedXAddress()) } // STA $0000, X
	// c.opCodes[0xa0] = func() { c.waitCycles = 2; c.ldy(c.immediateAddress()) }        // LDY #$00
	// c.opCodes[0xa1] = func() { c.waitCycles = 6; c.lda(c.indexedIndirectAddress()) }  // LDX ($00, X)
	// c.opCodes[0xa2] = func() { c.waitCycles = 2; c.ldx(c.immediateAddress()) }        // LDX #$00
	// c.opCodes[0xa4] = func() { c.waitCycles = 3; c.ldy(c.zeroPageAddress()) }         // LDY $00
	// c.opCodes[0xa5] = func() { c.waitCycles = 3; c.lda(c.zeroPageAddress()) }         // LDA $00
	// c.opCodes[0xa6] = func() { c.waitCycles = 3; c.ldx(c.zeroPageAddress()) }         // LDX $00
	// c.opCodes[0xa8] = func() { c.waitCycles = 2; c.tay() }                            // TAY
	// c.opCodes[0xa9] = func() { c.waitCycles = 2; c.lda(c.immediateAddress()) }        // LDA #$00
	// c.opCodes[0xaa] = func() { c.waitCycles = 2; c.tax() }                            // TAX
	// c.opCodes[0xac] = func() { c.waitCycles = 4; c.ldy(c.absoluteAddress()) }         // LDY $0000
	// c.opCodes[0xad] = func() { c.waitCycles = 4; c.lda(c.absoluteAddress()) }         // LDA $0000
	// c.opCodes[0xae] = func() { c.waitCycles = 4; c.ldx(c.absoluteAddress()) }         // LDX $0000
	// c.opCodes[0xb0] = func() { c.waitCycles = 2; c.bcs(c.relativeAddress()) }         // BCS $0000 (rel)
	// c.opCodes[0xb1] = func() { c.waitCycles = 5; c.lda(c.indirectIndexedAddress()) }  // LDA ($00), Y
	// c.opCodes[0xb4] = func() { c.waitCycles = 4; c.ldy(c.zeroPageXIndexedAddress()) } // LDY $00, X
	// c.opCodes[0xb5] = func() { c.waitCycles = 4; c.lda(c.zeroPageXIndexedAddress()) } // LDA $00, X
	// c.opCodes[0xb6] = func() { c.waitCycles = 4; c.ldx(c.zeroPageYIndexedAddress()) } // LDX $00, Y
	// c.opCodes[0xb8] = func() { c.waitCycles = 2; c.clv() }                            // CLV
	// c.opCodes[0xb9] = func() { c.waitCycles = 4; c.lda(c.absoluteIndexedYAddress()) } // LDA $0000, Y
	// c.opCodes[0xba] = func() { c.waitCycles = 2; c.tsx() }                            // TSX
	// c.opCodes[0xbc] = func() { c.waitCycles = 4; c.ldy(c.absoluteIndexedXAddress()) } // LDY $0000, X
	// c.opCodes[0xbd] = func() { c.waitCycles = 4; c.lda(c.absoluteIndexedXAddress()) } // LDA $0000, X
	// c.opCodes[0xbe] = func() { c.waitCycles = 2; c.ldx(c.absoluteIndexedYAddress()) } // LDX $0000, Y
	// c.opCodes[0xc0] = func() { c.waitCycles = 2; c.cpy(c.immediateAddress()) }        // CPY #$00
	// c.opCodes[0xc1] = func() { c.waitCycles = 6; c.cmp(c.indexedIndirectAddress()) }  // CMP ($00, X)
	// c.opCodes[0xc4] = func() { c.waitCycles = 3; c.cpy(c.zeroPageAddress()) }         // CPY $00
	// c.opCodes[0xc5] = func() { c.waitCycles = 3; c.cmp(c.zeroPageAddress()) }         // CMP $00
	// c.opCodes[0xc6] = func() { c.waitCycles = 5; c.dec(c.zeroPageAddress()) }         // DEC $00
	// c.opCodes[0xc8] = func() { c.waitCycles = 2; c.iny() }                            // INY
	// c.opCodes[0xc9] = func() { c.waitCycles = 2; c.cmp(c.immediateAddress()) }        // CMP #$00
	// c.opCodes[0xca] = func() { c.waitCycles = 2; c.dex() }                            // DEX
	// c.opCodes[0xcc] = func() { c.waitCycles = 4; c.cpy(c.absoluteAddress()) }         // CPY $0000
	// c.opCodes[0xcd] = func() { c.waitCycles = 4; c.cmp(c.absoluteAddress()) }         // CMP $0000
	// c.opCodes[0xd0] = func() { c.waitCycles = 3; c.bne(c.relativeAddress()) }         // BNE $0000 (rel)
	// c.opCodes[0xc9] = func() { c.waitCycles = 2; c.cmp(c.immediateAddress()) }        // CMP #$00
	// c.opCodes[0xce] = func() { c.waitCycles = 6; c.dec(c.absoluteAddress()) }         // DEC $0000
	// c.opCodes[0xd1] = func() { c.waitCycles = 5; c.cmp(c.indirectIndexedAddress()) }  // CMP ($00), Y
	// c.opCodes[0xd5] = func() { c.waitCycles = 2; c.cmp(c.zeroPageXIndexedAddress()) } // CMP $00, X
	// c.opCodes[0xd6] = func() { c.waitCycles = 6; c.dec(c.zeroPageXIndexedAddress()) } // DEC $00, X
	// c.opCodes[0xd8] = func() { c.waitCycles = 2; c.cld() }                            // CLD
	// c.opCodes[0xd9] = func() { c.waitCycles = 4; c.cmp(c.absoluteIndexedYAddress()) } // CMP $0000, Y
	// c.opCodes[0xda] = func() { c.waitCycles = 2 }                                     // NOP
	// c.opCodes[0xdd] = func() { c.waitCycles = 4; c.cmp(c.absoluteIndexedXAddress()) } // CMP $0000, X
	// c.opCodes[0xde] = func() { c.waitCycles = 7; c.dec(c.absoluteIndexedXAddress()) } // DEC $0000, X
	// c.opCodes[0xe0] = func() { c.waitCycles = 2; c.cpx(c.immediateAddress()) }        // CPX #$00
	// c.opCodes[0xe1] = func() { c.waitCycles = 6; c.sbc(c.indexedIndirectAddress()) }  // SBC ($00, X)
	// c.opCodes[0xe4] = func() { c.waitCycles = 3; c.cpx(c.zeroPageAddress()) }         // CPX $00
	// c.opCodes[0xe5] = func() { c.waitCycles = 3; c.sbc(c.zeroPageAddress()) }         // SBC $00
	// c.opCodes[0xe6] = func() { c.waitCycles = 5; c.inc(c.zeroPageAddress()) }         // INC $00
	// c.opCodes[0xe8] = func() { c.waitCycles = 2; c.inx() }                            // INX
	// c.opCodes[0xe9] = func() { c.waitCycles = 2; c.sbc(c.immediateAddress()) }        // SBC #$00
	// c.opCodes[0xea] = func() { c.waitCycles = 2 }                                     // NOP
	// c.opCodes[0xec] = func() { c.waitCycles = 4; c.cpx(c.absoluteAddress()) }         // CPX $0000
	// c.opCodes[0xed] = func() { c.waitCycles = 4; c.sbc(c.absoluteAddress()) }         // SBC $0000
	// c.opCodes[0xee] = func() { c.waitCycles = 6; c.inc(c.absoluteAddress()) }         // INC $0000
	// c.opCodes[0xf0] = func() { c.waitCycles = 2; c.beq(c.relativeAddress()) }         // BEQ (rel)
	// c.opCodes[0xf1] = func() { c.waitCycles = 5; c.sbc(c.indirectIndexedAddress()) }  // SBC ($00), Y
	// c.opCodes[0xf5] = func() { c.waitCycles = 4; c.sbc(c.zeroPageXIndexedAddress()) } // SBC $00, X
	// c.opCodes[0xf6] = func() { c.waitCycles = 6; c.inc(c.zeroPageXIndexedAddress()) } // INC $00, X
	// c.opCodes[0xf8] = func() { c.waitCycles = 2; c.sed() }                            // SED
	// c.opCodes[0xf9] = func() { c.waitCycles = 4; c.sbc(c.absoluteIndexedYAddress()) } // SBC $0000, Y
	// c.opCodes[0xfa] = func() { c.waitCycles = 2 }                                     // NOP
	// c.opCodes[0xfd] = func() { c.waitCycles = 4; c.sbc(c.absoluteIndexedXAddress()) } // SBC $0000, X
	// c.opCodes[0xfe] = func() { c.waitCycles = 7; c.inc(c.absoluteIndexedXAddress()) } // INC $0000, X
}
