//go:build windows
// +build windows

package networking

import (
	"os"
	"os/exec"
)

// Set platform specifiy custom process attributes
func SetProcessAttributes(cmd *exec.Cmd) {}

// Kills child processes on Linux. Windows doesn't provide a direct way to kill child processes, so we kill just the main process.
func KillProcess(process *os.Process) error {
	return windows.GenerateConsoleCtrlEvent(windows.CTRL_BREAK_EVENT, uint32(process.Pid))
}
