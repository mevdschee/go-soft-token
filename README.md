# go-soft-token

A software TOTP implementation that uses scrypt for key stretching and AES for TOTP secret storage.

See: https://en.wikipedia.org/wiki/Software_token

This software is written in Go and uses modules (available in Go 1.11+)

This software has a text user interface (TUI) and uses [tview](https://github.com/rivo/tview/), which is based on [tcell](https://github.com/gdamore/tcell).

