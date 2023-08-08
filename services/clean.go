package services

import "git.sda1.net/nexryai/bluemoon/core"

func CleanTmpFiles(path string) {
	core.ExecCommand("umount", []string{path})
	core.ExecCommand("rm", []string{"-rf", path})
}

func CleanAllTmpFiles() {
	core.ExecCommand("umount", []string{"/var/bluemoon/tmp/*"})
	core.ExecCommand("rm", []string{"-rf", "/var/bluemoon/tmp/"})
}
