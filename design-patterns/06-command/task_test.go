package main

import (
	"errors"
	"testing"
)

func TestCommandHistoryExecutesAndUndoesSet(t *testing.T) {
	store := NewValueStore()
	history := &CommandHistory{}
	if err := history.Execute(NewSetValueCommand(store, "theme", "dark")); err != nil {
		t.Fatal(err)
	}
	if got, ok := store.Get("theme"); !ok || got != "dark" {
		t.Fatalf("store = %q, %v", got, ok)
	}
	if err := history.Undo(); err != nil {
		t.Fatal(err)
	}
	if _, ok := store.Get("theme"); ok {
		t.Fatal("key should be removed")
	}
}
func TestCommandHistoryUndoesInReverseOrder(t *testing.T) {
	store := NewValueStore()
	history := &CommandHistory{}
	_ = history.Execute(NewSetValueCommand(store, "a", "1"))
	_ = history.Execute(NewSetValueCommand(store, "b", "2"))
	if err := history.Undo(); err != nil {
		t.Fatal(err)
	}
	if _, ok := store.Get("b"); ok {
		t.Fatal("last command should be undone")
	}
	if got, ok := store.Get("a"); !ok || got != "1" {
		t.Fatalf("first value = %q, %v", got, ok)
	}
}
func TestCommandHistoryReportsEmptyUndo(t *testing.T) {
	if err := (&CommandHistory{}).Undo(); !errors.Is(err, ErrNothingToUndo) {
		t.Fatalf("error = %v", err)
	}
}
func TestCommandHistoryRejectsNilCommand(t *testing.T) {
	if err := (&CommandHistory{}).Execute(nil); err == nil {
		t.Fatal("expected error")
	}
}
