package networking

import (
	"os"
	"os/exec"

	"golang.org/x/sys/unix"

	"github.com/seriousm4x/upsnap/logger"
)

// Set platform specifiy custom process attributes
func SetProcessAttributes(cmd *exec.Cmd) {
	cmd.SysProcAttr = &unix.SysProcAttr{Setpgid: true}
}

// Kills child processes on Linux. Windows doesn't provide a direct way to kill child processes, so we kill just the main process.
func KillProcess(process *os.Process) error {
	pgid, err := unix.Getpgid(process.Pid)
	logger.Debug.Println(pgid)
	if err != nil {
		return err
	}
	return unix.Kill(-pgid, unix.SIGKILL)
}
