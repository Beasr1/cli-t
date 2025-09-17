package command

import (
	"fmt"
	"sort"
	"sync"
)

// Registry holds all available commands
type Registry struct {
	mu       sync.RWMutex
	commands map[string]Command
}

// global registry instance
var defaultRegistry = &Registry{
	commands: make(map[string]Command),
}

// Register adds a command to the registry
func Register(cmd Command) error {
	defaultRegistry.mu.Lock()
	defer defaultRegistry.mu.Unlock()

	name := cmd.Name()
	if _, exists := defaultRegistry.commands[name]; exists {
		return fmt.Errorf("command %s already registered", name)
	}

	defaultRegistry.commands[name] = cmd
	return nil
}

// Get retrieves a command by name
func Get(name string) (Command, bool) {
	defaultRegistry.mu.RLock()
	defer defaultRegistry.mu.RUnlock()

	cmd, ok := defaultRegistry.commands[name]
	return cmd, ok
}

// List returns all registered command names
func List() []string {
	defaultRegistry.mu.RLock()
	defer defaultRegistry.mu.RUnlock()

	names := make([]string, 0, len(defaultRegistry.commands))
	for name := range defaultRegistry.commands {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// GetAll returns all registered commands
func GetAll() map[string]Command {
	defaultRegistry.mu.RLock()
	defer defaultRegistry.mu.RUnlock()

	// Return a copy to prevent external modification
	cmds := make(map[string]Command, len(defaultRegistry.commands))
	for k, v := range defaultRegistry.commands {
		cmds[k] = v
	}
	return cmds
}
