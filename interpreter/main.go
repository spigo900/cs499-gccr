package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

// How much memory (i.e. how many 1-byte cells) the Brainfuck machine has.
const MEMORY_SIZE = 30000

func findMatchingLoopStart(prog []byte, pc int) int {
	extraUnmatched := 0
	out := -1
	for i := pc - 1; i >= 0; i-- {
		if prog[i] == ']' {
			extraUnmatched++
		} else if prog[i] == '[' {
			if extraUnmatched == 0 {
				out = i
				break
			} else {
				extraUnmatched--
			}
		}
	}
	return out
}

func findMatchingLoopEnd(prog []byte, pc int) int {
	extraUnmatched := 0
	out := -1
	for i := pc + 1; i < len(prog); i++ {
		if prog[i] == '[' {
			extraUnmatched++
		} else if prog[i] == ']' {
			if extraUnmatched == 0 {
				out = i
				break
			} else {
				extraUnmatched--
			}
		}
	}
	return out
}

func main() {
	var bytes [MEMORY_SIZE]byte
	readWriteHead := 0

	if len(os.Args) == 1 {
		fmt.Println("A Brainfuck interpreter written in Go.")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Printf("    $ %v program.bf\n", os.Args[0])
		os.Exit(0)
	} else if len(os.Args) != 2 {
		os.Stderr.WriteString(
			os.Args[0] + " takes only a single file as an argument.\n")
		os.Exit(1)
	}

	prog, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		os.Stderr.WriteString("Couldn't read file `" + os.Args[1] + "`.\n")
		os.Exit(1)
	}

	pc := 0
	for {
		if pc >= len(prog) {
			break
		}

		instruction := prog[pc]
		switch instruction {
		// INC & DEC
		case '+':
			bytes[readWriteHead]++
		case '-':
			bytes[readWriteHead]--

		// RIGHT & LEFT
		case '>':
			readWriteHead++
		case '<':
			readWriteHead--

		// INPUT & OUTPUT
		case ',':
			reader := bufio.NewReader(os.Stdin)
			bytes[readWriteHead], err = reader.ReadByte()
			if err != nil {
				os.Stderr.WriteString("Couldn't get input: " + err.Error())
				os.Exit(1)
			}
		case '.':
			fmt.Print(string(bytes[readWriteHead]))

		// LOOP
		case '[':
			// If the loop's over, find and skip to the end.
			if bytes[readWriteHead] == 0 {
				loopEnd := findMatchingLoopEnd(prog, pc)

				if loopEnd == -1 {
					os.Stderr.WriteString("Unmatched '['!\n")
					os.Exit(1)
				}

				pc = loopEnd + 1
				continue
			}
		case ']':
			if bytes[readWriteHead] != 0 {
				loopStart := findMatchingLoopStart(prog, pc)

				if loopStart == pc {
					os.Stderr.WriteString("Unmatched ']'!\n")
					os.Exit(1)
				}
				pc = loopStart + 1
				continue
			}

		// WHITESPACE
		case ' ':
		case '\t':
		case '\n':
		case '\r':

		// UNRECOGNIZED
		default:
			os.Stderr.WriteString("Unrecognized instruction `" +
				string(bytes[pc]) + "` at character " + fmt.Sprint(pc) + "!\n")
			os.Exit(1)
		}

		pc++
	}
}
