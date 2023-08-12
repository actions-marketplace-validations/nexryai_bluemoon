package services

import (
	"fmt"
	"git.sda1.net/nexryai/bluemoon/core"
	"os"
)

func DownloadAndMountPackage(packageUrl string, id string, sizeLimit string) string {
	packageSfsPath := fmt.Sprintf("/var/bluemoon/images/%s.sfs", id)
	packageMountPoint := fmt.Sprintf("/var/bluemoon/tmp/%s/base", id)

	// Download if no exists
	_, err := os.Stat("/var/bluemoon/images")
	if err != nil {
		os.Mkdir("/var/bluemoon/images", 0644)
	}

	_, err = os.Stat(packageSfsPath)
	if err != nil {
		core.MsgInfo("Downloading package...")
		core.ExecCommand("wget", []string{packageUrl, "-O", packageSfsPath})
	}

	os.Mkdir(packageMountPoint, 0644)
	core.ExecCommand("mount", []string{"-t", "squashfs", packageSfsPath, packageMountPoint})
	return packageMountPoint
}

// baseDirを下層、tmpfsを上層ディレクトリとする書き込み可能なディレクトリを作成する
func CreateOverlay(baseDir string, id string) string {

	// これが最終的なコンテナのエントリーポイント
	mountDir := fmt.Sprintf("/var/bluemoon/tmp/%s/%s", id, id)

	overlayDir := fmt.Sprintf("/var/bluemoon/tmp/%s/over", id)
	workDir := fmt.Sprintf("/var/bluemoon/tmp/%s/work", id)

	os.Mkdir(mountDir, 0644)
	os.Mkdir(overlayDir, 0644)
	os.Mkdir(workDir, 0644)

	mountOpts := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", baseDir, overlayDir, workDir)
	core.ExecCommand("mount", []string{"-t", "overlay", "overlay", "-o", mountOpts, mountDir})

	return mountDir
}
