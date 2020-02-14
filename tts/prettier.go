package tts

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	pgs "github.com/lyft/protoc-gen-star"
)

// NewFormatter calls prettier on the generated client
func NewFormatter() pgs.PostProcessor { return formatter{} }

type formatter struct{}

// Match returns true only for Custom and Generated files (including templates).
func (cpp formatter) Match(a pgs.Artifact) bool {
	switch a := a.(type) {
	case pgs.GeneratorFile:
		return filepath.Ext(a.Name) == ".ts"
	default:
		return false
	}
}

// Process attaches the copyright header to the top of the input bytes
func (cpp formatter) Process(in []byte) ([]byte, error) {
	cmd := exec.Command(
		"prettier",
		"--stdin",
		"--parser",
		"typescript",
		"--trailing-comma",
		"all",
		"--tab-width",
		"4",
		"--single-quote",
		"true",
		"--print-width",
		"120",
	)
	cmd.Stdin = bytes.NewReader(in)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return in, nil
	}
	return out, nil
}
