package main

import (
	"fmt"

	"mal/reader"
	"mal/readline"
)

func read(input string) (reader.MalType, error) {
	return reader.ReadStr(input)
	// if err != nil {
	// 	log.Printf("%v\n", err)
	// }
	// return res, err
}

func eval(input reader.MalType, err error) (reader.MalType, error) {
	return input, err
}

func print(input reader.MalType, err error) string {
	if err != nil {
		fmt.Printf("%v", err)
		return ""
	}
	return input.String()
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
