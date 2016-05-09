package worker

import (
//	"fmt"
	"log"
	"os/exec"
	//"runtime"
	"bytes"
)

type Worker struct {
	Command string
	Args    string
	Stdout  chan string
	Stderr  chan string
}

func (command *Worker) Run () {
	cmd := exec.Command(command.Command, command.Args)
	
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
		command.Stdout <- stdout.String()
		command.Stderr <- stderr.String()
}


func (command *Worker) RunFalse(url string, dir string){
	
	var bin = "false"
	var args = []string{"clone", url, dir}

	worker := &Worker{Command: "php", Args: "slowService.php", Output: c}
	
	cmd := exec.Command(bin, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		logger.Println(err)

		return error
	}
	logger.Println("git clone:" + stdout.String() + stderr.String())
}