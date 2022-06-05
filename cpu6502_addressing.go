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

type addrMode struct {
	args int
	fmt  func(uint16, uint16) string           // func(argAddr, arg) => string
	addr func(*Cpu6502, uint16, uint16) uint16 // func(argAddr, arg) => address
}

var (
	none = addrMode{
		args: 0,
		fmt:  func(argAddr, arg uint16) string { return "" },
		addr: func(cpu *Cpu6502, argAddr, arg uint16) uint16 { return 0 },
	}
	imm = addrMode{
		args: 1,
		fmt:  func(argAddr, arg uint16) string { return fmt.Sprintf("#$%02x", arg) },
		addr: func(cpu *Cpu6502, argAddr, arg uint16) uint16 { return argAddr },
	}
	zp = addrMode{
		args: 1,
		fmt:  func(argAddr, arg uint16) string { return fmt.Sprintf("$%02x", arg) },
		addr: func(cpu *Cpu6502, argAddr, arg uint16) uint16 { return arg },
	}
	zpx = addrMode{
		args: 1,
		fmt:  func(argAddr, arg uint16) string { return fmt.Sprintf("$%02x,X", arg) },
		addr: func(cpu *Cpu6502, argAddr, arg uint16) uint16 { return uint16(uint8(arg) + cpu.X) },
	}
	zpy = addrMode{
		args: 1,
		fmt:  func(argAddr, arg uint16) string { return fmt.Sprintf("$%02x,Y", arg) },
		addr: func(cpu *Cpu6502, argAddr, arg uint16) uint16 { return uint16(uint8(arg) + cpu.Y) },
	}
	izx = addrMode{
		args: 1,
		fmt:  func(argAddr, arg uint16) string { return fmt.Sprintf("($%02x, X)", arg) },
		addr: func(cpu *Cpu6502, argAddr, arg uint16) uint16 {
			addr := uint8(arg) + cpu.X
			return uint16(cpu.bus.Read(uint16(addr))) + uint16(cpu.bus.Read(uint16(addr+1)))<<8
		},
	}
	izy = addrMode{
		args: 1,
		fmt:  func(argAddr, arg uint16) string { return fmt.Sprintf("($%02x),Y", arg) },
		addr: func(cpu *Cpu6502, argAddr, arg uint16) uint16 {
			a := uint16(cpu.bus.Read(arg)) + uint16(cpu.bus.Read(uint16(uint8(arg)+1)))<<8
			return a + uint16(cpu.Y)
		},
	}
	abs = addrMode{
		args: 2,
		fmt:  func(argAddr, arg uint16) string { return fmt.Sprintf("$%04x", arg) },
		addr: func(cpu *Cpu6502, argAddr, arg uint16) uint16 { return arg },
	}
	abx = addrMode{
		args: 2,
		fmt:  func(argAddr, arg uint16) string { return fmt.Sprintf("$%04x,X", arg) },
		addr: func(cpu *Cpu6502, argAddr, arg uint16) uint16 {
			addr := arg + uint16(cpu.X)
			if arg&0xff00 != addr&0xff00 {
				cpu.waitCycles += 1
			}
			return addr
		},
	}
	aby = addrMode{
		args: 2,
		fmt:  func(argAddr, arg uint16) string { return fmt.Sprintf("$%04x,Y", arg) },
		addr: func(cpu *Cpu6502, argAddr, arg uint16) uint16 {
			addr := arg + uint16(cpu.Y)
			if arg&0xff00 != addr&0xff00 {
				cpu.waitCycles += 1
			}
			return addr
		},
	}
	ind = addrMode{
		args: 2,
		fmt:  func(argAddr, arg uint16) string { return fmt.Sprintf("($%04x)", arg) },
		addr: func(cpu *Cpu6502, argAddr, arg uint16) uint16 {
			return uint16(cpu.bus.Read(arg&0xff00+uint16(uint8(arg)+1)))<<8 + uint16(cpu.bus.Read(arg))
		},
	}
	rel = addrMode{
		args: 1,
		fmt: func(argAddr, arg uint16) string {
			if arg < 0x80 {
				arg += argAddr
			} else {
				arg += argAddr - 0x100
			}
			return fmt.Sprintf("$%04x", arg+1)
		},
		addr: func(cpu *Cpu6502, argAddr, arg uint16) uint16 {
			if arg < 0x80 {
				arg += argAddr
			} else {
				arg += argAddr - 0x100
			}
			return arg + 1
		},
	}
)
