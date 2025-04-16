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

// TokenType defines types of tokens in .lpn language
type TokenType int

const (
	TOKEN_EOF TokenType = iota
	TOKEN_IDENTIFIER
	TOKEN_NUMBER
	TOKEN_ASSIGN    // =
	TOKEN_PLUS      // +
	TOKEN_MINUS     // -
	TOKEN_MULT      // *
	TOKEN_DIV       // /
	TOKEN_LPAREN    // (
	TOKEN_RPAREN    // )
	TOKEN_SEMICOLON // ;
	TOKEN_INICIO    // INICIO
	TOKEN_FIM       // FIM
	TOKEN_PROGRAMA  // PROGRAMA
	TOKEN_QUOTE     // "
	TOKEN_COLON     // :
)

type Token struct {
	Type   TokenType
	Lexeme string
}

type ASTNode interface{}

type NumberNode struct{ Value int }

type VarNode struct{ Name string }

type BinOpNode struct {
	Op    TokenType
	Left  ASTNode
	Right ASTNode
}

type Statement struct {
	VarName string
	Expr    ASTNode
}

type Program struct {
	Name       string
	Statements []*Statement
	Result     ASTNode
}

func Compile(inputPath, outputPath string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("falha ao ler %s: %w", inputPath, err)
	}
	lines := strings.Split(string(data), "\n")

	reProg := regexp.MustCompile(`^PROGRAMA\s+"(.+)"[:]?`)
	reAssignLit := regexp.MustCompile(`^(\w+)\s*=\s*(\d+)\s*$`)
	reAssignBin := regexp.MustCompile(`^(\w+)\s*=\s*(\w+)\s*([\+\-])\s*(\w+)\s*$`)

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
		if m := reProg.FindStringSubmatch(line); m != nil {
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
			r, l, o, ri := m[1], m[2], m[3], m[4]
			bins = append(bins, binAssign{r, l, o, ri})
			vars[l] = struct{}{}
			vars[ri] = struct{}{}
			vars[r] = struct{}{}
			continue
		}
	}

	// Garante diretório de saída
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("falha ao criar pasta para %s: %w", outputPath, err)
	}
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("falha ao criar %s: %w", outputPath, err)
	}
	defer f.Close()
	bw := bufio.NewWriter(f)

	// Seção DATA
	fmt.Fprintln(bw, ".DATA")
	// constantes únicas
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

	// Seção CODE
	fmt.Fprintln(bw, ".CODE")
	fmt.Fprintln(bw, ".ORG 0")
	// atribuições literais
	for _, L := range lits {
		fmt.Fprintf(bw, "LDA CONST_%d\n", L.val)
		fmt.Fprintf(bw, "STA %s\n", L.varName)
	}
	// atribuições binárias
	for _, B := range bins {
		fmt.Fprintf(bw, "LDA %s\n", B.left)
		if B.op == "+" {
			fmt.Fprintf(bw, "ADD %s\n", B.right)
		} else {
			fmt.Fprintf(bw, "SUB %s\n", B.right)
		}
		fmt.Fprintf(bw, "STA %s\n", B.res)
	}
	fmt.Fprintln(bw, "HLT")

	if err := bw.Flush(); err != nil {
		return fmt.Errorf("falha ao escrever .asm: %w", err)
	}
	return nil
}
