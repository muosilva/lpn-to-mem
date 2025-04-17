package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/muosilva/lpn-to-mem/src/assembler"
	"github.com/muosilva/lpn-to-mem/src/compiler"
	"github.com/muosilva/lpn-to-mem/src/emulator"
)

func main() {
	start := time.Now()
	fmt.Printf("%s | Iniciando emulador...\n\n", start.Format("2006-01-02 15:04:05"))
	lpnDir := "files/lpn"
	asmDir := "files/asm"
	memDir := "files/mem"

	if err := os.MkdirAll(lpnDir, 0755); err != nil {
		log.Fatalf("falha ao criar diretório: %v", err)
	}

	if err := os.MkdirAll(asmDir, 0755); err != nil {
		log.Fatalf("falha ao criar diretório: %v", err)
	}

	if err := os.MkdirAll(memDir, 0755); err != nil {
		log.Fatalf("falha ao criar diretório: %v", err)
	}
	lpnFile := filepath.Join(lpnDir, "sample.lpn")
	asmFile := filepath.Join(asmDir, "sample.asm")
	memFile := filepath.Join(memDir, "sample.mem")

	if err := RunCompiler(lpnFile, asmFile); err != nil {
		log.Fatalf("falha ao compilar: %v", err)
	}
	if err := RunAssembler(asmFile, memFile); err != nil {
		log.Fatalf("falha ao montar: %v", err)
	}
	if err := RunEmulator(memFile); err != nil {
		log.Fatalf("falha ao emular: %v", err)
	}
	fmt.Println("Emulador concluído com sucesso!")
	fmt.Println("Arquivos gerados:")
	fmt.Printf(" - %s\n", asmFile)
	fmt.Printf(" - %s\n", memFile)

	elapsed := time.Since(start)
	fmt.Printf("Tempo total: %s\n", elapsed)
}

func RunCompiler(inputPath, outputPath string) error {
	if err := compiler.Compile(inputPath, outputPath); err != nil {
		return fmt.Errorf("falha na compilação: %w", err)
	}
	return nil
}
func RunAssembler(inputPath, outputPath string) error {
	if err := assembler.Assemble(inputPath, outputPath); err != nil {
		return fmt.Errorf("falha na montagem: %w", err)
	}
	return nil
}

func RunEmulator(inputPath string) error {
	if err := emulator.Run(inputPath); err != nil {
		return fmt.Errorf("falha na emulação: %w", err)
	}
	return nil
}
