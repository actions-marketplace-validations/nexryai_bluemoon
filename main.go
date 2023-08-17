package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"git.sda1.net/nexryai/bluemoon/core"
	"git.sda1.net/nexryai/bluemoon/services"
	"github.com/google/uuid"
	"os"
	"os/exec"
)

func cleanAllFiles(tmpRoot string) {
	core.ExecCommand("umount", []string{"-l", tmpRoot})
}

func main() {
	flag.Parse()
	operation := flag.Arg(0)

	if operation == "start" {
		config := core.LoadConfig("./bluemoon.yml")
		id := fmt.Sprintf("%x", sha1.Sum([]byte(config.PackageUrl)))

		tmpRootDir := "/var/bluemoon/tmp/" + id
		core.MsgInfo("tmpRootDir = " + tmpRootDir)

		// ゴミが残ってたら消す
		_, err := os.Stat(tmpRootDir)
		if err != nil {
			exec.Command("umount", fmt.Sprintf("-l %s/%s", tmpRootDir, id))
			exec.Command("umount", fmt.Sprintf("-l %s", tmpRootDir))
		}

		services.CreateTmpfs(tmpRootDir, config.RamLimit)

		// どっかでエラー起こしてもクリーンアップされるようにする
		defer func() {
			if r := recover(); r != nil {
				core.MsgErr(fmt.Sprintf("Falat error occurred: %s", r.(error)))
				cleanAllFiles(tmpRootDir)
				os.Exit(1)
			}
		}()

		baseDir := services.DownloadAndMountPackage(config.PackageUrl, id, config.RamLimit)
		rootDir := services.CreateOverlay(baseDir, id)

		core.WriteToFile(config.Exec, fmt.Sprintf("%s/%s.sh", rootDir, id))

		nspawnOpts := core.NspawnOpts{
			RootDirPath:  rootDir,
			Uid:          "",
			BindDir:      config.BindDir,
			StartCommand: []string{"sh", fmt.Sprintf("/%s.sh", id)},
		}

		// コンテナ実行
		core.MsgInfo("Starting runtime...")
		services.StartNspawn(nspawnOpts)

		// クリーンアップ
		cleanAllFiles(tmpRootDir)

	} else if operation == "build" {
		core.MsgInfo("Build bluemoon package")

		defer func() {
			if r := recover(); r != nil {
				core.MsgErr(fmt.Sprintf("Falat error occurred: %s", r.(error)))
				os.Exit(1)
			}
		}()

		id := uuid.New()
		config := core.LoadBuildConfig("./bluemoon.build.yml")

		if config.SrcDir == "" {
			if config.SrcRepo != "" {
				core.MsgInfo(fmt.Sprintf("Build a package from repo (%s)", config.SrcRepo))

				config.SrcDir = "./src." + id.String()
				services.CloneRepo(config.SrcRepo, config.SrcDir)
			} else {
				core.MsgErr("Invalid build config file! There are no src.")
				os.Exit(1)
			}
		} else {
			core.MsgInfo(fmt.Sprintf("Build a package from src dir (%s)", config.SrcDir))
		}

		if config.EnableMultiPlatformBuild {
			core.MsgInfo("Use multiplatform build!")
		} else {
			config.Platforms = []string{"amd64"}
		}

		for _, p := range config.Platforms {
			tmpRootTar := services.BuildAndExtractDockerImage(config.SrcDir, p)
			core.MsgInfo("Building bluemoon package... - " + p)

			packageName := config.PackageName

			if config.EnableMultiPlatformBuild {
				packageName = fmt.Sprintf("%s-%s.sfs", packageName, p)
			} else {
				packageName = fmt.Sprintf("%s.sfs", packageName)
			}

			services.BuildBluemoonPackageFromTar(tmpRootTar, packageName)
		}

		core.MsgInfo("Done!")

	} else {
		core.MsgErr("Invalid args")
	}
}
