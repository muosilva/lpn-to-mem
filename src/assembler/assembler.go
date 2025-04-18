package assembler

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	HeaderSize = 4                // bytes de cabeçalho
	MemorySize = HeaderSize + 256 // total: header + 256 bytes
	DataStart  = 0x80             // início da seção de dados

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

// Assemble monta um .asm em um .mem com suporte a labels e opcodes
func Assemble(inputPath, outputPath string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("erro lendo %s: %w", inputPath, err)
	}
	lines := strings.Split(string(data), "\n")

	//coleta símbolos de DATA
	sym := map[string]int{}
	dataVals := map[string]byte{}
	addrData := DataStart
	mode := ""
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}
		up := strings.ToUpper(line)
		switch {
		case strings.HasPrefix(up, ".DATA"):
			mode = "data"
			continue
		case strings.HasPrefix(up, ".CODE"):
			mode = ""
			continue
		}
		if mode == "data" {
			parts := strings.Fields(line)
			if len(parts) >= 3 && strings.ToUpper(parts[1]) == "DB" {
				name := parts[0]
				valStr := parts[2]
				val := byte(0)
				if valStr != "?" {
					if v, err := strconv.Atoi(valStr); err == nil {
						val = byte(v)
					} else if v2, err2 := strconv.ParseInt(valStr, 0, 0); err2 == nil {
						val = byte(v2)
					}
				}
				sym[name] = addrData
				dataVals[name] = val
				addrData++
			}
		}
	}

	//coleta labels em CODE
	labels := map[string]int{}
	pc := 0
	mode = ""
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}
		up := strings.ToUpper(line)
		switch {
		case strings.HasPrefix(up, ".DATA"):
			mode = ""
			continue
		case strings.HasPrefix(up, ".CODE"):
			mode = "code"
			continue
		}
		if mode != "code" {
			continue
		}
		if strings.HasPrefix(up, ".ORG") {
			f := strings.Fields(line)
			if len(f) == 2 {
				if o, err := strconv.ParseInt(f[1], 0, 0); err == nil {
					pc = int(o)
				}
			}
			continue
		}
		parts := strings.Fields(line)
		if len(parts) > 0 && strings.HasSuffix(parts[0], ":") {
			lbl := strings.TrimSuffix(parts[0], ":")
			labels[lbl] = pc
			continue
		}
		if len(parts) == 1 {
			pc++
		} else {
			pc += 2
		}
	}
	// mescla labels em sym
	for lbl, addr := range labels {
		sym[lbl] = addr
	}

	mem := make([]byte, MemorySize)
	copy(mem[:HeaderSize], []byte{0x03, 'N', 'D', 'R'})

	for name, addr := range sym {
		if val, ok := dataVals[name]; ok {
			mem[HeaderSize+addr] = val
		}
	}

	mode = ""
	pc = 0
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}
		up := strings.ToUpper(line)
		switch {
		case strings.HasPrefix(up, ".DATA"):
			mode = ""
			continue
		case strings.HasPrefix(up, ".CODE"):
			mode = "code"
			continue
		}
		if mode != "code" {
			continue
		}
		if strings.HasPrefix(up, ".ORG") {
			f := strings.Fields(line)
			if len(f) == 2 {
				if o, err := strconv.ParseInt(f[1], 0, 0); err == nil {
					pc = int(o)
				}
			}
			continue
		}
		parts := strings.Fields(line)
		if len(parts) > 0 && strings.HasSuffix(parts[0], ":") {
			continue // ignora labels
		}
		op := strings.ToUpper(parts[0])
		var code byte
		switch op {
		case "NOP":
			code = OPCODE_NOP
		case "STA":
			code = OPCODE_STA
		case "LDA":
			code = OPCODE_LDA
		case "ADD":
			code = OPCODE_ADD
		case "SUB":
			code = OPCODE_SUB
		case "OR":
			code = OPCODE_OR
		case "AND":
			code = OPCODE_AND
		case "NOT":
			code = OPCODE_NOT
		case "JMP":
			code = OPCODE_JMP
		case "JN":
			code = OPCODE_JN
		case "JZ":
			code = OPCODE_JZ
		case "HLT":
			code = OPCODE_HLT
		case "MUL":
			code = OPCODE_MUL
		case "DIV":
			code = OPCODE_DIV
		default:
			return fmt.Errorf("opcode desconhecido '%s'", op)
		}

		base := HeaderSize + pc
		if len(parts) == 2 {
			addr, ok := sym[parts[1]]
			if !ok {
				return fmt.Errorf("símbolo não definido '%s'", parts[1])
			}
			mem[base] = code
			pc++
			mem[HeaderSize+pc] = byte(addr)
			pc++
		} else {
			mem[base] = code
			pc++
		}
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(outputPath, mem, 0644); err != nil {
		return fmt.Errorf("falha ao escrever %s: %w", outputPath, err)
	}
	return nil
}
