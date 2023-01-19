package process

import (
	"fmt"
	"testing"
	"time"
)

func TestFoo(t *testing.T) {
	// p := NewProcess("/Users/john/Code/GoGoPool/template-subnetevm/node/bin/avalanchego", "/Users/john/Code/GoGoPool/template-subnetevm/node/bin")
	done := make(chan error)
	p := NewProcess("sleep", []string{"4"}, "")
	go func() {
		err := p.Start()
		if err != nil {
			done <- err
			return
		}
		fmt.Printf("Process started with PID: %d\n", p.Process.Pid)
		done <- p.Wait()
	}()

	go func() {
		time.Sleep(time.Second * 1)
		p.Kill()
		done <- fmt.Errorf("killed")
	}()

	err := <-done
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Println("Process completed.")
	fmt.Println(p.ProcessState)

	t.Fatal()
}
