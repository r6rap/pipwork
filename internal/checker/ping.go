package checker

import (
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type PINGResult struct {
	Status 	string		   // UP or DOWN
	Latency time.Duration
	Error 	string
}

func CheckPING(address string) PINGResult {
	var cmd *exec.Cmd
	// detect the operating system
	if runtime.GOOS == "windows" {
		// execute command
		cmd = exec.Command("ping", "-n", "1", "-w", "1000", address)
	} else {
		cmd = exec.Command("ping", "-c", "1", "-w", "1", address)
	}

	// save time when checking
	start := time.Now()
	// running cmd
	output, err := cmd.CombinedOutput()
	// response time received
	latency := time.Since(start)

	if err != nil || !strings.Contains(string(output), "TTL") {
		
		return PINGResult{Status: "DOWN", Latency: latency, Error: err.Error()}
	}

	return PINGResult{Status: "UP", Latency: latency}
}