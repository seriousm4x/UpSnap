//go:build windows

package networking

import "golang.org/x/sys/windows"

func isRoot() bool {
	var sid *windows.SID
	sid, err := windows.CreateWellKnownSid(windows.WinBuiltinAdministratorsSid)
	if err != nil {
		return false
	}
	token := windows.GetCurrentProcessToken()
	isAdmin, err := token.IsMember(sid)
	return err == nil && isAdmin
}
