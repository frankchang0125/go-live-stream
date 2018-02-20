package rtmp

import (
	"github.com/frankchang0125/go-live-stream/rtmp/amf"
)

type Chunk struct {
	CSID      uint32
	Timestamp uint32
	Length    uint32
	TypeID    uint32
	StreamID  uint32 // Message stream ID
	Data      []byte
}

func NewPCMChunk(typeID uint32, data []byte) *Chunk {
	return &Chunk{
		CSID:     2,
		Length:   uint32(len(data)),
		TypeID:   typeID,
		StreamID: 0,
		Data:     data,
	}
}

func NewUCMChunk(data []byte) *Chunk {
	return &Chunk{
		CSID:     2,
		Length:   uint32(len(data)),
		TypeID:   typeIDusrCtrlMsg,
		StreamID: 0,
		Data:     data,
	}
}

func NewAMFCmdChunk(encoding float64, csid uint32, streamID uint32, data []byte) *Chunk {
	var typeID uint32

	if encoding == amf.AMF0 {
		typeID = typeIDCmdMsgAMF0
	} else {
		typeID = typeIDCmdMsgAMF3
	}

	return &Chunk{
		CSID:     csid,
		Length:   uint32(len(data)),
		TypeID:   typeID,
		StreamID: streamID,
		Data:     data,
	}
}

func NewAudioChunk(timestamp uint32, streamID uint32, data []byte) *Chunk {
	return &Chunk{
		CSID:      4,
		Length:    uint32(len(data)),
		TypeID:    typeIDAudioMsg,
		Timestamp: timestamp,
		Data:      data,
	}
}

func NewVideoChunk(timestamp uint32, streamID uint32, data []byte) *Chunk {
	return &Chunk{
		CSID:      6,
		Length:    uint32(len(data)),
		TypeID:    typeIDVideoMsg,
		Timestamp: timestamp,
		Data:      data,
	}
}
