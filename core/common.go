package core

import (
	"fmt"
	"github.com/google/uuid"
	"os"
)

var reset = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"
var yellow = "\033[33m"
var gray = "\033[37m"

// 現状だとデバッグモードはロギングにのみ影響するので公開していないが将来的にIsDebugModeになる可能性もある？
func isDebugMode() bool {
	return os.Getenv("LANCE_DEBUG_MODE") == "true"
}

func MsgInfo(text string) {
	fmt.Println(green + "✔ INFO: " + reset + text)
}

func MsgErr(text string) {
	fmt.Fprintln(os.Stderr, red+"✘ ERROR: "+text+reset)
}

func MsgWarn(text string) {
	fmt.Println(yellow + "⚠ WARNING: " + reset + text)
}

func MsgDebug(text string) {
	if isDebugMode() {
		fmt.Println(gray + "⚙ DEBUG: " + text + reset)
	}
}

func MsgDetail(text string) {
	fmt.Println(gray + "  ↳ " + reset + text)
}

func ExitOnError(err error, message string) {
	if err != nil {
		MsgErr(message)
		panic(err)
	}

	return
}

func GenUUID() string {
	id := uuid.New()
	return id.String()
}
