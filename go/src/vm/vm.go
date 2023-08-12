package vm

const (
	Load  = 0x01
	Store = 0x02
	Add   = 0x03
	Sub   = 0x04
	Halt  = 0xff
)

// Stretch goals
const (
	Addi = 0x05
	Subi = 0x06
	Jump = 0x07
	Beqz = 0x08
)

// Given a 256 byte array of "memory", run the stored program
// to completion, modifying the data in place to reflect the result
//
// The memory format is:
//
// 00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f ... ff
// __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ ... __
// ^==DATA===============^ ^==INSTRUCTIONS==============^
//

func compute(memory []byte) {

	registers := [3]byte{8, 0, 0} // PC, R1 and R2

	// Keep looping, like a physical computer's clock
	for {
		pc := registers[0]
		op := memory[pc]
		switch op {
		case Load:
			rx := memory[pc + 1]
			mx := memory[pc + 2]
			registers[rx] = memory[mx]
		        registers[0] += 3
		case Store:
			rx := memory[pc + 1]
			mx := memory[pc + 2]
			// modifying memory outside of data segment
			// Halt immediately
			if mx > 8 {
				return
			}
			memory[mx] = registers[rx]
		        registers[0] += 3
		case Sub:
			r1 := memory[pc + 1]
			r2 := memory[pc + 2]
			registers[r1] -= registers[r2]
		        registers[0] += 3
		case Add:
			r1 := memory[pc + 1]
			r2 := memory[pc + 2]
			registers[r1] += registers[r2]
		        registers[0] += 3
		case Halt:
			return
		case Addi:
			rx := memory[pc + 1]
			registers[rx] += memory[pc + 2]
		        registers[0] += 3
		case Subi:
			rx := memory[pc + 1]
			registers[rx] -= memory[pc + 2]
		        registers[0] += 3
		case Jump:
			jmp := memory[pc + 1]
			registers[0] = jmp
		case Beqz:
			rx := memory[pc + 1]
			offset := memory[pc + 2]
			if registers[rx] == 0 {
				registers[0] += offset + 3
			} else {
				registers[0] += 3
			}
		}
	}
}








