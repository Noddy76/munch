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

type Ram struct {
	bytes []uint8
}

func NewRam(size uint) *Ram {
	if size > 0x10000 {
		size = 0x10000
	}
	return &Ram{bytes: make([]uint8, size)}
}

func (r *Ram) Read(addr uint16) uint8 {
	return r.bytes[int(addr)%len(r.bytes)]
}

func (r *Ram) Write(addr uint16, v uint8) {
	r.bytes[int(addr)%len(r.bytes)] = v
}

type Rom struct {
	bytes []uint8
}

func NewRom(bytes []uint8) *Ram {
	return &Ram{bytes: bytes}
}

func (r *Rom) Read(addr uint16) uint8 {
	return r.bytes[int(addr)%len(r.bytes)]
}

func (r *Rom) Write(addr uint16, v uint8) {}
