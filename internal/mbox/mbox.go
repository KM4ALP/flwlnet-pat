package mbox

import (
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/la5nta/wl2k-go/fbb"
	"github.com/la5nta/wl2k-go/mailbox"
)

type message struct {
	Date    string
	Size    string
	From    string
	Subject string

	Selected bool
}

type Mbox struct {
	List     []string
	Messages map[string]*message
	Path     string
}

func New(path string) (*Mbox, error) {
	mbox := &Mbox{
		Path: path,
	}

	msgs, err := mailbox.LoadMessageDir(path)
	if err != nil {
		return nil, err
	}

	if len(msgs) < 1 {
		return mbox, nil
	}

	sort.Sort(fbb.ByDate(msgs))

	sorted := make([]string, 0, len(msgs))
	messages := make(map[string]*message, len(msgs))

	for _, msg := range msgs {
		sorted = append(sorted, msg.MID())

		messages[msg.MID()] = &message{
			Date:     msg.Date().Format(fbb.DateLayout),
			Size:     strconv.Itoa(msg.BodySize()),
			From:     strings.ToUpper(msg.From().Addr),
			Subject:  msg.Subject(),
			Selected: false,
		}
	}

	mbox.List = sorted
	mbox.Messages = messages

	return mbox, nil
}

func (m *Mbox) WriteRegistry(out io.WriteCloser) (int, int, error) {
	defer out.Close() //nolint

	lines := 0
	size := 0

	for _, id := range m.List {
		msg := m.Messages[id]
		if !msg.Selected {
			continue
		}

		fields := []string{
			id,
			msg.Date,
			msg.From,
			msg.From,
			"2",
			"False|False|||||",
			msg.Size,
			msg.Subject,
			"GNET",
			"", "", "", "", "", "", "", "",
			"0",
			"00010101000000",
		}

		n, err := io.WriteString(out, strings.Join(fields, "\x01")+"\r\n")
		if err != nil {
			return lines, size, err
		}

		size += n
		lines++
	}

	return lines, size, nil
}
