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

type opcode struct {
	mne      string
	wait     int
	addrMode addrMode
	exec     func(uint16)
}

func (c *Cpu6502) initOpcodes() {
	c.opCodes[0x00] = &opcode{mne: "BRK", wait: 7, addrMode: none, exec: c.brk}
	c.opCodes[0x01] = &opcode{mne: "ORA", wait: 6, addrMode: izx, exec: c.ora}
	c.opCodes[0x03] = &opcode{mne: "SLO*", wait: 8, addrMode: izx, exec: c.slo}
	c.opCodes[0x04] = &opcode{mne: "NOP*", wait: 3, addrMode: zp, exec: c.nop}
	c.opCodes[0x05] = &opcode{mne: "ORA", wait: 3, addrMode: zp, exec: c.ora}
	c.opCodes[0x06] = &opcode{mne: "ASL", wait: 5, addrMode: zp, exec: c.asl}
	c.opCodes[0x07] = &opcode{mne: "SLO*", wait: 5, addrMode: zp, exec: c.slo}
	c.opCodes[0x08] = &opcode{mne: "PHP", wait: 3, addrMode: none, exec: c.php}
	c.opCodes[0x09] = &opcode{mne: "ORA", wait: 2, addrMode: imm, exec: c.ora}
	c.opCodes[0x0a] = &opcode{mne: "ASL", wait: 2, addrMode: none, exec: c.aslAcc}
	c.opCodes[0x0b] = &opcode{mne: "ANC*", wait: 2, addrMode: imm, exec: c.anc}
	c.opCodes[0x0c] = &opcode{mne: "NOP*", wait: 4, addrMode: abs, exec: c.nop}
	c.opCodes[0x0d] = &opcode{mne: "ORA", wait: 4, addrMode: abs, exec: c.ora}
	c.opCodes[0x0e] = &opcode{mne: "ASL", wait: 6, addrMode: abs, exec: c.asl}
	c.opCodes[0x0f] = &opcode{mne: "SLO*", wait: 6, addrMode: abs, exec: c.slo}

	c.opCodes[0x10] = &opcode{mne: "BPL", wait: 2, addrMode: rel, exec: c.bpl}
	c.opCodes[0x11] = &opcode{mne: "ORA", wait: 5, addrMode: izy, exec: c.ora}
	c.opCodes[0x13] = &opcode{mne: "SLO*", wait: 8, addrMode: izy, exec: c.slo}
	c.opCodes[0x14] = &opcode{mne: "NOP*", wait: 3, addrMode: zpx, exec: c.nop}
	c.opCodes[0x15] = &opcode{mne: "ORA", wait: 4, addrMode: zpx, exec: c.ora}
	c.opCodes[0x16] = &opcode{mne: "ASL", wait: 6, addrMode: zpx, exec: c.asl}
	c.opCodes[0x17] = &opcode{mne: "SLO*", wait: 6, addrMode: zpx, exec: c.slo}
	c.opCodes[0x18] = &opcode{mne: "CLC", wait: 2, addrMode: none, exec: c.clc}
	c.opCodes[0x19] = &opcode{mne: "ORA", wait: 4, addrMode: aby, exec: c.ora}
	c.opCodes[0x1a] = &opcode{mne: "NOP", wait: 2, addrMode: none, exec: c.nop}
	c.opCodes[0x1b] = &opcode{mne: "SLO*", wait: 7, addrMode: aby, exec: c.slo}
	c.opCodes[0x1c] = &opcode{mne: "NOP*", wait: 4, addrMode: abx, exec: c.nop}
	c.opCodes[0x1d] = &opcode{mne: "ORA", wait: 4, addrMode: abx, exec: c.ora}
	c.opCodes[0x1e] = &opcode{mne: "ASL", wait: 7, addrMode: abx, exec: c.asl}
	c.opCodes[0x1f] = &opcode{mne: "SLO*", wait: 7, addrMode: abx, exec: c.slo}

	c.opCodes[0x20] = &opcode{mne: "JSR", wait: 6, addrMode: abs, exec: c.jsr}
	c.opCodes[0x21] = &opcode{mne: "AND", wait: 6, addrMode: izx, exec: c.and}
	c.opCodes[0x23] = &opcode{mne: "RLA*", wait: 8, addrMode: izx, exec: c.rla}
	c.opCodes[0x24] = &opcode{mne: "BIT", wait: 3, addrMode: zp, exec: c.bit}
	c.opCodes[0x25] = &opcode{mne: "AND", wait: 3, addrMode: zp, exec: c.and}
	c.opCodes[0x26] = &opcode{mne: "ROL", wait: 5, addrMode: zp, exec: c.rol}
	c.opCodes[0x27] = &opcode{mne: "RLA*", wait: 5, addrMode: zp, exec: c.rla}
	c.opCodes[0x28] = &opcode{mne: "PLP", wait: 4, addrMode: none, exec: c.plp}
	c.opCodes[0x29] = &opcode{mne: "AND", wait: 2, addrMode: imm, exec: c.and}
	c.opCodes[0x2a] = &opcode{mne: "ROL", wait: 2, addrMode: none, exec: c.rolAcc}
	c.opCodes[0x2b] = &opcode{mne: "ANC*", wait: 2, addrMode: imm, exec: c.anc}
	c.opCodes[0x2c] = &opcode{mne: "BIT", wait: 4, addrMode: abs, exec: c.bit}
	c.opCodes[0x2d] = &opcode{mne: "AND", wait: 4, addrMode: abs, exec: c.and}
	c.opCodes[0x2e] = &opcode{mne: "ROL", wait: 6, addrMode: abs, exec: c.rol}
	c.opCodes[0x2f] = &opcode{mne: "RLA*", wait: 6, addrMode: abs, exec: c.rla}

	c.opCodes[0x30] = &opcode{mne: "BMI", wait: 2, addrMode: rel, exec: c.bmi}
	c.opCodes[0x31] = &opcode{mne: "AND", wait: 5, addrMode: izy, exec: c.and}
	c.opCodes[0x33] = &opcode{mne: "RLA*", wait: 8, addrMode: izy, exec: c.rla}
	c.opCodes[0x34] = &opcode{mne: "NOP*", wait: 3, addrMode: zpx, exec: c.nop}
	c.opCodes[0x35] = &opcode{mne: "AND", wait: 4, addrMode: zpx, exec: c.and}
	c.opCodes[0x36] = &opcode{mne: "ROL", wait: 6, addrMode: zpx, exec: c.rol}
	c.opCodes[0x37] = &opcode{mne: "RLA*", wait: 6, addrMode: zpx, exec: c.rla}
	c.opCodes[0x38] = &opcode{mne: "SEC", wait: 2, addrMode: none, exec: c.sec}
	c.opCodes[0x39] = &opcode{mne: "AND", wait: 4, addrMode: aby, exec: c.and}
	c.opCodes[0x3a] = &opcode{mne: "NOP*", wait: 2, addrMode: none, exec: c.nop}
	c.opCodes[0x3b] = &opcode{mne: "RLA*", wait: 7, addrMode: aby, exec: c.rla}
	c.opCodes[0x3c] = &opcode{mne: "NOP*", wait: 4, addrMode: abx, exec: c.nop}
	c.opCodes[0x3d] = &opcode{mne: "AND", wait: 4, addrMode: abx, exec: c.and}
	c.opCodes[0x3e] = &opcode{mne: "ROL", wait: 7, addrMode: abx, exec: c.rol}
	c.opCodes[0x3f] = &opcode{mne: "RLA*", wait: 7, addrMode: abx, exec: c.rla}

	c.opCodes[0x40] = &opcode{mne: "RTI", wait: 6, addrMode: none, exec: c.rti}
	c.opCodes[0x41] = &opcode{mne: "EOR", wait: 6, addrMode: izx, exec: c.eor}
	c.opCodes[0x43] = &opcode{mne: "SRE*", wait: 8, addrMode: izx, exec: c.sre}
	c.opCodes[0x44] = &opcode{mne: "NOP*", wait: 3, addrMode: zp, exec: c.nop}
	c.opCodes[0x45] = &opcode{mne: "EOR", wait: 3, addrMode: zp, exec: c.eor}
	c.opCodes[0x46] = &opcode{mne: "LSR", wait: 5, addrMode: zp, exec: c.lsr}
	c.opCodes[0x47] = &opcode{mne: "SRE*", wait: 5, addrMode: zp, exec: c.sre}
	c.opCodes[0x48] = &opcode{mne: "PHA", wait: 3, addrMode: none, exec: c.pha}
	c.opCodes[0x49] = &opcode{mne: "EOR", wait: 2, addrMode: imm, exec: c.eor}
	c.opCodes[0x4a] = &opcode{mne: "LSR", wait: 2, addrMode: none, exec: c.lsrAcc}
	c.opCodes[0x4b] = &opcode{mne: "ALR*", wait: 2, addrMode: imm, exec: c.alr}
	c.opCodes[0x4c] = &opcode{mne: "JMP", wait: 3, addrMode: abs, exec: c.jmp}
	c.opCodes[0x4d] = &opcode{mne: "EOR", wait: 4, addrMode: abs, exec: c.eor}
	c.opCodes[0x4e] = &opcode{mne: "LSR", wait: 6, addrMode: abs, exec: c.lsr}
	c.opCodes[0x4f] = &opcode{mne: "SRE*", wait: 6, addrMode: abs, exec: c.sre}

	c.opCodes[0x50] = &opcode{mne: "BVC", wait: 2, addrMode: rel, exec: c.bvc}
	c.opCodes[0x51] = &opcode{mne: "EOR", wait: 5, addrMode: izy, exec: c.eor}
	c.opCodes[0x53] = &opcode{mne: "SRE*", wait: 8, addrMode: izy, exec: c.sre}
	c.opCodes[0x54] = &opcode{mne: "NOP*", wait: 3, addrMode: zpx, exec: c.nop}
	c.opCodes[0x55] = &opcode{mne: "EOR", wait: 4, addrMode: zpx, exec: c.eor}
	c.opCodes[0x56] = &opcode{mne: "LSR", wait: 6, addrMode: zpx, exec: c.lsr}
	c.opCodes[0x57] = &opcode{mne: "SRE*", wait: 6, addrMode: zpx, exec: c.sre}
	c.opCodes[0x58] = &opcode{mne: "CLI", wait: 2, addrMode: none, exec: c.cli}
	c.opCodes[0x59] = &opcode{mne: "EOR", wait: 4, addrMode: aby, exec: c.eor}
	c.opCodes[0x5a] = &opcode{mne: "NOP*", wait: 2, addrMode: none, exec: c.nop}
	c.opCodes[0x5b] = &opcode{mne: "SRE*", wait: 7, addrMode: aby, exec: c.sre}
	c.opCodes[0x5c] = &opcode{mne: "NOP*", wait: 4, addrMode: abx, exec: c.nop}
	c.opCodes[0x5d] = &opcode{mne: "EOR", wait: 4, addrMode: abx, exec: c.eor}
	c.opCodes[0x5e] = &opcode{mne: "LSR", wait: 7, addrMode: abx, exec: c.lsr}
	c.opCodes[0x5f] = &opcode{mne: "SRE*", wait: 7, addrMode: abx, exec: c.sre}

	c.opCodes[0x60] = &opcode{mne: "RTS", wait: 6, addrMode: none, exec: c.rts}
	c.opCodes[0x61] = &opcode{mne: "ADC", wait: 6, addrMode: izx, exec: c.adc}
	c.opCodes[0x63] = &opcode{mne: "RRA*", wait: 8, addrMode: izx, exec: c.rra}
	c.opCodes[0x64] = &opcode{mne: "NOP*", wait: 3, addrMode: zp, exec: c.nop}
	c.opCodes[0x65] = &opcode{mne: "ADC", wait: 3, addrMode: zp, exec: c.adc}
	c.opCodes[0x66] = &opcode{mne: "ROR", wait: 5, addrMode: zp, exec: c.ror}
	c.opCodes[0x67] = &opcode{mne: "RRA*", wait: 5, addrMode: zp, exec: c.rra}
	c.opCodes[0x68] = &opcode{mne: "PLA", wait: 4, addrMode: none, exec: c.pla}
	c.opCodes[0x69] = &opcode{mne: "ADC", wait: 2, addrMode: imm, exec: c.adc}
	c.opCodes[0x6a] = &opcode{mne: "ROR", wait: 2, addrMode: none, exec: c.rorAcc}
	c.opCodes[0x6b] = &opcode{mne: "ARR*", wait: 2, addrMode: imm, exec: c.arr}
	c.opCodes[0x6c] = &opcode{mne: "JMP", wait: 5, addrMode: ind, exec: c.jmp}
	c.opCodes[0x6d] = &opcode{mne: "ADC", wait: 4, addrMode: abs, exec: c.adc}
	c.opCodes[0x6e] = &opcode{mne: "ROR", wait: 6, addrMode: abs, exec: c.ror}
	c.opCodes[0x6f] = &opcode{mne: "RRA*", wait: 6, addrMode: abs, exec: c.rra}

	c.opCodes[0x70] = &opcode{mne: "BVS", wait: 2, addrMode: rel, exec: c.bvs}
	c.opCodes[0x71] = &opcode{mne: "ADC", wait: 5, addrMode: izy, exec: c.adc}
	c.opCodes[0x73] = &opcode{mne: "RRA*", wait: 8, addrMode: izy, exec: c.rra}
	c.opCodes[0x74] = &opcode{mne: "NOP*", wait: 3, addrMode: zpx, exec: c.nop}
	c.opCodes[0x75] = &opcode{mne: "ADC", wait: 4, addrMode: zpx, exec: c.adc}
	c.opCodes[0x76] = &opcode{mne: "ROR", wait: 6, addrMode: zpx, exec: c.ror}
	c.opCodes[0x77] = &opcode{mne: "RRA*", wait: 6, addrMode: zpx, exec: c.rra}
	c.opCodes[0x78] = &opcode{mne: "SEI", wait: 2, addrMode: none, exec: c.sei}
	c.opCodes[0x79] = &opcode{mne: "ADC", wait: 4, addrMode: aby, exec: c.adc}
	c.opCodes[0x7a] = &opcode{mne: "NOP*", wait: 2, addrMode: none, exec: c.nop}
	c.opCodes[0x7b] = &opcode{mne: "RRA*", wait: 7, addrMode: aby, exec: c.rra}
	c.opCodes[0x7c] = &opcode{mne: "NOP*", wait: 4, addrMode: abx, exec: c.nop}
	c.opCodes[0x7d] = &opcode{mne: "ADC", wait: 4, addrMode: abx, exec: c.adc}
	c.opCodes[0x7e] = &opcode{mne: "ROR", wait: 7, addrMode: abx, exec: c.ror}
	c.opCodes[0x7f] = &opcode{mne: "RRA*", wait: 7, addrMode: abx, exec: c.rra}

	c.opCodes[0x80] = &opcode{mne: "NOP*", wait: 2, addrMode: imm, exec: c.nop}
	c.opCodes[0x81] = &opcode{mne: "STA", wait: 6, addrMode: izx, exec: c.sta}
	c.opCodes[0x82] = &opcode{mne: "NOP*", wait: 2, addrMode: imm, exec: c.nop}
	c.opCodes[0x83] = &opcode{mne: "SAX*", wait: 6, addrMode: izx, exec: c.sax}
	c.opCodes[0x84] = &opcode{mne: "STY", wait: 3, addrMode: zp, exec: c.sty}
	c.opCodes[0x85] = &opcode{mne: "STA", wait: 3, addrMode: zp, exec: c.sta}
	c.opCodes[0x86] = &opcode{mne: "STX", wait: 3, addrMode: zp, exec: c.stx}
	c.opCodes[0x87] = &opcode{mne: "SAX*", wait: 3, addrMode: zp, exec: c.sax}
	c.opCodes[0x88] = &opcode{mne: "DEY", wait: 2, addrMode: none, exec: c.dey}
	c.opCodes[0x89] = &opcode{mne: "NOP*", wait: 2, addrMode: imm, exec: c.nop}
	c.opCodes[0x8a] = &opcode{mne: "TXA", wait: 2, addrMode: none, exec: c.txa}
	c.opCodes[0x8c] = &opcode{mne: "STY", wait: 4, addrMode: abs, exec: c.sty}
	c.opCodes[0x8d] = &opcode{mne: "STA", wait: 4, addrMode: abs, exec: c.sta}
	c.opCodes[0x8e] = &opcode{mne: "STX", wait: 4, addrMode: abs, exec: c.stx}
	c.opCodes[0x8f] = &opcode{mne: "SAX*", wait: 4, addrMode: abs, exec: c.sax}

	c.opCodes[0x90] = &opcode{mne: "BCC", wait: 2, addrMode: rel, exec: c.bcc}
	c.opCodes[0x91] = &opcode{mne: "STA", wait: 6, addrMode: izy, exec: c.sta}
	c.opCodes[0x94] = &opcode{mne: "STY", wait: 4, addrMode: zpx, exec: c.sty}
	c.opCodes[0x95] = &opcode{mne: "STA", wait: 4, addrMode: zpx, exec: c.sta}
	c.opCodes[0x96] = &opcode{mne: "STX", wait: 4, addrMode: zpy, exec: c.stx}
	c.opCodes[0x97] = &opcode{mne: "SAX*", wait: 4, addrMode: zpy, exec: c.sax}
	c.opCodes[0x98] = &opcode{mne: "TYA", wait: 2, addrMode: none, exec: c.tya}
	c.opCodes[0x99] = &opcode{mne: "STA", wait: 5, addrMode: aby, exec: c.sta}
	c.opCodes[0x9a] = &opcode{mne: "TXS", wait: 2, addrMode: none, exec: c.txs}
	c.opCodes[0x9c] = &opcode{mne: "SHY*", wait: 5, addrMode: abs, exec: c.shy}
	c.opCodes[0x9d] = &opcode{mne: "STA", wait: 5, addrMode: abx, exec: c.sta}
	c.opCodes[0x9e] = &opcode{mne: "SHX", wait: 5, addrMode: abs, exec: c.shx}

	c.opCodes[0xa0] = &opcode{mne: "LDY", wait: 2, addrMode: imm, exec: c.ldy}
	c.opCodes[0xa1] = &opcode{mne: "LDA", wait: 6, addrMode: izx, exec: c.lda}
	c.opCodes[0xa2] = &opcode{mne: "LDX", wait: 2, addrMode: imm, exec: c.ldx}
	c.opCodes[0xa3] = &opcode{mne: "LAX*", wait: 6, addrMode: izx, exec: c.lax}
	c.opCodes[0xa4] = &opcode{mne: "LDY", wait: 3, addrMode: zp, exec: c.ldy}
	c.opCodes[0xa5] = &opcode{mne: "LDA", wait: 3, addrMode: zp, exec: c.lda}
	c.opCodes[0xa6] = &opcode{mne: "LDX", wait: 3, addrMode: zp, exec: c.ldx}
	c.opCodes[0xa7] = &opcode{mne: "LAX*", wait: 3, addrMode: zp, exec: c.lax}
	c.opCodes[0xa8] = &opcode{mne: "TAY", wait: 2, addrMode: none, exec: c.tay}
	c.opCodes[0xa9] = &opcode{mne: "LDA", wait: 2, addrMode: imm, exec: c.lda}
	c.opCodes[0xaa] = &opcode{mne: "TAX", wait: 2, addrMode: none, exec: c.tax}
	c.opCodes[0xab] = &opcode{mne: "LAX*", wait: 2, addrMode: imm, exec: c.lax}
	c.opCodes[0xac] = &opcode{mne: "LDY", wait: 4, addrMode: abs, exec: c.ldy}
	c.opCodes[0xad] = &opcode{mne: "LDA", wait: 4, addrMode: abs, exec: c.lda}
	c.opCodes[0xae] = &opcode{mne: "LDX", wait: 4, addrMode: abs, exec: c.ldx}
	c.opCodes[0xaf] = &opcode{mne: "LAX*", wait: 4, addrMode: abs, exec: c.lax}

	c.opCodes[0xb0] = &opcode{mne: "BCS", wait: 2, addrMode: rel, exec: c.bcs}
	c.opCodes[0xb1] = &opcode{mne: "LDA", wait: 5, addrMode: izy, exec: c.lda}
	c.opCodes[0xb3] = &opcode{mne: "LAX*", wait: 5, addrMode: izy, exec: c.lax}
	c.opCodes[0xb4] = &opcode{mne: "LDY", wait: 4, addrMode: zpx, exec: c.ldy}
	c.opCodes[0xb5] = &opcode{mne: "LDA", wait: 4, addrMode: zpx, exec: c.lda}
	c.opCodes[0xb6] = &opcode{mne: "LDX", wait: 4, addrMode: zpy, exec: c.ldx}
	c.opCodes[0xb7] = &opcode{mne: "LAX*", wait: 4, addrMode: zpy, exec: c.lax}
	c.opCodes[0xb8] = &opcode{mne: "CLV", wait: 2, addrMode: none, exec: c.clv}
	c.opCodes[0xb9] = &opcode{mne: "LDA", wait: 4, addrMode: aby, exec: c.lda}
	c.opCodes[0xba] = &opcode{mne: "TSX", wait: 2, addrMode: none, exec: c.tsx}
	c.opCodes[0xbc] = &opcode{mne: "LDY", wait: 4, addrMode: abx, exec: c.ldy}
	c.opCodes[0xbd] = &opcode{mne: "LDA", wait: 4, addrMode: abx, exec: c.lda}
	c.opCodes[0xbe] = &opcode{mne: "LDX", wait: 2, addrMode: aby, exec: c.ldx}
	c.opCodes[0xbf] = &opcode{mne: "LAX*", wait: 4, addrMode: aby, exec: c.lax}

	c.opCodes[0xc0] = &opcode{mne: "CPY", wait: 2, addrMode: imm, exec: c.cpy}
	c.opCodes[0xc1] = &opcode{mne: "CMP", wait: 6, addrMode: izx, exec: c.cmp}
	c.opCodes[0xc2] = &opcode{mne: "NOP*", wait: 2, addrMode: imm, exec: c.nop}
	c.opCodes[0xc3] = &opcode{mne: "DCP*", wait: 8, addrMode: izx, exec: c.dcp}
	c.opCodes[0xc4] = &opcode{mne: "CPY", wait: 3, addrMode: zp, exec: c.cpy}
	c.opCodes[0xc5] = &opcode{mne: "CMP", wait: 3, addrMode: zp, exec: c.cmp}
	c.opCodes[0xc6] = &opcode{mne: "DEC", wait: 5, addrMode: zp, exec: c.dec}
	c.opCodes[0xc7] = &opcode{mne: "DCP*", wait: 5, addrMode: zp, exec: c.dcp}
	c.opCodes[0xc8] = &opcode{mne: "INY", wait: 2, addrMode: none, exec: c.iny}
	c.opCodes[0xc9] = &opcode{mne: "CMP", wait: 2, addrMode: imm, exec: c.cmp}
	c.opCodes[0xca] = &opcode{mne: "DEX", wait: 2, addrMode: none, exec: c.dex}
	c.opCodes[0xcb] = &opcode{mne: "AXS*", wait: 2, addrMode: imm, exec: c.axs}
	c.opCodes[0xcc] = &opcode{mne: "CPY", wait: 4, addrMode: abs, exec: c.cpy}
	c.opCodes[0xcd] = &opcode{mne: "CMP", wait: 4, addrMode: abs, exec: c.cmp}
	c.opCodes[0xce] = &opcode{mne: "DEC", wait: 6, addrMode: abs, exec: c.dec}
	c.opCodes[0xcf] = &opcode{mne: "DCP*", wait: 6, addrMode: abs, exec: c.dcp}

	c.opCodes[0xd0] = &opcode{mne: "BNE", wait: 2, addrMode: rel, exec: c.bne}
	c.opCodes[0xd1] = &opcode{mne: "CMP", wait: 5, addrMode: izy, exec: c.cmp}
	c.opCodes[0xd3] = &opcode{mne: "DCP*", wait: 8, addrMode: izy, exec: c.dcp}
	c.opCodes[0xd4] = &opcode{mne: "NOP*", wait: 3, addrMode: zpx, exec: c.nop}
	c.opCodes[0xd5] = &opcode{mne: "CMP", wait: 2, addrMode: zpx, exec: c.cmp}
	c.opCodes[0xd6] = &opcode{mne: "DEC", wait: 6, addrMode: zpx, exec: c.dec}
	c.opCodes[0xd7] = &opcode{mne: "DCP*", wait: 6, addrMode: zpx, exec: c.dcp}
	c.opCodes[0xd8] = &opcode{mne: "CLD", wait: 2, addrMode: none, exec: c.cld}
	c.opCodes[0xd9] = &opcode{mne: "CMP", wait: 4, addrMode: aby, exec: c.cmp}
	c.opCodes[0xda] = &opcode{mne: "NOP*", wait: 2, addrMode: none, exec: c.nop}
	c.opCodes[0xdb] = &opcode{mne: "DCP*", wait: 7, addrMode: aby, exec: c.dcp}
	c.opCodes[0xdc] = &opcode{mne: "NOP*", wait: 4, addrMode: abx, exec: c.nop}
	c.opCodes[0xdd] = &opcode{mne: "CMP", wait: 4, addrMode: abx, exec: c.cmp}
	c.opCodes[0xde] = &opcode{mne: "DEC", wait: 7, addrMode: abx, exec: c.dec}
	c.opCodes[0xdf] = &opcode{mne: "DCP*", wait: 7, addrMode: abx, exec: c.dcp}

	c.opCodes[0xe0] = &opcode{mne: "CPX", wait: 2, addrMode: imm, exec: c.cpx}
	c.opCodes[0xe1] = &opcode{mne: "SBC", wait: 6, addrMode: izx, exec: c.sbc}
	c.opCodes[0xe2] = &opcode{mne: "NOP*", wait: 2, addrMode: imm, exec: c.nop}
	c.opCodes[0xe3] = &opcode{mne: "ISC*", wait: 8, addrMode: izx, exec: c.isc}
	c.opCodes[0xe4] = &opcode{mne: "CPX", wait: 3, addrMode: zp, exec: c.cpx}
	c.opCodes[0xe5] = &opcode{mne: "SBC", wait: 3, addrMode: zp, exec: c.sbc}
	c.opCodes[0xe6] = &opcode{mne: "INC", wait: 5, addrMode: zp, exec: c.inc}
	c.opCodes[0xe7] = &opcode{mne: "ISC*", wait: 5, addrMode: zp, exec: c.isc}
	c.opCodes[0xe8] = &opcode{mne: "INX", wait: 2, addrMode: none, exec: c.inx}
	c.opCodes[0xe9] = &opcode{mne: "SBC", wait: 2, addrMode: imm, exec: c.sbc}
	c.opCodes[0xea] = &opcode{mne: "NOP", wait: 2, addrMode: none, exec: c.nop}
	c.opCodes[0xeb] = &opcode{mne: "SBC*", wait: 2, addrMode: imm, exec: c.sbc}
	c.opCodes[0xec] = &opcode{mne: "CPX", wait: 4, addrMode: abs, exec: c.cpx}
	c.opCodes[0xed] = &opcode{mne: "SBC", wait: 4, addrMode: abs, exec: c.sbc}
	c.opCodes[0xee] = &opcode{mne: "INC", wait: 6, addrMode: abs, exec: c.inc}
	c.opCodes[0xef] = &opcode{mne: "ISC*", wait: 6, addrMode: abs, exec: c.isc}

	c.opCodes[0xf0] = &opcode{mne: "BEQ", wait: 2, addrMode: rel, exec: c.beq}
	c.opCodes[0xf1] = &opcode{mne: "SBC", wait: 5, addrMode: izy, exec: c.sbc}
	c.opCodes[0xf3] = &opcode{mne: "ISC*", wait: 8, addrMode: izy, exec: c.isc}
	c.opCodes[0xf4] = &opcode{mne: "NOP*", wait: 3, addrMode: zpx, exec: c.nop}
	c.opCodes[0xf5] = &opcode{mne: "SBC", wait: 4, addrMode: zpx, exec: c.sbc}
	c.opCodes[0xf6] = &opcode{mne: "INC", wait: 6, addrMode: zpx, exec: c.inc}
	c.opCodes[0xf7] = &opcode{mne: "ISC*", wait: 6, addrMode: zpx, exec: c.isc}
	c.opCodes[0xf8] = &opcode{mne: "SED", wait: 2, addrMode: none, exec: c.sed}
	c.opCodes[0xf9] = &opcode{mne: "SBC", wait: 4, addrMode: aby, exec: c.sbc}
	c.opCodes[0xfa] = &opcode{mne: "NOP*", wait: 2, addrMode: none, exec: c.nop}
	c.opCodes[0xfb] = &opcode{mne: "ISC*", wait: 7, addrMode: aby, exec: c.isc}
	c.opCodes[0xfc] = &opcode{mne: "NOP*", wait: 4, addrMode: abx, exec: c.nop}
	c.opCodes[0xfd] = &opcode{mne: "SBC", wait: 4, addrMode: abx, exec: c.sbc}
	c.opCodes[0xfe] = &opcode{mne: "INC", wait: 7, addrMode: abx, exec: c.inc}
	c.opCodes[0xff] = &opcode{mne: "ISC*", wait: 7, addrMode: abx, exec: c.isc}
}
