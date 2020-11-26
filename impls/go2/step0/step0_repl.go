package main

import (
	"fmt"
	//"log"

	"mal/readline"
)

func read(input string) string {
	return input
}

func eval(input string) string {
	return input
}

func print(input string) string {
	return input
}

func rep(input string) string {
	return print(eval(read(input)))
}

func main() {
	defer readline.Close()
	for {
		var line string
		var err error
		//fmt.Print("user> ")
		//reader := bufio.NewReader(os.Stdin)
		line, err = readline.Prompt("user> ")
		if err != nil { // io.EOF, readline.ErrInterrupt
			break
			//log.Fatalf("%v",err)
		}
		fmt.Printf("%v\n", rep(line))
	}
}
