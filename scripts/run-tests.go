package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// Run tests with coverage for each package
	packages := []string{
		"./internal/models/...",
		"./internal/auth/...",
		"./internal/registry/...",
		"./internal/context/...",
		"./internal/storage/...",
		"./internal/server/...",
		"./internal/errors/...",
	}

	packageCount := 0

	for _, pkg := range packages {
		fmt.Printf("\n=== Testing %s ===\n", pkg)

		cmd := exec.Command("go", "test", "-v", "-short", "-coverprofile=coverage.out", pkg)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("Tests failed for %s: %v\n", pkg, err)
			continue
		}

		// Get coverage
		coverageCmd := exec.Command("go", "tool", "cover", "-func=coverage.out")
		output, err := coverageCmd.Output()
		if err == nil {
			fmt.Printf("Coverage for %s:\n%s\n", pkg, string(output))
		}

		packageCount++
	}

	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("Tested %d packages\n", packageCount)
}
