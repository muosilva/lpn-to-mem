# lpn-to-mem

## Aluno: Murilo Oliveira

Sistema compilador, montador e emulador que converte programas LPN (linguagem de programação simples) em código assembly e os executa em uma máquina virtual.

## Visão Geral

Este projeto é composto por três componentes principais:

1. **Compiler**: Converte arquivos `.lpn` em arquivos assembly (`.asm`)
2. **Assembler**: Converte arquivos assembly (`.asm`) em código de máquina (`.mem`)
3. **Emulator**: Executa os arquivos de código de máquina (`.mem`) em uma CPU virtual

## Como Usar

O projeto pode ser executado usando make:

```sh
make
```

Isso irá:
1. Ler o arquivo de entrada em `files/lpn/sample.lpn`
2. Gerar assembly em `files/asm/sample.asm`
3. Criar código de máquina em `files/mem/sample.mem`
4. Executar o programa no emulador

## Linguagem LPN

LPN é uma linguagem de programação simples que suporta:

- Declarações de variáveis
- Literais inteiros
- Operações aritméticas básicas (+, -, *, /)

Exemplo de programa:
```lpn
PROGRAMA "ExprMult":
INICIO
a = 12
b = 2
c = a + b
d = 6
RES = c + d
FIM
```

## Todo: 
- Permitir o uso de parentenses como operador. Atualmente isto NÃO funciona.
- Permitir a passagem de multiplas variaveis de uma vez em .lpn, ao invés de fazer como no sample.lpn