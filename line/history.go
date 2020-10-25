package line

import "github.com/xenobyter/xbsh/db"

type history struct {
	id      int64
	pending string
}

func newHistory() *history {
	return &history{id: 0}
}

func (h *history) scroll(dy int64) (line string) {
	h.id += dy
	max := db.GetMaxID()
	switch {
	case h.id < 0:
		h.id = max
		line, h.id = db.HistoryRead(h.id)
	case h.id == 0:
		line = h.pending
	case h.id>max:
		line = h.pending
		h.id = 0
	default:
		line, h.id = db.HistoryRead(h.id)
	}
	return
}
