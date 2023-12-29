package main

import (
	"errors"
	"fmt"
)

// Main instruction functions

// 0x00: No operation performed
func nop_0(cpu *CPU) {
}

// 0x01: No operation performed
func nop_1(cpu *CPU) {
}

// 0x02: Toggle the least significant bit of CX
func or1(cpu *CPU) {
	cpu.CX ^= 1
}

// 0x03: Shift the CX register left by one bit
func shl(cpu *CPU) {
	cpu.CX = cpu.CX << 1
}

// 0x04: Set CX register to zero
func zero(cpu *CPU) {
	cpu.CX = 0
}

// 0x05: Skip next instruction if CX is zero
func if_cz(cpu *CPU) {
	if cpu.CX == 0 {
		cpu.IP++
	}
}

// 0x06: Subtract BX from AX, store in CX
func sub_ab(cpu *CPU) {
	cpu.CX = cpu.AX - cpu.BX
}

// 0x07: Subtract CX from AX, store in AX
func sub_ac(cpu *CPU) {
	cpu.AX = cpu.AX - cpu.CX
}

// 0x08: Increment AX by 1
func inc_a(cpu *CPU) {
	cpu.AX++
}

// 0x09: Increment BX by 1
func inc_b(cpu *CPU) {
	cpu.BX++
}

// 0x0a: Decrement CX by 1
func dec_c(cpu *CPU) {
	cpu.CX--
}

// 0x0b: Increment CX by 1
func inc_c(cpu *CPU) {
	cpu.CX++
}

// 0x0c: Push AX onto the stack
func push_ax(cpu *CPU) {
	if STACK_SIZE > cpu.SP+1 {
		cpu.ST[cpu.SP] = cpu.AX
		cpu.SP++
	} else {
		cpu.FL = true
	}
}

// 0x0d: Push BX onto the stack
func push_bx(cpu *CPU) {
	if STACK_SIZE > cpu.SP+1 {
		cpu.ST[cpu.SP] = cpu.BX
		cpu.SP++
	} else {
		cpu.FL = true
	}
}

// 0x0e: Push CX onto the stack
func push_cx(cpu *CPU) {
	if STACK_SIZE > cpu.SP+1 {
		cpu.ST[cpu.SP] = cpu.CX
		cpu.SP++
	} else {
		cpu.FL = true
	}
}

// 0x0f: Push DX onto the stack
func push_dx(cpu *CPU) {
	if STACK_SIZE > cpu.SP+1 {
		cpu.ST[cpu.SP] = cpu.DX
		cpu.SP++
	} else {
		cpu.FL = true
	}
}

// 0x10: Pop from stack into AX
func pop_ax(cpu *CPU) {
	if cpu.SP > 0 {
		cpu.SP--
		cpu.AX = cpu.ST[cpu.SP]
	} else {
		cpu.FL = true
	}
}

// 0x11: Pop from stack into BX
func pop_bx(cpu *CPU) {
	if cpu.SP > 0 {
		cpu.SP--
		cpu.BX = cpu.ST[cpu.SP]
	} else {
		cpu.FL = true
	}
}

// 0x12: Pop from stack into CX
func pop_cx(cpu *CPU) {
	if cpu.SP > 0 {
		cpu.SP--
		cpu.CX = cpu.ST[cpu.SP]
	} else {
		cpu.FL = true
	}
}

// 0x13: Pop from stack into DX
func pop_dx(cpu *CPU) {
	if cpu.SP > 0 {
		cpu.SP--
		cpu.DX = cpu.ST[cpu.SP]
	} else {
		cpu.FL = true
	}
}

// 0x14: Jump forward
func jmp(cpu *CPU) {
	fmt.Println("DEBUG: jmp not yet implemented")
}

// 0x15: Jump backward
func jmpb(cpu *CPU) {
	fmt.Println("DEBUG: jmpb not yet implemented")
}

// 0x16: Call a subroutine
func call(cpu *CPU) {
	fmt.Println("DEBUG: call not yet implemented")
}

// 0x17: Return from a subroutine
func ret(cpu *CPU) {
	fmt.Println("DEBUG: ret not yet implemented")
}

// 0x18: Move value from CX to DX
func mov_cd(cpu *CPU) {
	cpu.DX = cpu.CX
}

// 0x19: Move value from AX to BX
func mov_ab(cpu *CPU) {
	cpu.BX = cpu.AX
}

// 0x1a: Move instruction from address in BX to address in AX
func mov_iab(cpu *CPU, soup []byte) {
	instruction := soup[wrapSoupAddress(cpu.BX)]
	soup[wrapSoupAddress(cpu.AX)] = instruction
}

// 0x1b: Find address of the nearest template, store in AX
func adr(cpu *CPU, soup []byte) {
	fmt.Println("DEBUG: adr not yet implemented")
}

// 0x1c: Search backward for a template, store in AX
func adrb(cpu *CPU, soup []byte) {
	template, size := resolveTemplateInput(cpu.IP, soup)
	templateAddress, err := searchComplementBackward(cpu, soup, template, size)
	if err != nil {
		cpu.FL = true
	} else {
		cpu.AX = templateAddress
	}
}

// 0x1d: Search forward for a template, store in AX
func adrf(cpu *CPU, soup []byte) {
	template, size := resolveTemplateInput(cpu.IP, soup)
	templateAddress, err := searchComplementForward(cpu, soup, template, size)
	if err != nil {
		cpu.FL = true
	} else {
		cpu.AX = templateAddress
	}
}

// 0x1e: Allocate memory for a new creature
func mal(cpu *CPU) {
	fmt.Println("DEBUG: mal not yet implemented")
}

// 0x1f: Trigger cell division
func divide(cpu *CPU) {
	fmt.Println("DEBUG: divide not yet implemented")
}

// Helper functions that assist the main instruction functions

func wrapSoupAddress(address int) int {
	// Make sure that an address wraps to the soup address space
	return (address%SOUP_SIZE + SOUP_SIZE) % SOUP_SIZE
}

func resolveTemplateInput(ip int, soup []byte) (byte, int) {
	// Find the template by searching ahead for all matching NOPs
	searchIndex := 1
	var template byte
	soupLocation := (ip + searchIndex + SOUP_SIZE) % SOUP_SIZE
	for soup[soupLocation] == 0 || soup[soupLocation] == 1 {
		template = template << 1
		template = template | soup[soupLocation]
		searchIndex++
		soupLocation = (ip + searchIndex + SOUP_SIZE) % SOUP_SIZE
	}
	return template, searchIndex - 1
}

func searchComplementBackward(cpu *CPU, soup []byte, template byte, size int) (int, error) {
	templateComplement := complementOfTemplate(template, size)

	for i := 1; i < SOUP_SIZE; i++ {
		foundTemplate, foundSize := resolveTemplateInput(cpu.IP-i, soup)
		if (size == foundSize) && (templateComplement == foundTemplate) {

			return cpu.IP - i + 1, nil
		}
	}
	return 0, errors.New("no matching template complement found")
}

func searchComplementForward(cpu *CPU, soup []byte, template byte, size int) (int, error) {
	templateComplement := complementOfTemplate(template, size)

	for i := 1; i < SOUP_SIZE; i++ {
		foundTemplate, foundSize := resolveTemplateInput(cpu.IP+i, soup)
		if (size == foundSize) && (templateComplement == foundTemplate) {

			return cpu.IP + i + 1, nil
		}
	}
	return 0, errors.New("no matching template complement found")
}

func complementOfTemplate(template byte, size int) byte {
	var mask byte
	mask = ^(^mask << size)
	output := (^template) & mask
	return output
}
