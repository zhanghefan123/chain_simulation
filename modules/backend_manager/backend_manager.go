package backend_manager

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"chain_simulation/configs"
)

var defaultManager = NewBackendManager()

type BackendManager struct {
	mu        sync.Mutex
	cancel    context.CancelFunc
	command   *exec.Cmd
	waitGroup sync.WaitGroup
}

func NewBackendManager() *BackendManager {
	return &BackendManager{}
}

// StartBackendService starts the backend process for one experiment. The
// process runs asynchronously and is owned by the default backend manager.
func StartBackendService(experimentIndex int) error {
	return defaultManager.Start(experimentIndex)
}

// StopBackendService stops the active backend process and waits for it to exit.
// Calling StopBackendService when no process is running is safe.
func StopBackendService() error {
	return defaultManager.Stop()
}

func (manager *BackendManager) Start(experimentIndex int) error {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	if manager.command != nil {
		return fmt.Errorf("backend service is already running")
	}

	ctx, cancel := context.WithCancel(context.Background())
	command := exec.CommandContext(
		ctx,
		"./cmd",
		"http_service",
		"-e",
		fmt.Sprintf("%d", experimentIndex),
	)
	command.Dir = configs.TopConfigInstance.PathConfig.Cmd
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Start(); err != nil {
		cancel()
		return fmt.Errorf("start backend service in %s: %w", command.Dir, err)
	}

	manager.cancel = cancel
	manager.command = command
	manager.waitGroup.Add(1)
	go manager.waitForExit(ctx, command)
	return nil
}

func (manager *BackendManager) Stop() error {
	manager.mu.Lock()
	cancel := manager.cancel
	manager.mu.Unlock()

	if cancel == nil {
		return nil
	}

	cancel()
	manager.waitGroup.Wait()
	return nil
}

func (manager *BackendManager) waitForExit(ctx context.Context, command *exec.Cmd) {
	defer manager.waitGroup.Done()

	err := command.Wait()
	if err != nil && ctx.Err() == nil {
		fmt.Printf("backend service exited with error: %v\n", err)
	}

	manager.mu.Lock()
	if manager.command == command {
		manager.command = nil
		manager.cancel = nil
	}
	manager.mu.Unlock()
}
