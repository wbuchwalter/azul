package core

import (
	"os"
	"os/exec"
)

func Deploy(filePath string) error {
	cmd := exec.Command("sh", "-c", "GOOS=windows GOARCH=386 go build "+filePath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

//push to https://${name}.scm.azurewebsites.net/api/vfs/site/wwwroot/${funcName}/${filepath}
/*
$username = "`$website"
$password = "pwd"
$base64AuthInfo = [Convert]::ToBase64String([Text.Encoding]::ASCII.GetBytes(("{0}:{1}" -f $username,$password)))
*/
