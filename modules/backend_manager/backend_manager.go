package backend_manager

import (
	"chain_simulation/configs"
	"chain_simulation/modules/thread_manager"
	"chain_simulation/utils/dir"
	"context"
	"fmt"
	"os"
	"os/exec"
)

var BackendManagerInstance = NewBackendManager()

type BackendManager struct {
	Cancel context.CancelFunc
}

func NewBackendManager() *BackendManager {
	return &BackendManager{}
}

func StartBackendService(experimentIndex int) {
	thread_manager.ThreadManagerInstance.Add()
	go func() {
		defer func() {
			thread_manager.ThreadManagerInstance.Done()
		}()
		startBackendService(experimentIndex)
	}()
}

func StopBackendService() {
	BackendManagerInstance.Cancel()
}

func startBackendService(experimentIndex int) {
	var cmd *exec.Cmd
	var cmdPath = configs.TopConfigInstance.PathConfig.Cmd
	var ctx context.Context
	fmt.Println(cmdPath)
	{
		err := dir.WithContextManager(cmdPath, func() error {
			// 1. 创建命令
			ctx, BackendManagerInstance.Cancel = context.WithCancel(context.Background())
			cmd = exec.CommandContext(ctx, "./cmd", "http_service", "-e", fmt.Sprintf("%d", experimentIndex))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			// 2. 后台启动
			if err := cmd.Start(); err != nil {
				fmt.Printf("start backend service error: %v\n", err)
			}

			go func() {
				// 3. 阻塞直到完成
				err := cmd.Wait()
				if err != nil {
					fmt.Printf("backend service error: %v\n", err)
				}
			}()
			return nil
		})
		if err != nil {
			fmt.Printf("start backend service error: %v\n", err)
		}
	}
}
