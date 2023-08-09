package services

import (
	"fmt"
	"git.sda1.net/nexryai/bluemoon/core"
)

func DownloadAndMountPackage(packageUrl string, id string, sizeLimit string) string {
	packageSfsPath := fmt.Sprintf("/var/bluemoon/tmp/%s/image/root.sfs", id)
	packageMountPoint := fmt.Sprintf("/var/bluemoon/tmp/%s/base", id)

	core.ExecCommand("wget", []string{packageUrl, "-o", packageSfsPath})

	CreateTmpfs(packageMountPoint, sizeLimit)
	core.ExecCommand("mount", []string{"-t", "squashfs", packageSfsPath, packageMountPoint})
	return packageMountPoint
}

// baseDirを下層、tmpfsを上層ディレクトリとする書き込み可能なディレクトリを作成する
func CreateOverlay(baseDir string, id string) string {

	// これが最終的なコンテナのエントリーポイント
	mountDir := fmt.Sprintf("/var/bluemoon/tmp/%s/root", id)

	overlayDir := fmt.Sprintf("/var/bluemoon/tmp/%s/over", id)
	workDir := fmt.Sprintf("/var/bluemoon/tmp/%s/work", id)

	CreateTmpfs(overlayDir, "100M")
	CreateTmpfs(workDir, "500M")

	mountOpts := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", baseDir, overlayDir, workDir)
	core.ExecCommand("mount", []string{"-t", "overlay", "overlay", "-o", mountOpts, mountDir})

	return mountDir
}
