package main

import (
	"fmt"
	"libnf/examples"
	"log"
	"os"
	"runtime/pprof"
)

func helpAndExit() {
	fmt.Println("Usage: <program> <option> [--profile]")
	fmt.Println("Options:")
	fmt.Println("  reader - Run reader()")
	fmt.Println("  writer - Run writer()")
	fmt.Println("  filtering - Run filtering()")
	fmt.Println("  --profile - Run with CPU and memory profiling")
	os.Exit(1)
}

func writeHeapProfile(f *os.File) {
	if f == nil {
		return
	}
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}

func main() {
	// Check if an argument is provided
	argc := len(os.Args)
	if argc < 2 || argc > 3 {
		helpAndExit()
	}
	// Parse the argument
	if argc == 3 {
		if os.Args[2] == "--profile" {
			os.Mkdir(".prof", 0755)
			cpu, err := os.Create(".prof/cpu.prof")
			if err != nil {
				log.Fatal("could not create CPU profile: ", err)
			}
			defer cpu.Close()

			if err := pprof.StartCPUProfile(cpu); err != nil {
				log.Fatal("could not start CPU profile: ", err)
			}
			defer pprof.StopCPUProfile()

			mem, err := os.Create(".prof/mem.prof")
			if err != nil {
				log.Fatal("could not create memory profile: ", err)
			}
			defer mem.Close()

			defer writeHeapProfile(mem)
		} else {
			fmt.Println("Invalid option.")
			helpAndExit()
		}
	}

	if os.Args[1] == "reader" {
		examples.Reader()
	} else if os.Args[1] == "writer" {
		examples.Writer()
	} else if os.Args[1] == "filtering" {
		examples.Filtering()
	} else if os.Args[1] == "sorting" {
		examples.Sorting()
	} else {
		fmt.Println("Invalid option.")
		helpAndExit()
	}
}
