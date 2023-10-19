package networking

import (
	"os"
	"os/exec"

	"github.com/seriousm4x/upsnap/logger"
	"golang.org/x/sys/unix"
)

// Set platform specifiy custom process attributes
func SetProcessAttributes(cmd *exec.Cmd) {
	cmd.SysProcAttr = &unix.SysProcAttr{Setpgid: true}
}

// Kills child processes on Linux. Windows doesn't provide a direct way to kill child processes, so we kill just the main process.
func KillProcess(process *os.Process) error {
	logger.Warning.Println("Your shutdown cmd didn't finish in 2 minutes. It will be killed.")
	pgid, err := unix.Getpgid(process.Pid)
	if err != nil {
		return err
	}

	err = unix.Kill(-pgid, unix.SIGTERM)
	if err != nil {
		return unix.Kill(-pgid, unix.SIGKILL)
	}

	_, err = process.Wait()

	return err
}
