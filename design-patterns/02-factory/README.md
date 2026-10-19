# Factory

Implement `NewNotifier(channel Channel) (Notifier, error)`. The factory must return the right notifier for `EmailChannel` and `SMSChannel`, and reject unknown channels. A notifier validates that both notification recipient and body are present, then returns `email:<recipient>:<body>` or `sms:<recipient>:<body>`.

Callers must depend only on the `Notifier` interface; they must not need to know the concrete email or SMS type.

## Public API

Implement `NewNotifier(Channel) (Notifier, error)` and preserve the provided `Channel`, `Notification`, and `Notifier` types. For example, an email notifier sending to `ada@example.com` with body `hello` returns `email:ada@example.com:hello`.

## Acceptance criteria

- Supported channels create behaviorally distinct notifiers.
- Unknown channels return an error.
- Validation belongs to the selected notifier.

## Tags

`Design Patterns` `Factory`
