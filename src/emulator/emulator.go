package emulator

import (
	"fmt"
	"os"
)

const (
	HeaderSize = 4
	DataSize   = 256

	OPCODE_NOP = 0x00
	OPCODE_STA = 0x10
	OPCODE_LDA = 0x20
	OPCODE_ADD = 0x30
	OPCODE_SUB = 0x31
	OPCODE_OR  = 0x40
	OPCODE_AND = 0x50
	OPCODE_NOT = 0x60
	OPCODE_JMP = 0x80
	OPCODE_JN  = 0x90
	OPCODE_JZ  = 0xA0
	OPCODE_HLT = 0xF0
	OPCODE_MUL = 0xB0
	OPCODE_DIV = 0xC0
)

type CPU struct {
	AC  byte
	PC  byte
	Mem []byte
}

func (c *CPU) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if len(data) < HeaderSize+DataSize {
		return fmt.Errorf("imagem pequena demais: %d bytes", len(data))
	}
	c.Mem = make([]byte, DataSize)
	copy(c.Mem, data[HeaderSize:HeaderSize+DataSize])
	c.AC, c.PC = 0, 0
	return nil
}

func (c *CPU) Run() {
	for {
		instr := c.Mem[c.PC]
		switch instr {
		case OPCODE_NOP:
			c.PC++
		case OPCODE_HLT:
			return
		default:
			operand := c.Mem[c.PC+1]
			switch instr {
			case OPCODE_STA:
				c.Mem[operand] = c.AC
			case OPCODE_LDA:
				c.AC = c.Mem[operand]
			case OPCODE_ADD:
				c.AC += c.Mem[operand]
			case OPCODE_SUB:
				c.AC -= c.Mem[operand]
			case OPCODE_OR:
				c.AC |= c.Mem[operand]
			case OPCODE_AND:
				c.AC &= c.Mem[operand]
			case OPCODE_NOT:
				c.AC = ^c.AC
			case OPCODE_JMP:
				c.PC = operand
				continue
			case OPCODE_JN:
				if c.AC&0x80 != 0 {
					c.PC = operand
					continue
				}
			case OPCODE_JZ:
				if c.AC == 0 {
					c.PC = operand
					continue
				}
			case OPCODE_MUL:
				c.AC *= c.Mem[operand]
			case OPCODE_DIV:
				if c.Mem[operand] != 0 {
					c.AC /= c.Mem[operand]
				} else {
					fmt.Println("divisão por zero")
				}
			default:
				fmt.Printf("opcode desconhecido 0x%X em PC=0x%X\n", instr, c.PC)
				return
			}
			c.PC += 2
		}
	}
}

func (c *CPU) Dump() {
	fmt.Printf("AC = 0x%02X, AC_valor = %d | PC = 0x%02X\n", c.AC, c.AC, c.PC)
	fmt.Println("Memória final:")
	for i := 0; i < len(c.Mem); i++ {
		if i%16 == 0 {
			fmt.Printf("\n%02X:", i)
		}
		fmt.Printf(" %02X", c.Mem[i])
	}
	fmt.Println()
}

func Run(path string) error {
	cpu := &CPU{}
	if err := cpu.Load(path); err != nil {
		return err
	}
	cpu.Run()
	cpu.Dump()
	return nil
}
