package services

import (
	"git.sda1.net/nexryai/bluemoon/core"
	"os"
)

func BuildBluemoonPackageFromTar(src string, out string) {
	tmpDir := "tmp." + core.GenUUID()

	core.ExecCommand("mkdir", []string{tmpDir})
	core.ExecCommand("tar", []string{"-C", tmpDir, "-xf", src})

	// GitHub Actionとかだとストレージに限りがあるので削除
	err := os.Remove(src)
	if err != nil {
		core.ExitOnError(err, "failed to remove tmp file!")
	}

	core.MsgInfo("Compressing files....")
	core.ExecCommandWithStdout("mksquashfs", []string{tmpDir, out, "-comp", "xz"})
	core.ExecCommand("rm", []string{"-rf", tmpDir})
}
