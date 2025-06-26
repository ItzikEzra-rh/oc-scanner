package scanner

import (
	"fmt"
	"os/exec"
)

type PodScanner struct {
	Namespace string
}

func (p PodScanner) Scan() error {
	fmt.Println("Scanning pods")

	cmd := exec.Command("oc", "get", "pods", "-n", p.Namespace, "-o", "wide")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running oc: %v", err)
	}

	fmt.Println(string(output))
	return nil
}
