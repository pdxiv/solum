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
| `00` | `nop_0` | no operation |
| `01` | `nop_1` | no operation |
| `02` | `or1` | flip low order bit of cx, cx ^= 1 |
| `03` | `shl` | shift left cx register, cx <<= 1 |
| `04` | `zero` | set cx register to zero, cx = 0 |
| `05` | `if_cz` | if cx==0 execute next instruction |
| `06` | `sub_ab` | subtract bx from ax, cx = ax - bx |
| `07` | `sub_ac` | subtract cx from ax, ax = ax - cx |
| `08` | `inc_a` | increment ax, ax = ax + 1 |
| `09` | `inc_b` | increment bx, bx = bx + 1 |
| `0a` | `dec_c` | decrement cx, cx = cx - 1 |
| `0b` | `inc_c` | increment cx, cx = cx + 1 |
| `0c` | `push_ax` | push ax on stack |
| `0d` | `push_bx` | push bx on stack |
| `0e` | `push_cx` | push cx on stack |
| `0f` | `push_dx` | push dx on stack |
| `10` | `pop_ax` | pop top of stack into ax |
| `11` | `pop_bx` | pop top of stack into bx |
| `12` | `pop_cx` | pop top of stack into cx |
| `13` | `pop_dx` | pop top of stack into dx |
| `14` | `jmp` | move ip to template |
| `15` | `jmpb` | move ip backward to template |
| `16` | `call` | call a procedure |
| `17` | `ret` | return from a procedure |
| `18` | `mov_cd` | move cx to dx, dx = cx |
| `19` | `mov_ab` | move ax to bx, bx = ax |
| `1a` | `mov_iab` | move instruction at address in bx to address in ax |
| `1b` | `adr` | address of nearest template to ax |
| `1c` | `adrb` | search backward for template |
| `1d` | `adrf` | search forward for template |
| `1e` | `mal` | allocate memory for daughter cell |
| `1f` | `divide` | cell division |

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
