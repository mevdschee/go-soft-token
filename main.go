package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mevdschee/go-soft-token/file"
	"github.com/xlzd/gotp"

	"encoding/base32"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Config struct that holds account info
type Config struct {
	Accounts []Account `json:"accounts"`
}

// Account struct which contains a name
// and a secret
type Account struct {
	Name   string `json:"name"`
	Secret string `json:"secret"`
}

func readConfig(password, filename string) (Config, error) {
	var config Config

	data, err := file.Read(password, filename)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func writeConfig(config Config, filename string) error {

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = file.Write(password, filename, data)
	if err != nil {
		return err
	}

	return nil
}

const filename = "config.txt"

var config Config
var password string
var selectedIndex int
var spinnerIndex int

func main() {
	selectedIndex = 0
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, false).
				AddItem(nil, 0, 1, false), width, 1, false).
			AddItem(nil, 0, 1, false)
	}

	app := tview.NewApplication()
	pages := tview.NewPages()

	passwordForm := tview.NewForm()
	tokenForm := tview.NewForm()
	tokenText := tview.NewTextView()
	buttons := tview.NewForm()
	spinner := tview.NewTextView()
	confirm := tview.NewModal()
	warning := tview.NewModal()

	frame := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tokenText, 0, 1, false).
		AddItem(buttons, 3, 1, true)

	drawToken := func() {
		if len(config.Accounts) > 0 {
			selectedIndex = selectedIndex % len(config.Accounts)
			frame.SetTitle(fmt.Sprintf(" Token %d/%d ", selectedIndex+1, len(config.Accounts)))
			a := config.Accounts[selectedIndex]
			totp := gotp.NewDefaultTOTP(a.Secret)
			code, t := totp.NowWithExpiration()
			text := fmt.Sprintf("%s\n\n[yellow]%s[white] (%02d)", a.Name, code, (time.Now().Unix()-t)*-1)
			tokenText.SetText(text)
		} else {
			frame.SetTitle(" Token 0/0 ")
			tokenText.SetText("")
		}
	}

	drawSpinner := func() {
		spinnerIndex = (spinnerIndex + 1) % 4
		str := ""
		for i := 0; i < 4; i++ {
			if i == spinnerIndex {
				str += "o"
			} else {
				str += "."
			}
		}
		spinner.SetText(str)
	}

	updateTimer := func() {
		for {
			time.Sleep(500 * time.Millisecond)
			app.QueueUpdateDraw(func() {
				drawToken()
				drawSpinner()
			})
		}
	}

	confirmOkHandler := func() {
		pages.ShowPage("spinner")
		s := config.Accounts
		i := selectedIndex
		config.Accounts = append(s[:i], s[i+1:]...)
		err := writeConfig(config, filename)
		if err != nil {
			app.Stop()
			return
		}
		pages.HidePage("spinner")
		drawToken()
		buttons.SetFocus(0)
		app.SetFocus(buttons)
	}

	confirm.
		AddButtons([]string{"Ok", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.HidePage("confirm")
			if buttonLabel == "Ok" {
				go confirmOkHandler()
			} else {
				drawToken()
				buttons.SetFocus(0)
				app.SetFocus(buttons)
			}
		})

	warning.
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.HidePage("warning")
			tokenForm.SetFocus(1)
			app.SetFocus(tokenForm)
		})

	spinner.
		SetText("....").
		SetTextAlign(tview.AlignCenter).
		SetTitle(" Loading ").
		SetTitleColor(tcell.ColorYellow).
		SetBorder(true).
		SetBackgroundColor(tcell.NewHexColor(0x222222)).
		SetBorderPadding(4, 4, 4, 4)

	buttons.
		SetButtonsAlign(tview.AlignCenter).
		AddButton("Next", func() {
			if len(config.Accounts) == 0 {
				return
			}
			selectedIndex = (selectedIndex + 1) % len(config.Accounts)
			drawToken()
		}).
		AddButton("Add", func() {
			pages.ShowPage("tokenForm")
			tokenForm.SetFocus(0)
			app.SetFocus(tokenForm)
		}).
		AddButton("Del", func() {
			if len(config.Accounts) == 0 {
				return
			}
			a := config.Accounts[selectedIndex]
			confirm.SetText(fmt.Sprintf("Are you sure you want to delete token\nwith name '%s'?", a.Name))
			pages.ShowPage("confirm")
			confirm.SetFocus(0)
			app.SetFocus(confirm)
		}).
		SetBackgroundColor(tcell.NewHexColor(0x222222))

	tokenAddHandler := func() {
		pages.ShowPage("spinner")
		var a Account
		nameInput := tokenForm.GetFormItem(0).(*tview.InputField)
		secretInput := tokenForm.GetFormItem(1).(*tview.InputField)
		a.Name = nameInput.GetText()
		a.Secret = secretInput.GetText()
		a.Secret = strings.ToUpper(a.Secret)
		_, err := base32.StdEncoding.DecodeString(a.Secret)
		if err != nil {
			warning.SetText(fmt.Sprintf("error: %s", err))
			pages.HidePage("spinner")
			pages.ShowPage("warning")
			warning.SetFocus(0)
			app.SetFocus(warning)
			return
		}
		nameInput.SetText("")
		secretInput.SetText("")
		config.Accounts = append(config.Accounts, a)
		err = writeConfig(config, filename)
		if err != nil {
			app.Stop()
			return
		}
		pages.HidePage("spinner")
		pages.HidePage("tokenForm")
		drawToken()
		buttons.SetFocus(0)
		app.SetFocus(buttons)
	}

	tokenForm.
		SetButtonsAlign(tview.AlignRight).
		AddInputField("Name", "", 30, nil, nil).
		AddInputField("Secret", "", 30, nil, nil).
		AddButton("Ok", func() {
			go tokenAddHandler()
		}).
		AddButton("Cancel", func() {
			pages.HidePage("tokenForm")
			drawToken()
			buttons.SetFocus(0)
			app.SetFocus(buttons)
		}).
		SetTitle(" Add Token ").
		SetTitleColor(tcell.ColorYellow).
		SetBorder(true).
		SetBackgroundColor(tcell.NewHexColor(0x222222)).
		SetBorderPadding(2, 2, 3, 3)

	passwordSubmitHandler := func() {
		pages.ShowPage("spinner")
		passwordInput := passwordForm.GetFormItem(0).(*tview.InputField)
		password = passwordInput.GetText()
		var err error
		if _, err = os.Stat(filename); err != nil {
			writeConfig(config, filename)
			pages.HidePage("spinner")
			passwordInput.SetText("")
			passwordForm.SetFocus(0)
			app.SetFocus(passwordForm)
			return
		}
		config, err = readConfig(password, filename)
		if err != nil {
			pages.HidePage("spinner")
			passwordInput.SetText("")
			passwordForm.SetFocus(0)
			app.SetFocus(passwordForm)
		} else {
			pages.HidePage("spinner")
			pages.HidePage("passwordForm")
			drawToken()
			app.SetFocus(buttons)
		}
	}

	passwordForm.
		SetButtonsAlign(tview.AlignCenter).
		SetItemPadding(3).
		AddPasswordField("Password", "", 26, '*', nil).
		AddButton("Ok", func() { go passwordSubmitHandler() }).
		SetTitle(" Unlock ").
		SetTitleColor(tcell.ColorYellow).
		SetBorder(true).
		SetBackgroundColor(tcell.NewHexColor(0x222222)).
		SetBorderPadding(3, 0, 4, 4)

	passwordForm.GetFormItem(0).(*tview.InputField).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			go passwordSubmitHandler()
		}
	})

	tokenText.
		SetWrap(true).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetBackgroundColor(tcell.NewHexColor(0x222222))

	frame.
		SetTitle(" Token ").
		SetTitleColor(tcell.ColorYellow).
		SetBorder(true).
		SetBackgroundColor(tcell.NewHexColor(0x222222)).
		SetBorderPadding(1, 0, 2, 2)

	pages.
		AddPage("tokenText", modal(frame, 45, 11), true, true).
		AddPage("tokenForm", modal(tokenForm, 45, 11), true, false).
		AddPage("confirm", confirm, true, false).
		AddPage("warning", warning, true, false).
		AddPage("passwordForm", modal(passwordForm, 45, 11), true, true).
		AddPage("spinner", modal(spinner, 45, 11), true, false)

	tokenForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages.HidePage("tokenForm")
			drawToken()
			buttons.SetFocus(0)
			app.SetFocus(buttons)
			return nil
		}
		return event
	})

	exitOnEsc := func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRight {
			return tcell.NewEventKey(tcell.KeyTab, '\t', tcell.ModNone)
		}
		if event.Key() == tcell.KeyLeft {
			return tcell.NewEventKey(tcell.KeyBacktab, '\t', tcell.ModShift)
		}
		if event.Key() == tcell.KeyEscape {
			app.Stop()
			return nil
		}
		return event
	}
	buttons.SetInputCapture(exitOnEsc)
	passwordForm.SetInputCapture(exitOnEsc)

	go updateTimer()
	if err := app.SetRoot(pages, true).SetFocus(passwordForm).Run(); err != nil {
		panic(err)
	}
}
