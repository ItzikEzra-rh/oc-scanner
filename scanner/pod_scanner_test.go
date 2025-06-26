package scanner_test

import (
	"bytes"
	"errors"
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"oc-scanner/scanner"
)

var _ = Describe("PodScanner", func() {
	var originalExecCommand func(string, ...string) *exec.Cmd

	BeforeEach(func() {
		originalExecCommand = scanner.ExecCommand
	})

	AfterEach(func() {
		scanner.ExecCommand = originalExecCommand
	})

	It("should run kubectl with correct arguments", func() {
		var calledArgs []string

		scanner.ExecCommand = func(name string, args ...string) *exec.Cmd {
			calledArgs = append([]string{name}, args...)

			return fakeCmd("mock output", nil)
		}

		ps := scanner.PodScanner{Namespace: "test-ns"}
		err := ps.Scan()

		Expect(err).To(BeNil())
		Expect(calledArgs).To(Equal([]string{
			"oc", "get", "pods", "-n", "test-ns", "-o", "wide",
		}))
	})

	It("should return error if oc fails", func() {
		scanner.ExecCommand = func(name string, args ...string) *exec.Cmd {
			return fakeCmd("", errors.New("command failed"))
		}

		ps := scanner.PodScanner{Namespace: "fail-ns"}
		err := ps.Scan()

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("error running kubectl"))
	})
})

func fakeCmd(output string, execErr error) *exec.Cmd {
	var cmdName string
	var cmdArgs []string

	if execErr != nil {
		cmdName = "false"
		cmdArgs = []string{}
	} else {
		cmdName = "echo"
		cmdArgs = []string{output}
	}

	cmd := exec.Command(cmdName, cmdArgs...)

	if execErr != nil {
		cmd.Stdout = bytes.NewBufferString("")
		cmd.Stderr = bytes.NewBufferString(execErr.Error())
	}

	return cmd
}
