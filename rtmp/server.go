package rtmp

import (
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	newStreamer chan *Conn
	newViewer   chan *Conn
	channels    sync.Map // Map<Stream Name>*Channel
}

func NewRTMPServer() *Server {
	return &Server{
		newStreamer: make(chan *Conn),
		newViewer:   make(chan *Conn),
	}
}

func (s *Server) HandleRTMPRequest(netConn net.Conn) {
	conn := NewConn(netConn, s.newStreamer, s.newViewer)
	defer func() {
		if conn.info != nil {
			streamName := conn.info.Name
			isPublisher := conn.isPublisher

			if ch, ok := s.channels.Load(streamName); ok {
				// Remove connection from channel list
				channel := ch.(*Channel)

				if isPublisher {
					channel.lock.Lock()
					channel.streamer = nil
					channel.lock.Unlock()
				} else {
					if result := channel.removeViewer(conn); !result {
						log.Warn("Cannot find connection in channel list.")
					}
				}

				// Remove channel if there're no streamer and viewers
				channel.lock.RLock()
				if channel.streamer == nil && len(channel.viewers) == 0 {
					log.WithField("streamName", streamName).Info("Channel removed.")
					s.channels.Delete(conn.info.Name)
				}
				channel.lock.RUnlock()
			}
		}

		log.Info("Connection closed.")
		conn.Close()
		conn = nil
	}()

	err := conn.Handshake()
	if err != nil {
		return
	}

	conn.Serve()
}

func (s *Server) Monitor() {
	for {
		select {
		case conn := <-s.newStreamer:
			streamName := conn.info.Name

			if ch, ok := s.channels.Load(streamName); !ok {
				// Channel not exists, create a new channel
				newChannel := NewChannel(streamName)
				newChannel.streamer = conn
				s.channels.Store(conn.info.Name, newChannel)
				conn.channel = newChannel
				log.WithField("streamName", streamName).Info("New streamer connected.")
			} else {
				channel := ch.(*Channel)

				if channel.streamer != nil {
					// Channel already existed, kick off the existed streamer and replace with new streamer
					log.WithField("streamName",
						streamName).Info("Duplicate streamer detected, disconnect existing streamer.")
					channel.streamer.Close()

					channel.lock.Lock()
					channel.streamer = conn
					conn.channel = channel
					channel.lock.Unlock()

					log.WithField("streamName", streamName).Info("New streamer connected.")
				} else {
					// Channel already existed, which was created by pending viewers
					channel.lock.Lock()
					channel.streamer = conn
					conn.channel = channel
					channel.lock.Unlock()

					log.WithField("streamName", streamName).Info("New streamer connected.")
				}
			}

			conn.channelCreated <- true
		case conn := <-s.newViewer:
			streamName := conn.info.Name

			if ch, ok := s.channels.Load(streamName); !ok {
				// Channel not exists, create a new channel
				newChannel := NewChannel(streamName)
				newChannel.viewers = append(newChannel.viewers, conn)
				s.channels.Store(streamName, newChannel)
				conn.channel = newChannel

				log.WithField("streamName", streamName).Info("New player connected.")
			} else {
				// Channel already existed, add connection to viewers list
				channel := ch.(*Channel)

				channel.lock.Lock()
				channel.viewers = append(channel.viewers, conn)
				channel.lock.Unlock()

				log.WithField("streamName", streamName).Info("New player connected.")
			}

			conn.channelCreated <- true
		}
	}
}
