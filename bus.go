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

type Ticker interface {
	Tick() error
}

type Addressable interface {
	Read(addr uint16) uint8
	Write(addr uint16, v uint8)
}

type memoryDevice struct {
	start  uint16
	end    uint16
	device Addressable
}

type Bus struct {
	tickCount uint64

	tickers []Ticker
	mems    []memoryDevice
}

func NewBus() *Bus {
	return &Bus{tickers: make([]Ticker, 0)}
}

// One tick of the clock
func (b *Bus) Tick() error {
	b.tickCount++
	for _, t := range b.tickers {
		if err := t.Tick(); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bus) TickCount() uint64 { return b.tickCount }

func (b *Bus) Addressable(start uint16, end uint16, device Addressable) {
	b.mems = append(b.mems, memoryDevice{start: start, end: end, device: device})
}

func (b *Bus) Ticker(t Ticker) {
	b.tickers = append(b.tickers, t)
}

func (b *Bus) Read(addr uint16) uint8 {
	d := b.findDevice(addr)
	return d.device.Read(addr - d.start)
}

func (b *Bus) Write(addr uint16, v uint8) {
	d := b.findDevice(addr)
	d.device.Write(addr-d.start, v)
}

func (b *Bus) findDevice(addr uint16) memoryDevice {
	for _, mem := range b.mems {
		if addr >= mem.start && addr <= mem.end {
			return mem
		}
	}
	return memoryDevice{start: 0, end: 0xffff, device: &NullDevice{}}
}

type NullDevice struct{}

func (*NullDevice) Read(uint16) uint8   { return 0 }
func (*NullDevice) Write(uint16, uint8) {}
