package services

import "git.sda1.net/nexryai/bluemoon/core"

func StartNspawn(opts core.NspawnOpts) {
	var args []string

	args = append(args, "-D", opts.RootDirPath)

	if opts.Uid != "" {
		args = append(args, "--private-users="+opts.Uid)
	} else {
		args = append(args, "--private-users=pick")
	}

	if len(opts.BindDir) != 0 {
		for _, d := range opts.BindDir {
			args = append(args, "--bind="+d)
		}
	}

	args = append(args, "--as-pid2")
	args = append(args, opts.StartCommand...)

	core.ExecCommandWithStdout("systemd-nspawn", args)
}
