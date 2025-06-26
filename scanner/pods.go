package scanner

import (
	"fmt"
	"os/exec"
)

var ExecCommand = exec.Command

type PodScanner struct {
	Namespace string
}

func (p PodScanner) Scan() error {
	fmt.Println("Scanning pods")

	cmd := ExecCommand("oc", "get", "pods", "-n", p.Namespace, "-o", "wide")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running kubectl: %v", err)
	}

	fmt.Println(string(output))
	return nil
}
