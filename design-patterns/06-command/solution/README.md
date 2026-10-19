# Command: Solution

## Intent

Command turns a request into an object. The object contains everything needed to execute and later reverse the request, while `CommandHistory` manages sequencing rather than knowing each operation's details.

## Go mapping

`setValueCommand` captures the receiver (`ValueStore`), input, and prior state. `CommandHistory` stores only the interface, pushes after a successful execution, and invokes `Undo` in last-in-first-out order.

## Why use it

Use Command for undo/redo, job queues, transaction-like actions, macro recording, and audit trails. It decouples the operation requester from the receiver.

## Trade-offs

Undo requires enough captured state to restore the prior value. Commands can become memory-heavy when their state is large, so real systems may use snapshots, event logs, or compensating actions instead. This compact example deliberately keeps redo out of scope.
