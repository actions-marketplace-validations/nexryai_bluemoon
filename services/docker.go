package services

import "git.sda1.net/nexryai/bluemoon/core"

func ExtractDockerImage(imageName string, path string) {
	core.ExecCommandWithStdout("/usr/libexec/docker-image-extract", []string{"-o", path, imageName})
}
