package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

//RCCConf commands are the commands to complete the compliation, Output is the file or directory with the output file(s)
type RCCConf struct {
	Commands []string
	Output   string
}

func compile(workingdir string) (string, string, error) {

	conf := &RCCConf{}

	b, err := ioutil.ReadFile(fmt.Sprintf("%s/rcc.json", workingdir))
	if err != nil {
		return "", "", err
	}

	err = json.Unmarshal(b, conf)
	if err != nil {
		return "", "", err
	}

	var output bytes.Buffer

	for i := 0; i < len(conf.Commands); i++ {
		args := strings.Split(conf.Commands[i], " ")
		var cmd *exec.Cmd
		if len(args) > 1 {
			cmd = exec.Command(args[0], args[1:]...)
		} else {
			cmd = exec.Command(args[0])
		}

		var outb, errb bytes.Buffer

		cmd.Dir = workingdir

		cmd.Stdout = &outb
		cmd.Stderr = &errb
		err := cmd.Run()
		if err != nil {
			output.WriteString(err.Error() + "\n")
		}
		output.WriteString(fmt.Sprintf("%s\n%s\n", outb.String(), errb.String()))
	}

	toReturn := strings.Replace(output.String(), "\n", "~n~", -1)

	return toReturn, fmt.Sprintf("%s/%s", workingdir, conf.Output), nil
}
