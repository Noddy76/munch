# munch

A 6502 emulator.

## Usage

```go
package main

import (
	"github.com/noddy76/munch"
)

var rom = []uint8{
	0xa9, 0x01, //       LDA #$01
	0x8d, 0x00, 0x02, // STA $0200
	0xa9, 0x05, //       LDA #$05
	0x8d, 0x01, 0x02, // STA $0201
	0xa9, 0x08, //       LDA #$08
	0x8d, 0x02, 0x02, // STA $0202
}

func main() {
	bus := munch.NewBus()
	bus.Addressable(0x0000, 0x3fff, munch.NewRam(0x40000))
	bus.Addressable(0x8000, 0x8fff, munch.NewRom(rom))
	bus.Addressable(0xfffa, 0xffff, munch.NewRom([]uint8{0x00, 0x00, 0x00, 0x80, 0x00, 0x00}))

	cpu := munch.NewCpu6502(bus)

	cpu.Debug = true

	for cpu.PC != uint16(0x8000+len(rom)) {
		bus.Tick()
	}
}
```

When run (because debug output from the CPU was turned on) you should see.

```
8000    LDA #$01                  PC:8002 A:01 X:00 Y:00 SP:fd nv‑BdIzc
8002    STA $0200                 PC:8005 A:01 X:00 Y:00 SP:fd nv‑BdIzc
8005    LDA #$05                  PC:8007 A:05 X:00 Y:00 SP:fd nv‑BdIzc
8007    STA $0201                 PC:800a A:05 X:00 Y:00 SP:fd nv‑BdIzc
800a    LDA #$08                  PC:800c A:08 X:00 Y:00 SP:fd nv‑BdIzc
800c    STA $0202                 PC:800f A:08 X:00 Y:00 SP:fd nv‑BdIzc
```

The central element is the `bus`. You can attach `Addressable`s to the bus sunch as the `Ram` and
`Rom` objects in this package. When the `Cpu6502` is created registers itself as a `Ticker` with
the bus. From then on every time `Tick()` is called on the `Bus` the call is propogated to every
registered `Ticker` object.

## References

* [Fergulator](https://github.com/scottferg/Fergulator) A NES emulator written in Go
* [redcode's 6502](https://github.com/redcode/6502) A portable 6502 emulator written in C
* [Nesdev Wiki](https://www.nesdev.org/wiki/Nesdev_Wiki)
