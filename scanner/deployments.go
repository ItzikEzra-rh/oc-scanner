package scanner

import (
	"fmt"
	"os/exec"
)

type DeploymentScanner struct {
	Namespace string
}

func (d DeploymentScanner) Scan() error {
	fmt.Println("Scanning deployments")

	cmd := exec.Command("oc", "get", "deployments", "-n", d.Namespace, "-o", "wide")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running oc: %v", err)
	}

	fmt.Println(string(output))
	return nil
}
