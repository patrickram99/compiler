package main

import (
	"fmt"
	"os"
	"os/user"

	"main/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("¡Bienvenido %s! al primer lenguaje de programación de Taylor Swift!\n",
		user.Username)
	fmt.Printf("Burn some commands\n")
	repl.Start("main.sp", os.Stdout)
}
