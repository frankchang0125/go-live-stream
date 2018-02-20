package rtmp

import (
	"io"
	"encoding/binary"
	"errors"

	"github.com/frankchang0125/go-live-stream/rtmp/amf"
	bin "github.com/frankchang0125/go-live-stream/binary"
	log "github.com/sirupsen/logrus"
)

// Message Type ID
const (
	// Protocol Control Messages
	typeIDSetChunkSize     = 1
	typeIDAbortMsg         = 2
	typeIDAck              = 3
	typeIDWindowAckSize    = 5
	typeIDSetPeerBandwidth = 6

	// User Control Messages
	typeIDusrCtrlMsg = 4

	// RTMP Command Messages
	typeIDAudioMsg            = 8
	typeIDVideoMsg            = 9
	typeIDCmdMsgAMF0          = 20
	typeIDCmdMsgAMF3          = 17
	typeIDDataMsgAMF0         = 18
	typeIDDataMsgAMF3         = 15
	typeISSharedObjectMsgAMF0 = 19
	typeIDSharedObjectMsgAMF3 = 16
	typeIDAggregateMsg        = 22
)

type ChunkStreamStatus struct {
	chunk			*Chunk // Latest received/sent chunk on the chunk stream
	timestampDelta	uint32 // Timestamp delta between previous received/sent chunk on the chunk stream
}

type ChunkStream struct {
	conn     				*Conn
	curRead  				map[uint32]*ChunkStreamStatus // Latest status of receiver side chunk stream
	curWrite 				map[uint32]*ChunkStreamStatus // Latest status of sender side chunk stream
}

func NewChunkStream(conn *Conn) *ChunkStream {
	return &ChunkStream{
		conn: 		conn,
		curRead:	make(map[uint32]*ChunkStreamStatus),
		curWrite:	make(map[uint32]*ChunkStreamStatus),
	}
}

func (cs *ChunkStream) readChunk() error {
	chunk := &Chunk{}

	buf, err := cs.readBytes(1)
	if err != nil {
		if err == io.EOF {
			return err
		}

		log.WithField("err", err).Error("Error while reading chunks.")
		return err
	}

	header := bin.U8(buf)
	format := uint32(header >> 6)
	csid := uint32(header) & 0x3F

	switch csid {
	case 0:
		// Chunk basic header 2, read 1 more byte to compose CSID
		buf, err = cs.readBytes(1)
		if err != nil {
			log.WithField("err", err).Error("Error while reading chunk basic header.")
			return err
		}

		csid = uint32(bin.U8(buf))
	case 1:
		// Chunk basic header 3, read 2 more bytes to compose CSID
		buf, err = cs.readBytes(2)
		if err != nil {
			log.WithField("err", err).Error("Error while reading chunk basic header.")
			return err
		}

		csid = 1<<8 | uint32(bin.U8(buf))
	}

	chunk.CSID = csid

	// Format type
	switch format {
	case 0:
		buf, err = cs.readBytes(11)
		if err != nil {
			log.WithField("err", err).Error("Error while reading chunk message header.")
			return err
		}

		timestamp := bin.U24BE(buf[:3])

		if timestamp >= 0xFFFFFF {
			// Read extended timestamp
			var extTimestamp uint32
			err = binary.Read(cs.conn, binary.BigEndian, &extTimestamp)
			if err != nil {
				log.WithField("err", err).Error("Error while reading chunk message header.")
				return err
			}

			chunk.Timestamp = extTimestamp
		} else {
			chunk.Timestamp = timestamp
		}

		chunk.Length = bin.U24BE(buf[3:6])
		chunk.TypeID = uint32(bin.U8(buf[6:7]))
		chunk.StreamID = bin.U32LE(buf[7:11])
	case 1:
		buf, err = cs.readBytes(7)
		if err != nil {
			log.WithField("err", err).Error("Error while reading chunk message header.")
			return err
		}

		timestampDelta := bin.U24BE(buf[:3])

		if cur, ok := cs.curRead[chunk.CSID]; ok {
			if timestampDelta >= 0xFFFFFF {
				// Read extended timestamp
				var extTimestamp uint32
				err = binary.Read(cs.conn, binary.BigEndian, &extTimestamp)
				if err != nil {
					log.WithField("err", err).Error("Error while reading chunk message header.")
					return err
				}

				chunk.Timestamp = cur.chunk.Timestamp + extTimestamp
			} else {
				chunk.Timestamp = cur.chunk.Timestamp + timestampDelta
			}

			// Copy chunk message header information from previous chunk
			chunk.StreamID = cur.chunk.StreamID
		} else {
			err = errors.New("Calculating received chunk's timestamp by timestamp delta, " +
				"but cannot find previous chunk's timestamp.")
			log.WithField("err", err).Error("Error while reading chunk message header.")
			return err
		}

		chunk.Length = bin.U24BE(buf[3:6])
		chunk.TypeID = uint32(bin.U8(buf[6:7]))
	case 2:
		buf, err = cs.readBytes(3)
		if err != nil {
			log.WithField("err", err).Error("Error while reading chunk message header.")
			return err
		}

		timestampDelta := bin.U24BE(buf)

		if cur, ok := cs.curRead[chunk.CSID]; ok {
			chunk.Timestamp = cur.chunk.Timestamp + timestampDelta

			// Copy chunk message header information from previous chunk
			chunk.Length = cur.chunk.Length
			chunk.TypeID = cur.chunk.TypeID
			chunk.StreamID = cur.chunk.StreamID
		} else {
			err = errors.New("Calculating received chunk's timestamp by timestamp delta, " +
				"but cannot find previous chunk's timestamp")
			log.WithField("err", err).Error("Error while reading chunk message header.")
			return err
		}
	case 3:
		if cur, ok := cs.curRead[chunk.CSID]; ok {
			// Copy chunk message header information from previous chunk
			chunk.Timestamp = cur.chunk.Timestamp
			chunk.Length = cur.chunk.Length
			chunk.TypeID = cur.chunk.TypeID
			chunk.StreamID = cur.chunk.StreamID
		} else {
			err = errors.New("Reading type 3 Chunk Message Header, but cannot find correspond previous chunk")
			log.WithField("err", err).Error("Error while reading chunk message header.")
			return err
		}
	}

	// 1. Save the timestamp delta between current chunk and previous chunk on the chunk stream
	// 2. Update the latest received chunk on chunk stream
	if cur, ok := cs.curRead[chunk.CSID]; ok {
		cur.chunk = chunk
		cur.timestampDelta = chunk.Timestamp - cur.chunk.Timestamp
	} else {
		cs.curRead[chunk.CSID] = &ChunkStreamStatus{
			chunk:	chunk,
			timestampDelta: 0,
		}
	}

	// FIXME: Remove me
	log.WithFields(log.Fields{
		"timestamp": chunk.Timestamp,
		"typeID": chunk.TypeID,
		"length": chunk.Length,
		}).Debug("readChunk()")

	switch chunk.TypeID {
	case typeIDSetChunkSize:
		err = cs.conn.setChunkSize()
	case typeIDWindowAckSize:
		err = cs.conn.setWindowAckSize()
	case typeIDusrCtrlMsg:
		err = cs.conn.handleUsrCtrlMsg(cs, chunk)
	case typeIDAudioMsg:
		err = cs.conn.handleAudioMsg(cs, chunk)
	case typeIDVideoMsg:
		err = cs.conn.handleVideoMsg(cs, chunk)
	case typeIDCmdMsgAMF0:
		if cs.conn.amfEncoding == amf.AMF0 {
			err = cs.conn.handleCmdMsg(cs, chunk)
		} else {
			log.Warning("Cannot process AMF3 command under AMF0 encoding.")
		}
	case typeIDCmdMsgAMF3:
		if cs.conn.amfEncoding == amf.AMF3 {
			err = cs.conn.handleCmdMsg(cs, chunk)
		} else {
			log.Warning("Cannot process AMF0 command under AMF3 encoding.")
		}
	case typeIDDataMsgAMF0:
		if cs.conn.amfEncoding == amf.AMF0 {
			err = cs.conn.handleDataMsg(cs, chunk)
		} else {
			log.Warning("Cannot process AMF3 command under AMF0 encoding.")
		}
	case typeIDDataMsgAMF3:
		if cs.conn.amfEncoding == amf.AMF3 {
			err = cs.conn.handleDataMsg(cs, chunk)
		} else {
			log.Warning("Cannot process AMF0 command under AMF3 encoding.")
		}
	default:
		log.WithField("Type ID", chunk.TypeID).Warning("Unknown Message Type ID.")
	}

	if err != nil {
		return err
	}

	return nil
}

func (cs *ChunkStream) writeChunk(chunk *Chunk, chunkSize uint32) error {
	numOfChunks := uint32(chunk.Length / chunkSize)

	if chunk.Length % chunkSize != 0 {
		numOfChunks++
	}

	remain := chunk.Length
	var i uint32

	for i = 0; i < numOfChunks; i++ {
		chunkHeader, err := cs.createChunkHeader(chunk)
		if err != nil {
			return err
		}

		start := i * chunkSize
		var end, size uint32

		if remain < chunkSize {
			size = remain
		} else {
			size = chunkSize
		}

		end = start + size

		c := append(chunkHeader, chunk.Data[start:end]...)

		_, err = cs.conn.Write(c)
		if err != nil {
			log.WithField("err", err).Error("Fail to write chunk.")
			return err
		}

		remain -= size

		// 1. Save the timestamp delta between current chunk and previous chunk
		// 2. Update the latest sent chunk on chunk stream
		if cur, ok := cs.curWrite[chunk.CSID]; ok {
			cur.chunk = chunk
			cur.timestampDelta = chunk.Timestamp - cur.chunk.Timestamp
		} else {
			cs.curWrite[chunk.CSID] = &ChunkStreamStatus{
				chunk: chunk,
				timestampDelta: 0,
			}
		}
	}

	return nil
}

func (cs *ChunkStream) createChunkHeader(chunk *Chunk) (header []byte, err error) {
	var timestamp, timestampDelta, extTimestamp uint32
	var headerType int

	if cur, ok := cs.curWrite[chunk.CSID]; ok {
		timestampDelta = chunk.Timestamp - cur.chunk.Timestamp

		if chunk.StreamID == cur.chunk.StreamID {
			if (chunk.Length == cur.chunk.Length) && (chunk.TypeID == cur.chunk.TypeID) {
				if timestampDelta == cur.timestampDelta {
					// Type 3
					headerType = 3
				} else {
					// Type 2
					headerType = 2
				}
			} else {
				// Type 1
				headerType = 1
			}
		} else {
			// Type 0
			headerType = 0
		}
	} else {
		// Type 0
		headerType = 0
	}

	switch headerType {
	case 0:
		if chunk.Timestamp > 0xFFFFFF {
			extTimestamp = chunk.Timestamp
			timestamp = 0xFFFFFFF
		} else {
			timestamp = chunk.Timestamp
		}
	case 1:
		fallthrough
	case 2:
		if timestampDelta > 0xFFFFFF {
			extTimestamp = timestampDelta
			timestampDelta = 0xFFFFFFF
		}
	}

	var basicHeader []byte

	if chunk.CSID < 64 {
		// Chunk basic header 1
		basicHeader = make([]byte, 1)
		basicHeader[0] = uint8(chunk.CSID)
		basicHeader[0] |= uint8(headerType << 6)
	} else if chunk.CSID < 320 {
		// Chunk basic header 2
		basicHeader = make([]byte, 2)
		basicHeader[0] = uint8(headerType << 6)
		basicHeader[1] = uint8(chunk.CSID - 64)
	} else if chunk.CSID < 65560 {
		// Chunk basic header 3
		basicHeader = make([]byte, 3)
		basicHeader[0] = uint8(headerType<<6) | 0x3F
		basicHeader[1] = uint8((chunk.CSID - 64) >> 8)
		basicHeader[2] = uint8(chunk.CSID - 64)
	} else {
		// Invalidate chunk stream ID
		log.WithField("CSID", chunk.CSID).Error("Invalidate chunk stream ID.")
		return nil, errors.New("invalidate chunk stream ID")
	}

	var messageHeader []byte

	switch headerType {
	case 0:
		messageHeader = make([]byte, 11)
		err := bin.PutUBE(messageHeader[:3], timestamp)
		if err != nil {
			log.WithField("timestamp", timestamp).Error("Invalidate chunk timestamp")
			return nil, err
		}

		err = bin.PutUBE(messageHeader[3:6], chunk.Length)
		if err != nil {
			log.WithField("length", chunk.Length).Error("Invalidate chunk length.")
			return nil, err
		}

		messageHeader[6] = uint8(chunk.TypeID)

		err = bin.PutULE(messageHeader[7:], chunk.StreamID)
		if err != nil {
			log.WithField("Stream ID", chunk.StreamID).Error("Invalidate chunk Stream ID.")
			return nil, err
		}
	case 1:
		messageHeader = make([]byte, 7)
		err := bin.PutUBE(messageHeader[:3], timestampDelta)
		if err != nil {
			log.WithField("timestampDelta", timestampDelta).Error("Invalidate chunk timestamp delta.")
		}

		err = bin.PutUBE(messageHeader[3:6], chunk.Length)
		if err != nil {
			log.WithField("length", chunk.Length).Error("Invalidate chunk length.")
			return nil, err
		}

		messageHeader[6] = uint8(chunk.TypeID)
	case 2:
		messageHeader = make([]byte, 3)

		err := bin.PutUBE(messageHeader[:3], timestampDelta)
		if err != nil {
			log.WithField("timestampDelta", timestampDelta).Error("Invalidate chunk timestamp delta.")
		}
	case 3:
		messageHeader = make([]byte, 0)
	}

	var extTimestampHeader []byte

	if extTimestamp > 0 {
		extTimestampHeader = make([]byte, 4)
	} else {
		extTimestampHeader = make([]byte, 0)
	}

	return append(append(basicHeader, messageHeader...), extTimestampHeader...), nil
}

func (cs *ChunkStream) readAMFBody(length uint32, clientChunkSize uint32) ([]byte, error) {
	var results = make([]byte, 0)
	var n uint32
	remain := length

	for {
		if remain > clientChunkSize {
			n = clientChunkSize
		} else {
			n = remain
		}

		buf, err := cs.readBytes(n)
		if err != nil {
			return nil, nil
		}

		results = append(results, buf...)

		remain -= n
		if remain == 0 {
			return results, err
		}

		// AMF body exceeds client's chunk size and is divided into
		// serveral chunks with 1-byte chunk header ahead
		_, err = cs.readBytes(1)
		if err != nil {
			return nil, err
		}
	}
}

func (cs *ChunkStream) readBytes(n uint32) ([]byte, error) {
	buf := make([]byte, n)
	_, err := cs.conn.Read(buf)
	return buf, err
}
