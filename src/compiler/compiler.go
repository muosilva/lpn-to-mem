package compiler

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Compile lê um .lpn, suporta + - * / e gera .asm com MUL/DIV
func Compile(inputPath, outputPath string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("falha ao ler %s: %w", inputPath, err)
	}
	lines := strings.Split(string(data), "\n")

	reProg := regexp.MustCompile(`^PROGRAMA\s+"(.+)"[:]?`)
	reAssignLit := regexp.MustCompile(`^(\w+)\s*=\s*(\d+)\s*$`)
	reAssignBin := regexp.MustCompile(`^(\w+)\s*=\s*(\w+)\s*([\+\-\*/])\s*(\w+)\s*$`)

	type litAssign struct {
		varName string
		val     int
	}
	type binAssign struct{ res, left, op, right string }
	var lits []litAssign
	var bins []binAssign
	vars := map[string]struct{}{}

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || line == "INICIO" || line == "FIM" {
			continue
		}
		if reProg.MatchString(line) {
			continue
		}
		if m := reAssignLit.FindStringSubmatch(line); m != nil {
			v := m[1]
			n, _ := strconv.Atoi(m[2])
			lits = append(lits, litAssign{v, n})
			vars[v] = struct{}{}
			continue
		}
		if m := reAssignBin.FindStringSubmatch(line); m != nil {
			bins = append(bins, binAssign{m[1], m[2], m[3], m[4]})
			vars[m[2]] = struct{}{}
			vars[m[4]] = struct{}{}
			vars[m[1]] = struct{}{}
			continue
		}
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("falha ao criar pasta para %s: %w", outputPath, err)
	}
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("falha ao criar %s: %w", outputPath, err)
	}
	defer f.Close()
	bw := bufio.NewWriter(f)

	// .DATA
	fmt.Fprintln(bw, ".DATA")
	seen := map[int]struct{}{}
	for _, L := range lits {
		if _, ok := seen[L.val]; !ok {
			fmt.Fprintf(bw, "CONST_%d DB %d\n", L.val, L.val)
			seen[L.val] = struct{}{}
		}
	}
	for v := range vars {
		fmt.Fprintf(bw, "%s DB ?\n", v)
	}
	fmt.Fprintln(bw)

	// .CODE
	fmt.Fprintln(bw, ".CODE")
	fmt.Fprintln(bw, ".ORG 0")
	for _, L := range lits {
		fmt.Fprintf(bw, "LDA CONST_%d\n", L.val)
		fmt.Fprintf(bw, "STA %s\n", L.varName)
	}
	for _, B := range bins {
		fmt.Fprintf(bw, "LDA %s\n", B.left)
		switch B.op {
		case "+":
			fmt.Fprintf(bw, "ADD %s\n", B.right)
		case "-":
			fmt.Fprintf(bw, "SUB %s\n", B.right)
		case "*":
			fmt.Fprintf(bw, "MUL %s\n", B.right)
		case "/":
			fmt.Fprintf(bw, "DIV %s\n", B.right)
		}
		fmt.Fprintf(bw, "STA %s\n", B.res)
	}
	fmt.Fprintln(bw, "HLT")

	if err := bw.Flush(); err != nil {
		return fmt.Errorf("falha ao escrever .asm: %w", err)
	}
	return nil
}
