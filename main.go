package main

import (
	"fmt"
	"os"
	"sync"

	"oc-scanner/scanner"

	"github.com/texttheater/golang-levenshtein/levenshtein"
)

func suggestClosest(input string, options []string) string {
	// Set minimum distance threshold - only suggest if difference is less than 3 characters
	minDistance := 3
	closest := ""

	// Iterate through all valid options and calculate edit distance
	for _, opt := range options {
		// Calculate Levenshtein distance between input and current option
		// Using runes to properly handle Unicode characters
		dist := levenshtein.DistanceForStrings([]rune(input), []rune(opt), levenshtein.DefaultOptions)

		// If this option is closer than our current best match, update it
		if dist < minDistance {
			minDistance = dist
			closest = opt
		}
	}

	return closest
}

func main() {
	// Validate command line arguments - need at least 4 args: program name, command, namespace, and at least one resource type
	if len(os.Args) < 4 {
		fmt.Println("Usage: kube-scanner scan <namespace> <resource-type> [<resource-type>...]")
		return
	}

	// Parse command line arguments
	command := os.Args[1]    // First argument should be the command (e.g., "scan")
	namespace := os.Args[2]  // Second argument is the Kubernetes namespace to scan
	resources := os.Args[3:] // Remaining arguments are the resource types to scan

	// Validate that the command is supported (currently only "scan" is implemented)
	if command != "scan" {
		fmt.Println("Error: supported command is 'scan'")
		return
	}

	// Resource type mapping - maps resource type names to factory functions that create appropriate scanners
	// This design pattern allows easy addition of new resource types by just adding entries to this map
	scannerMap := map[string]func(string) scanner.Scanner{
		"pods": func(ns string) scanner.Scanner {
			return scanner.PodScanner{Namespace: ns}
		},
		"deployments": func(ns string) scanner.Scanner {
			return scanner.DeploymentScanner{Namespace: ns}
		},
	}

	// Create a WaitGroup to coordinate concurrent scanning operations
	var wg sync.WaitGroup

	// Extract valid resource types from the scanner map for error suggestion purposes
	validTypes := make([]string, 0, len(scannerMap))
	for k := range scannerMap {
		validTypes = append(validTypes, k)
	}

	// Process each requested resource type
	for _, resource := range resources {
		// Check if the requested resource type is supported
		factory, ok := scannerMap[resource]
		if !ok {
			// Resource type not found - try to suggest a similar one
			closest := suggestClosest(resource, validTypes)
			if closest != "" {
				fmt.Printf("Unknown resource type: %s\nDid you mean: %s ?\n", resource, closest)
			} else {
				fmt.Printf("Unknown resource type: %s\n", resource)
			}
			continue // Skip this resource and move to the next one
		}

		// Create scanner instance for this resource type
		s := factory(namespace)

		// Launch concurrent scan operation for this resource
		wg.Add(1) // Increment wait group counter before starting goroutine
		go func(r string, sc scanner.Scanner) {
			defer wg.Done() // Ensure wait group is decremented when goroutine completes

			// Execute the scan with user-friendly status messages
			fmt.Printf("🔍 Scanning %s in namespace '%s'...\n", r, namespace)
			if err := sc.Scan(); err != nil {
				fmt.Printf("❌ Error scanning %s: %v\n", r, err)
			}
		}(resource, s) // Pass current values to avoid closure issues
	}

	// Wait for all concurrent scan operations to complete
	fmt.Println("⏳ Waiting for scans to complete...")
	wg.Wait()
	fmt.Println("✅ All scans completed.")
}
