// File: src/assembler/assembler.go
package assembler

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	HeaderSize = 4                // bytes de cabeçalho ("NDR")
	MemorySize = HeaderSize + 256 // total: header + 256 bytes
	DataStart  = 0x80             // onde começa a seção de dados
	// opcodes
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
)

// Assemble monta um .asm em um .mem com instruções de 2 bytes.
func Assemble(inputPath, outputPath string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("erro lendo %s: %w", inputPath, err)
	}
	lines := strings.Split(string(data), "\n")

	// --- PASSO 1: coleta símbolos da seção .DATA ---
	sym := map[string]int{}
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
				// nome DB valor
				name, valStr := parts[0], parts[2]
				val := 0
				if valStr != "?" {
					if v, e := strconv.Atoi(valStr); e == nil {
						val = v
					} else if v2, e2 := strconv.ParseInt(valStr, 0, 0); e2 == nil {
						val = int(v2)
					} else {
						return fmt.Errorf("valor inválido '%s': %w", valStr, e2)
					}
				}
				_ = val
				sym[name] = addrData
				addrData++
			}
		}
	}

	// --- PASSO 2: gera memória com duas passagens ---
	mem := make([]byte, MemorySize)
	// cabeçalho Neander
	copy(mem[:HeaderSize], []byte{0x03, 'N', 'D', 'R'})

	// primeiro, grava valores iniciais de DATA
	for k, addr := range sym {
		// leitura repetida de arquivo para pegar o valor DB
		// (poderia armazenar no passo 1, mas para simplificar:)
		for _, raw := range lines {
			if strings.HasPrefix(strings.TrimSpace(raw), k+" DB") {
				parts := strings.Fields(raw)
				valStr := parts[2]
				v := 0
				if valStr != "?" {
					if x, _ := strconv.Atoi(valStr); true {
						v = x
					}
				}
				mem[HeaderSize+addr] = byte(v)
				break
			}
		}
	}

	// agora monta código
	mode = ""
	pc := 0
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
		// .ORG
		if strings.HasPrefix(up, ".ORG") {
			f := strings.Fields(line)
			if len(f) == 2 {
				if o, e := strconv.ParseInt(f[1], 0, 0); e == nil {
					pc = int(o)
				}
			}
			continue
		}

		parts := strings.Fields(line)
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
		default:
			return fmt.Errorf("opcode desconhecido '%s'", op)
		}

		base := HeaderSize + pc
		// instrução com operando
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
			// apenas opcode
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
