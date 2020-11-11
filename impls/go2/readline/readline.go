package readline

import "github.com/peterh/liner"

var rline *liner.State

func init() {
	rline = liner.NewLiner()
	rline.SetCtrlCAborts(true)
}

func Prompt(prompt string) (string, error) {
	line, err := rline.Prompt(prompt)
	if err != nil { // io.EOF, readline.ErrInterrupt
		return "", err
		//log.Fatalf("%v",err)
	}
	rline.AppendHistory(line)
	return line, nil
}

func Close() {
	rline.Close()
}