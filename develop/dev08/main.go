package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mitchellh/go-ps"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах
*/

var ErrNoPath = errors.New("path unspecified")

func execUnaryCmd(input string) error {
	input = strings.TrimRight(input, "\r\n")

	args := strings.Split(input, " ")

	switch args[0] {
	case "cd":
		if len(args) < 2 {
			return ErrNoPath
		}
		return os.Chdir(args[1])

	case "pwd":
		s, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(s)
		return nil

	case "echo":
		fmt.Println(strings.Join(args[1:], " "))
		return nil

	case "ps":
		p, err := ps.Processes()
		if err != nil {
			return err
		}
		fmt.Println("pid\tppid\tname")
		for _, pp := range p {
			fmt.Printf("%d\t%d\t%s\n", pp.Pid(), pp.PPid(), pp.Executable())
		}
		return nil

	case "kill":
		pid, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		p, err := os.FindProcess(pid)
		if err != nil {
			return err
		}
		return p.Kill()

	case "\\quit":
		os.Exit(0)
	}

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func main() {

	reader := bufio.NewReader(os.Stdin)
	for {
		s, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("\033[32m" + s + "\033[0m")
		fmt.Print("₽ ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		if strings.Contains(input, "|") {
			if err := execWithPipes(input); err != nil {
				fmt.Println(err)
			}
		} else {
			if err := execUnaryCmd(input); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func execWithPipes(input string) error {
	cmdsStrings := strings.Split(input, "|")

	cmds := make([]*exec.Cmd, 0, len(cmdsStrings))
	for _, cmdString := range cmdsStrings {

		cmdString = strings.Trim(cmdString, " \r\n")
		args := strings.Split(cmdString, " ")

		cmd := exec.Command(args[0], args[1:]...)
		cmds = append(cmds, cmd)
	}

	for i := 0; i < len(cmds)-1; i++ {
		rc, err := cmds[i].StdoutPipe()
		if err != nil {
			return err
		}
		cmds[i+1].Stdin = rc
	}

	// buf := &bytes.Buffer{}
	// cmds[len(cmds)-1].

	for i := 0; i < len(cmds)-1; i++ {
		if err := cmds[i].Start(); err != nil {
			return err
		}
		defer cmds[i].Wait()
	}

	b, err := cmds[len(cmds)-1].CombinedOutput()
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}
