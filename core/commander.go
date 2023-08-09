package core

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func ExecCommand(command string, args []string) {
	cmd := exec.Command(command, args...)
	stderr, err := cmd.CombinedOutput()

	ExitOnError(err, fmt.Sprintf("Failed to exec.　| \"%s\" >>> %s",
		strings.Join(cmd.Args, " "),
		string(stderr)))

}

// コマンドを実行して結果を一行ずつ配列に格納して返す
func ExecCommandGetResult(command string, args []string) []string {
	cmd := exec.Command(command, args...)
	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout
	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr

	// 実行
	err := cmd.Run()

	ExitOnError(err, fmt.Sprintf("Failed to exec.　| \"%s\" >>> %s",
		strings.Join(cmd.Args, " "),
		stderr.String()))

	output := strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\n")
	return output
}

// コマンドの実行結果をリアルタイムでstdoutに流す
func ExecCommandWithStdout(command string, args []string) {
	cmd := exec.Command(command, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		ExitOnError(err, "failed to create pipe(s)")
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		ExitOnError(err, "failed to create pipe(s)")
	}

	err = cmd.Start()
	if err != nil {
		ExitOnError(err, "failed to start command")
	}

	go copyOutput(stdout, os.Stdout)
	go copyOutput(stderr, os.Stderr)

	err = cmd.Wait()
	if err != nil {
		MsgErr("Command exit code != 0")
	}

	return
}

func copyOutput(src io.Reader, dest io.Writer) {
	_, err := io.Copy(dest, src)
	if err != nil {
		log.Fatal(err)
	}
}
