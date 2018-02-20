package rtmp

const (
    typeVideo = iota
    typeAudio
)

type Packet struct {
    packetType  int
    timestamp   uint32
    streamID    uint32
    data        []byte
}

func NewPacket(packetType int, timestamp uint32, streamID uint32, data []byte) *Packet {
    return &Packet{
        packetType: packetType,
        timestamp:  timestamp,
        streamID:   streamID,
        data:       data,
    }
}
