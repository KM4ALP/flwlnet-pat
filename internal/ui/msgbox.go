package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type messageDialog struct {
	title  *widget.Label
	label  *widget.Label
	button *widget.Button

	box *widget.PopUp
}

func NewMessageDialog(w fyne.Window) *messageDialog {
	title := widget.NewLabel("")
	title.Wrapping = fyne.TextWrapWord
	title.Alignment = fyne.TextAlignLeading
	title.TextStyle = fyne.TextStyle{Monospace: true}

	label := widget.NewLabel("")
	label.Wrapping = fyne.TextWrapWord
	label.Alignment = fyne.TextAlignLeading
	label.TextStyle = fyne.TextStyle{Monospace: true}

	button := widget.NewButton("Ok", nil)
	button.Importance = widget.HighImportance

	header := container.NewHBox(
		widget.NewIcon(theme.NewErrorThemedResource(theme.InfoIcon())),
		title,
	)

	inner := container.NewVBox(header, label, button)

	content := container.NewPadded(container.NewGridWrap(fyne.NewSize(360, 0), inner))

	box := widget.NewModalPopUp(content, w.Canvas())

	button.OnTapped = func() {
		fyne.Do(func() {
			box.Hide()
		})
	}

	newDiag := &messageDialog{
		title:  title,
		label:  label,
		button: button,
		box:    box,
	}

	return newDiag
}

func (m *messageDialog) Show(title string, msg string) {
	fyne.Do(func() {
		m.title.SetText(title)
		m.title.Refresh()
		m.label.SetText(msg)
		m.label.Refresh()
		m.box.Show()
	})
}
