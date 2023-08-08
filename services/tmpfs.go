package services

import (
	"git.sda1.net/nexryai/bluemoon/core"
	"os"
)

func CreateTmpfs(path string, sizeLimit string) {
	err := os.MkdirAll(path, 0644)
	core.ExitOnError(err, "Failed to create tmpDir")

	core.ExecCommand("mount", []string{"-t", "tmpfs", "-o", "size=" + sizeLimit, "tmpfs", path})
}
