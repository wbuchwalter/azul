package bluefunc

import (
	"os"
	"os/exec"
)

//compile file with GOOS=windows
//push (replace) to azure function

func Deploy(filePath string) error {
	cmd := exec.Command("go", "run", filePath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
