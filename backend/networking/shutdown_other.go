//go:build !windows

package networking

import (
	"os"
	"os/exec"

	"ganigeorgiev.com/pocketbase/pocketbase/core"
	"github.com/seriousm4x/upsnap/logging"
	"golang.org/x/sys/unix"
)

// Set platform specifiy custom process attributes
func SetProcessAttributes(cmd *exec.Cmd) {
	cmd.SysProcAttr = &unix.SysProcAttr{Setpgid: true}
}

// Kills child processes on Linux. Windows doesn't provide a direct way to kill child processes, so we kill just the main process.
func KillProcess(app core.App, process *os.Process) error {
	logging.Logger(app).Warn("Command did not finish in time and will be killed")
	pgid, err := unix.Getpgid(process.Pid)
	if err != nil {
		return err
	}

	err = unix.Kill(-pgid, unix.SIGTERM)
	if err != nil {
		return unix.Kill(-pgid, unix.SIGKILL)
	}

	return err
}
