package main

import (
	"flag"
	"fmt"
	"git.sda1.net/nexryai/bluemoon/core"
	"git.sda1.net/nexryai/bluemoon/services"
	"github.com/google/uuid"
	"os"
)

func main() {
	flag.Parse()
	operation := flag.Arg(0)

	if operation == "start" {
		id := uuid.New()
		config := core.LoadConfig("./bluemoon.yml")

		tmpRootDir := "/var/bluemoon/tmp/" + id.String()
		core.MsgInfo("tmpRootDir = " + tmpRootDir)

		services.CreateTmpfs(tmpRootDir, config.RamLimit)

		// どっかでエラー起こしてもクリーンアップされるようにする
		defer func() {
			if r := recover(); r != nil {
				core.MsgErr(fmt.Sprintf("Falat error occurred: %s", r.(error)))
				services.CleanTmpFiles(tmpRootDir)
				os.Exit(1)
			}
		}()

		services.ExtractDockerImage(config.DockerImage, tmpRootDir)

		core.WriteToFile(config.Exec, fmt.Sprintf("%s/%s.sh", tmpRootDir, id.String()))

		nspawnOpts := core.NspawnOpts{
			RootDirPath:  tmpRootDir,
			Uid:          "",
			BindDir:      config.BindDir,
			StartCommand: []string{"sh", fmt.Sprintf("/%s.sh", id.String())},
		}

		// コンテナ実行
		core.MsgInfo("Starting runtime...")
		services.StartNspawn(nspawnOpts)

		// クリーンアップ
		services.CleanTmpFiles(tmpRootDir)

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

		tmpRootTar := services.BuildAndExtractDockerImage(config.SrcDir)

		core.MsgInfo("Building bluemoon package...")
		services.BuildBluemoonPackageFromTar(tmpRootTar, config.PackageName+".sfs")

		core.MsgInfo("Done!")

	} else if operation == "clean" {
		core.MsgWarn("This operation MUST be executed after making sure that there is no running container !!!")
		services.CleanAllTmpFiles()
	} else {
		core.MsgErr("Invalid args")
	}
}
