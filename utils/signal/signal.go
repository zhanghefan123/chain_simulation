package signal

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func SetupSignalHandler(cmd *exec.Cmd) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Printf("Received signal: %v", sig)

		if cmd != nil && cmd.Process != nil {
			// 向整个进程组发送信号
			err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
			if err != nil {
				fmt.Printf("Failed to terminate process: %v", err)
			}

			// 等待进程退出
			time.AfterFunc(5*time.Second, func() {
				err = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
				if err != nil {
					fmt.Printf("Failed to terminate process: %v", err)
				}
			})
		}
		os.Exit(0)
	}()
}
