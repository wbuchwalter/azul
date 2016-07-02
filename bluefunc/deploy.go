package bluefunc

import (
	"os"
	"os/exec"
)

//compile file with GOOS=windows
//push (replace) to azure function

func Deploy(filePath string) error {
	cmd := exec.Command("sh", "-c", "GOOS=windows GOARCH=386 go build "+filePath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
