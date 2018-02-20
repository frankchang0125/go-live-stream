package rtmp

import (
	"sync"
)

type Channel struct {
	lock     *sync.RWMutex
	name     string
	streamer *Conn
	viewers  []*Conn
}

func NewChannel(name string) *Channel {
	return &Channel{
		lock:    new(sync.RWMutex),
		name:    name,
		viewers: make([]*Conn, 0),
	}
}

func (ch *Channel) removeViewer(conn *Conn) bool {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	for i, viewer := range ch.viewers {
		if viewer == conn {
			// Remove connection from viewers list
			ch.viewers = append(ch.viewers[:i], ch.viewers[i+1:]...)
			return true
		}
	}

	return false
}
