# go-soft-token

A text-based cross-platform software TOTP implementation (compatible with Google Authenticator and Microsoft Authenticator) written in Go.
This software TOTP implementation uses scrypt for key stretching and AES for TOTP secret storage.

See: https://en.wikipedia.org/wiki/Software_token

This software is written in Go and uses modules (available in Go 1.11+)

This software has a text user interface (TUI) and uses [tview](https://github.com/rivo/tview/), which is based on [tcell](https://github.com/gdamore/tcell).

### Run and build

To run the software install Go and run:

    go run .

To create an executable, run:

    go build
    
Or to create a release in the `dist` folder, run:

    bash build.sh
    
This automatically increments the build version.

### Download

Go to the [releases](https://github.com/mevdschee/go-soft-token/releases) section to download binaries and source code (Assets).
