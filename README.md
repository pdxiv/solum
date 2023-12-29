# A Tierra implementation in Go

## Overview

Please note: **work in progress**

This document outlines an attempt to make a Go-based implementation of the Tierra artificial life simulation. In this simulation, virtual creatures exist within a digital environment, each powered by its own simulated CPU. These creatures interact with a shared memory space, termed the "soup", in a concurrent and controlled manner.

## CPU

Each creature is equipped with a CPU. The CPU is designed with the following components:

```golang
type CPU struct {
	AX int     // Address register
	BX int     // Address register
	CX int     // Numerical register
	DX int     // Numerical register
	SP byte    // Stack pointer
	ST [10]int // Ten-word stack
	IP int     // Instruction pointer
	FL bool    // Flags register for error conditions
}
```

## Instructions

| opcode | mnemonic | description |
| ------ | -------- | ----------- |
| `00` | `nop_0` | No operation performed |
| `01` | `nop_1` | No operation performed |
| `02` | `or1` | Toggle the least significant bit of CX |
| `03` | `shl` | Shift the CX register left by one bit |
| `04` | `zero` | Set CX register to zero |
| `05` | `if_cz` | Skip next instruction if CX is zero |
| `06` | `sub_ab` | Subtract BX from AX, store in CX |
| `07` | `sub_ac` | Subtract CX from AX, store in AX |
| `08` | `inc_a` | Increment AX by 1 |
| `09` | `inc_b` | Increment BX by 1 |
| `0a` | `dec_c` | Decrement CX by 1 |
| `0b` | `inc_c` | Increment DX by 1 |
| `0c` | `push_ax` | Push AX onto the stack |
| `0d` | `push_bx` | Push BX onto the stack |
| `0e` | `push_cx` | Push CX onto the stack |
| `0f` | `push_dx` | Push DX onto the stack |
| `10` | `pop_ax` | Pop from stack into AX |
| `11` | `pop_bx` | Pop from stack into BX |
| `12` | `pop_cx` | Pop from stack into CX |
| `13` | `pop_dx` | Pop from stack into DX |
| `14` | `jmp` | Jump forward |
| `15` | `jmpb` | Jump backward |
| `16` | `call` | Call a subroutine |
| `17` | `ret` | Return from a subroutine |
| `18` | `mov_cd` | Move value from CX to DX |
| `19` | `mov_ab` | Move value from AX to BX |
| `1a` | `mov_iab` | Move instruction from address in BX to address in AX |
| `1b` | `adr` | Find address of the nearest template, store in AX |
| `1c` | `adrb` | Search backward for a template, store in AX |
| `1d` | `adrf` | Search forward for a template, store in AX |
| `1e` | `mal` | Allocate memory for a new creature |
| `1f` | `divide` | Trigger cell division |

## Communication between the CPUs and the soup

To manage concurrent access to the shared "soup" memory by multiple CPUs, a dedicated soupAccess goroutine handles all read and write operations. This approach ensures synchronization and prevents access conflicts.

### Read Operation Structure

CPUs read from the soup using the `readSoup` channel:

```golang
type readSoup struct {
	index  int       // Address in the soup
	result chan byte // Channel for receiving read results
}
```

### Write Operation Structure

CPUs write to the soup using the `writeSoup` channel:

```golang
type writeSoup struct {
	index   int       // Address in the soup
	value   byte      // Value to write
	success chan bool // Channel for acknowledging write success
}
```

##  Reaping CPUs

Whenever a new CPU is "birthed", it is allocated a section of contigous memory and is added to the beginning of a "reaping queue". If there is not enough free space to allocate the contigous memory of the new CPU, the last entries in the "reaping queue" are removed until enough memory is available.

## Memory Management and CPU Lifecycle

### Reaping Queue

* New CPUs are added to a "reaping queue" upon creation.
* If insufficient memory is available, older CPUs at the end of the queue are removed to free up space.
* Each new CPU is allocated a contiguous memory block.

```golang
type CPUAllocation struct {
	StartIndex int       // Start index of the CPU's memory in the soup
	Length     int       // Length of the allocated memory block
	Terminate  chan bool // Channel to signal termination
}
```

```golang
var reaperQueue []CPUAllocation
```
