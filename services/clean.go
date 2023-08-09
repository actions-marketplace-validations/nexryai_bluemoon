package services

import (
	"git.sda1.net/nexryai/bluemoon/core"
	"os/exec"
)

func CleanTmpFiles(path string) {
	// マウントされて無くてエラーになっても続行したいので普通にexecする
	exec.Command("umount", path)
	core.ExecCommand("rm", []string{"-rf", path})
}
