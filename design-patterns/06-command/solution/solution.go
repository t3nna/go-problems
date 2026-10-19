package main

import (
	"errors"
	"fmt"
)

var ErrNothingToUndo = errors.New("nothing to undo")

type Command interface {
	Execute() error
	Undo() error
}
type ValueStore struct{ values map[string]string }

func NewValueStore() *ValueStore                    { return &ValueStore{values: make(map[string]string)} }
func (s *ValueStore) Get(key string) (string, bool) { value, ok := s.values[key]; return value, ok }

type CommandHistory struct{ commands []Command }

func (h *CommandHistory) Execute(command Command) error {
	if command == nil {
		return fmt.Errorf("command is required")
	}
	if err := command.Execute(); err != nil {
		return err
	}
	h.commands = append(h.commands, command)
	return nil
}
func (h *CommandHistory) Undo() error {
	if len(h.commands) == 0 {
		return ErrNothingToUndo
	}
	index := len(h.commands) - 1
	command := h.commands[index]
	if err := command.Undo(); err != nil {
		return err
	}
	h.commands = h.commands[:index]
	return nil
}
func NewSetValueCommand(store *ValueStore, key, value string) Command {
	return &setValueCommand{store: store, key: key, value: value}
}

type setValueCommand struct {
	store      *ValueStore
	key, value string
	previous   string
	existed    bool
}

func (c *setValueCommand) Execute() error {
	if c.store == nil {
		return fmt.Errorf("store is required")
	}
	if c.key == "" {
		return fmt.Errorf("key is required")
	}
	c.previous, c.existed = c.store.values[c.key]
	c.store.values[c.key] = c.value
	return nil
}
func (c *setValueCommand) Undo() error {
	if c.store == nil {
		return fmt.Errorf("store is required")
	}
	if c.existed {
		c.store.values[c.key] = c.previous
	} else {
		delete(c.store.values, c.key)
	}
	return nil
}
