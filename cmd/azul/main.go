package main

import (
	_ "github.com/wbuchwalter/azul/cmd/azul/delete"
	_ "github.com/wbuchwalter/azul/cmd/azul/deploy"
	_ "github.com/wbuchwalter/azul/cmd/azul/logs"
	"github.com/wbuchwalter/azul/cmd/azul/root"
)

func main() {
	root.Execute()
}
