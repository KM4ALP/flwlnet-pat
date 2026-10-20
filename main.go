package main

import (
	"context"
	"os/signal"
	"syscall"

	"flwlnet-pat/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	a := app.NewWithID("com.github.km4alp.flwlnet-pat")
	w := a.NewWindow("FLWLNET - Pat Registry")

	u := ui.New(a, w)

	w.SetContent(u.Content())

	w.Resize(fyne.NewSize(1024, 768))

	go func() {
		<-ctx.Done()
		w.Close()
	}()

	w.ShowAndRun()
}
