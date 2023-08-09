package services

import "git.sda1.net/nexryai/bluemoon/core"

func CloneRepo(repo string, dir string) {
	core.ExecCommandWithStdout("git", []string{"clone", repo, dir})
}
