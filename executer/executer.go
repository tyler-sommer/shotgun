package executer

import (
	"github.com/tyler-sommer/shotgun/config"
	"github.com/tyler-sommer/shotgun/model"
	"code.google.com/p/go.crypto/ssh"

	"fmt"
	"sync"
	"os"
)

type ExecutionMode int

const (
	Run ExecutionMode = iota
	Pause
	Halt
)

type Executer struct {
	config *config.Config

	servers []model.Server

	currMode ExecutionMode

	Mode chan ExecutionMode
}

func New(conf *config.Config, servers []model.Server) *Executer {
	return &Executer{conf, servers, Run, make(chan ExecutionMode)}
}

func (e *Executer) tick() {
	select {
	case mode := <- e.Mode:
		e.currMode = mode

	default:
		return
	}
}

func (e *Executer) Execute() {
	wg := sync.WaitGroup{}
	for _, server := range e.servers {
		wg.Add(1)
		fmt.Println("Spawning goroutine")
		go e.handle(server)
	}

	wg.Wait()
}

func (e *Executer) handle(server model.Server) {
	fmt.Println("Connecting to ", server)
	session := e.Connect(server)
	e.tick()
	if e.currMode == Halt {
		return
	}

	for _, script := range server.Scripts {
		if script.Enabled != true {
			continue
		}

		fmt.Println("Executing script")
		stdinPipe, _ := session.StdinPipe()
		var run func(command string) error
		if script.RequiresSudo == true {
			session.Start("sudo -i")
			run = func(command string) error {
				stdinPipe.Write([]byte(command+"\n"))

				return nil
			}
		} else {
			run = session.Run
		}

		for _, cmd := range script.Commands {
			e.tick()
			if e.currMode == Halt {
				return
			}

			fmt.Println("Executing command ", string(cmd))
			run(string(cmd))
		}

		if script.RequiresSudo == true {
			stdinPipe.Write([]byte("exit\n"))
		}
	}
}

func (e *Executer) Connect(server model.Server) *ssh.Session {
	conf := e.config.NewClientConfig(server.User)

	client, err := ssh.Dial("tcp", server.Host+":22", conf)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}

	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}

	session.RequestPty("bash", 80, 40, nil)
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	return session
}
