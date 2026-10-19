# Command

Implement undoable key-value updates. `CommandHistory.Execute` runs a command and records it only after success; `Undo` reverses the most recently successful command. `NewSetValueCommand` must restore a previous value on undo, or remove a key that did not exist before execution.

Return `ErrNothingToUndo` when the history is empty, and reject a nil command.

## Public API and example

Implement `CommandHistory.Execute`, `CommandHistory.Undo`, and `NewSetValueCommand` using the given `Command` and `ValueStore` types. Executing a command that sets `theme` to `dark`, then undoing it, removes `theme` if it was absent originally.

## Acceptance criteria

- Only successful commands enter history.
- Undo uses last-in-first-out order.
- An overwritten value is restored rather than deleted.

## Tags

`Design Patterns` `Command` `Undo`
