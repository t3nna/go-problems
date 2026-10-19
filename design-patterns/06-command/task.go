package main

import "errors"

var ErrNothingToUndo = errors.New("nothing to undo")

type Command interface {
	Execute() error
	Undo() error
}
type ValueStore struct{ values map[string]string }

func NewValueStore() *ValueStore                    { return &ValueStore{values: make(map[string]string)} }
func (s *ValueStore) Get(key string) (string, bool) { value, ok := s.values[key]; return value, ok }

type CommandHistory struct{}

func (h *CommandHistory) Execute(Command) error              { return errors.New("TODO: execute command") }
func (h *CommandHistory) Undo() error                        { return ErrNothingToUndo }
func NewSetValueCommand(*ValueStore, string, string) Command { return todoCommand{} }

type todoCommand struct{}

func (todoCommand) Execute() error { return errors.New("TODO: set value") }
func (todoCommand) Undo() error    { return errors.New("TODO: undo value") }
