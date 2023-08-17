package services

import "git.sda1.net/nexryai/bluemoon/core"

func ExtractDockerImage(imageName string, path string) {
	core.ExecCommandWithStdout("/usr/libexec/docker-image-extract", []string{"-o", path, imageName})
}

func BuildAndExtractDockerImage(srcDir string, platform string) string {
	// ToDo これ被らないようにする
	tmpFileName := "bm.tmproot.tar"
	containerTmpTag := "bm.buildtmp"

	core.MsgInfo("Building image with Docker...")

	dockerBuildArgs := []string{"build", srcDir, "-t", containerTmpTag, "--platform", "linux/" + platform}

	core.ExecCommandWithStdout("docker", dockerBuildArgs)

	core.MsgInfo("Extracting files from built image....")
	containerId := core.ExecCommandGetResult("docker", []string{"run", "--detach", containerTmpTag, "false"})
	core.ExecCommandWithStdout("docker", []string{"export", "-o", tmpFileName, containerId[0]})

	return tmpFileName
}
