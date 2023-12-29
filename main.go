package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const STACK_SIZE = 10
const SOUP_SIZE = 60_000

type CPU struct {
	AX int             // Address register
	BX int             // Address register
	CX int             // Numerical register
	DX int             // Numerical register
	SP byte            // Stack pointer
	ST [STACK_SIZE]int // Stack
	IP int             // Instruction pointer
	FL bool            // Flags register for error conditions
}

func main() {
	soup := make([]byte, SOUP_SIZE)
	sizeOfCreature, creatureCpu, err := loadCreatureFromFile("ancestral_creature.asm", soup, 0)
	if err != nil {
		log.Fatal("Error reading creature from file:", err)
	}
	fmt.Println("DEBUG: size of loaded creature:", sizeOfCreature)
	fmt.Println("DEBUG: creature CPU:", creatureCpu)

	runCreature(&creatureCpu, soup)

}

func runCreature(cpu *CPU, soup []byte) {
	var opcode byte

	for {
		opcode = soup[cpu.IP]

		switch opcode {
		case 0x00:
			nop_0(cpu) // No operation performed
		case 0x01:
			nop_1(cpu) // No operation performed
		case 0x02:
			or1(cpu) // Toggle the least significant bit of CX
		case 0x03:
			shl(cpu) // Shift the CX register left by one bit
		case 0x04:
			zero(cpu) // Set CX register to zero
		case 0x05:
			if_cz(cpu) // Skip next instruction if CX is zero
		case 0x06:
			sub_ab(cpu) // Subtract BX from AX, store in CX
		case 0x07:
			sub_ac(cpu) // Subtract CX from AX, store in AX
		case 0x08:
			inc_a(cpu) // Increment AX by 1
		case 0x09:
			inc_b(cpu) // Increment BX by 1
		case 0x0a:
			dec_c(cpu) // Decrement CX by 1
		case 0x0b:
			inc_c(cpu) // Increment DX by 1
		case 0x0c:
			push_ax(cpu) // Push AX onto the stack
		case 0x0d:
			push_bx(cpu) // Push BX onto the stack
		case 0x0e:
			push_cx(cpu) // Push CX onto the stack
		case 0x0f:
			push_dx(cpu) // Push DX onto the stack
		case 0x10:
			pop_ax(cpu) // Pop from stack into AX
		case 0x11:
			pop_bx(cpu) // Pop from stack into BX
		case 0x12:
			pop_cx(cpu) // Pop from stack into CX
		case 0x13:
			pop_dx(cpu) // Pop from stack into DX
		case 0x14:
			jmp(cpu) // Jump forward
		case 0x15:
			jmpb(cpu) // Jump backward
		case 0x16:
			call(cpu) // Call a subroutine
		case 0x17:
			ret(cpu) // Return from a subroutine
		case 0x18:
			mov_cd(cpu) // Move value from CX to DX
		case 0x19:
			mov_ab(cpu) // Move value from AX to BX
		case 0x1a:
			mov_iab(cpu, soup) // Move instruction from address in BX to address in AX
		case 0x1b:
			adr(cpu, soup) // Find address of the nearest template, store in AX
		case 0x1c:
			adrb(cpu, soup) // Search backward for a template, store in AX
		case 0x1d:
			adrf(cpu, soup) // Search forward for a template, store in AX
		case 0x1e:
			mal(cpu) // Allocate memory for a new creature
		case 0x1f:
			divide(cpu) // Trigger cell division
		}
		cpu.IP++
	}
}

func loadCreatureFromFile(filePath string, soup []byte, baseIndex int) (int, CPU, error) {
	file, err := os.Open(filePath)
	var cpu CPU
	if err != nil {
		return 0, cpu, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	offset := 0
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, "nop_0"):
			soup[baseIndex+offset] = 0x00
		case strings.HasPrefix(line, "nop_1"):
			soup[baseIndex+offset] = 0x01
		case strings.HasPrefix(line, "or1"):
			soup[baseIndex+offset] = 0x02
		case strings.HasPrefix(line, "shl"):
			soup[baseIndex+offset] = 0x03
		case strings.HasPrefix(line, "zero"):
			soup[baseIndex+offset] = 0x04
		case strings.HasPrefix(line, "if_cz"):
			soup[baseIndex+offset] = 0x05
		case strings.HasPrefix(line, "sub_ab"):
			soup[baseIndex+offset] = 0x06
		case strings.HasPrefix(line, "sub_ac"):
			soup[baseIndex+offset] = 0x07
		case strings.HasPrefix(line, "inc_a"):
			soup[baseIndex+offset] = 0x08
		case strings.HasPrefix(line, "inc_b"):
			soup[baseIndex+offset] = 0x09
		case strings.HasPrefix(line, "dec_c"):
			soup[baseIndex+offset] = 0x0a
		case strings.HasPrefix(line, "inc_c"):
			soup[baseIndex+offset] = 0x0b
		case strings.HasPrefix(line, "push_ax"):
			soup[baseIndex+offset] = 0x0c
		case strings.HasPrefix(line, "push_bx"):
			soup[baseIndex+offset] = 0x0d
		case strings.HasPrefix(line, "push_cx"):
			soup[baseIndex+offset] = 0x0e
		case strings.HasPrefix(line, "push_dx"):
			soup[baseIndex+offset] = 0x0f
		case strings.HasPrefix(line, "pop_ax"):
			soup[baseIndex+offset] = 0x10
		case strings.HasPrefix(line, "pop_bx"):
			soup[baseIndex+offset] = 0x11
		case strings.HasPrefix(line, "pop_cx"):
			soup[baseIndex+offset] = 0x12
		case strings.HasPrefix(line, "pop_dx"):
			soup[baseIndex+offset] = 0x13
		case strings.HasPrefix(line, "jmpb"):
			soup[baseIndex+offset] = 0x15
		case strings.HasPrefix(line, "jmp"):
			soup[baseIndex+offset] = 0x14
		case strings.HasPrefix(line, "call"):
			soup[baseIndex+offset] = 0x16
		case strings.HasPrefix(line, "ret"):
			soup[baseIndex+offset] = 0x17
		case strings.HasPrefix(line, "mov_cd"):
			soup[baseIndex+offset] = 0x18
		case strings.HasPrefix(line, "mov_ab"):
			soup[baseIndex+offset] = 0x19
		case strings.HasPrefix(line, "mov_iab"):
			soup[baseIndex+offset] = 0x1a
		case strings.HasPrefix(line, "adrb"):
			soup[baseIndex+offset] = 0x1c
		case strings.HasPrefix(line, "adrf"):
			soup[baseIndex+offset] = 0x1d
		case strings.HasPrefix(line, "adr"):
			soup[baseIndex+offset] = 0x1b
		case strings.HasPrefix(line, "mal"):
			soup[baseIndex+offset] = 0x1e
		case strings.HasPrefix(line, "divide"):
			soup[baseIndex+offset] = 0x1f
		default:
			offset--
		}
		offset++
	}

	// Check for errors during Scan
	if err := scanner.Err(); err != nil {
		return 0, cpu, err
	}

	cpu.IP = baseIndex

	return offset, cpu, nil
}
