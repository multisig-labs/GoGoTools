package process

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type process struct {
	*exec.Cmd
}

// TODO better error handling, show user if process cant start etc

func NewProcess(command string, args []string, root string) *process {
	proc := &process{
		exec.Command(command, args...),
	}
	proc.Dir = root
	return proc
}

func (p *process) Running() bool {
	return p.Process != nil && p.ProcessState == nil
}

func (p *process) Start() error {
	fmt.Println(p.Cmd.String())
	return p.Cmd.Start()
}

func (p *process) Wait() error {
	return p.Cmd.Wait()
}

func (p *process) Kill() {
	if p.Running() {
		// p.signal(syscall.SIGKILL)
		p.signal(syscall.SIGTERM)
	}
}

func (p *process) signal(sig os.Signal) {
	group, err := os.FindProcess(-p.Process.Pid)
	if err != nil {
		return
	}

	if err = group.Signal(sig); err != nil {
		// fmt.Println(err)
	}
}
