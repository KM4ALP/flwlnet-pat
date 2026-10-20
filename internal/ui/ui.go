package ui

import (
	"fmt"
	"os"

	"flwlnet-pat/internal/mbox"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type UI struct {
	app     fyne.App
	win     fyne.Window
	content fyne.CanvasObject

	title   *widget.Label
	msgList *widget.List
	btnSave *widget.Button
	btnLoad *widget.Button

	msgDiag  *messageDialog
	mailBox  *mbox.Mbox
	selected int
}

func New(app fyne.App, win fyne.Window) *UI {
	newUI := &UI{
		app: app,
		win: win,
	}

	newUI.build()

	return newUI
}

func (u *UI) Content() fyne.CanvasObject {
	return u.content
}

func (u *UI) build() {
	styleTextMono := fyne.TextStyle{Monospace: true}

	u.msgDiag = NewMessageDialog(u.win)

	selector := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil {
			u.msgDiag.Show("Mailbox Open", err.Error())
			return
		}

		if uri != nil {
			go u.loadMbox(uri.Path(), true)
		}
	}, u.win)

	saver := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			u.msgDiag.Show("File Save", err.Error())
			return
		}

		if writer != nil {
			go u.saveMbox(writer)
		}
	}, u.win)

	title := widget.NewLabel("")
	title.TextStyle = styleTextMono

	u.title = title

	header := container.NewHBox(
		widget.NewIcon(theme.NewThemedResource(theme.QuestionIcon())),
		title,
	)

	u.msgList = widget.NewList(
		func() int {
			if u.mailBox != nil {
				return len(u.mailBox.List)
			} else {
				return 0
			}
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewCheck("", nil),
				widget.NewLabel("template"),
			)
		},
		func(id widget.ListItemID, o fyne.CanvasObject) {
			box := o.(*fyne.Container)
			chk := box.Objects[0].(*widget.Check)
			lbl := box.Objects[1].(*widget.Label)

			msg := u.mailBox.Messages[u.mailBox.List[id]]

			lbl.TextStyle = styleTextMono
			lbl.Alignment = fyne.TextAlignLeading
			lbl.SetText(fmt.Sprintf("%15s%10s     %s", msg.Date, msg.From, msg.Subject))

			chk.OnChanged = nil
			chk.SetChecked(msg.Selected)

			chk.OnChanged = func(v bool) {
				msg.Selected = v
				if v {
					u.selected++
				} else {
					u.selected--
				}

				if u.selected > 0 {
					u.btnSave.Enable()
				} else {
					u.btnSave.Disable()
				}
			}
		},
	)

	btnSave := widget.NewButton("Write Registry", func() {
		home, err := os.UserHomeDir()
		if err == nil && len(home) > 0 {
			uriHome, err := storage.ListerForURI(storage.NewFileURI(home))
			if err == nil {
				saver.SetLocation(uriHome)
			}
		}

		saver.Show()
	})

	u.btnSave = btnSave

	btnLoad := widget.NewButton("Load Mailbox", func() {
		if u.mailBox != nil && len(u.mailBox.Path) > 0 {
			uriList, err := storage.ListerForURI(storage.NewFileURI(u.mailBox.Path))
			if err == nil {
				selector.SetLocation(uriList)
			}
		}

		selector.Show()
	})

	u.btnLoad = btnLoad

	btnQuit := widget.NewButton("Quit", func() { u.app.Quit() })

	actions := container.NewHBox(btnSave, btnLoad, btnQuit)
	bottom := container.NewBorder(nil, nil, nil, actions, nil)

	u.loadMbox(mbox.GetDefaultInbox(), false)

	u.content = container.NewBorder(header, bottom, nil, nil, u.msgList)
}

func (u *UI) loadMbox(path string, selector bool) {
	fyne.Do(func() {
		u.btnSave.Disable()
		u.btnLoad.Disable()
	})

	mailBox, err := mbox.New(path)
	if err != nil && selector {
		u.msgDiag.Show("Error loading mailbox", err.Error())
	}

	u.mailBox = mailBox

	fyne.Do(func() {
		if u.mailBox == nil {
			u.title.SetText("Use Load Mailbox button to load Pat mailbox")
		} else {
			u.title.SetText(fmt.Sprintf("Inbox of %s: %d messages found", u.mailBox.Path, len(u.mailBox.List)))
		}

		u.title.Refresh()

		u.msgList.UnselectAll()
		u.msgList.Refresh()

		u.btnLoad.Enable()
	})
}

func (u *UI) saveMbox(writer fyne.URIWriteCloser) {
	fyne.Do(func() {
		u.btnSave.Disable()
		u.btnLoad.Disable()
	})

	lines, size, err := u.mailBox.WriteRegistry(writer)
	if err != nil {
		u.msgDiag.Show("Error writing registry", err.Error())
	} else {
		info := fmt.Sprintf("%s: %d entries (%d bytes)", writer.URI().Path(), lines, size)
		u.msgDiag.Show("Wrote Registry", info)
	}

	fyne.Do(func() {
		u.btnSave.Enable()
		u.btnLoad.Enable()
	})
}
