package flv

import (
    bin "github.com/frankchang0125/go-live-stream/binary"
)

type AudioTagHeader struct {
    SoundFormat     uint8
    SoundRate       uint8
    SoundSize       uint8
    SoundType       uint8
    AACPacketType   uint8
}

func (header *AudioTagHeader) Encode() []byte {
    var result []byte

    if header.SoundFormat == 10 {
        result = make([]byte, 2)
        result[0] = (header.SoundFormat << 4) |
                (header.SoundRate << 2) |
                (header.SoundSize << 1) |
                header.SoundType
        result[1] = header.AACPacketType
    } else {
        result = make([]byte, 1)
        result[0] = (header.SoundFormat << 4) |
                (header.SoundRate << 2) |
                (header.SoundSize << 1) |
                header.SoundType
    }

    return result
}

type AudioBody struct {
    AudioTagHeader
    Data []byte
}

type VideoTagHeader struct {
    FrameType           uint8
    CodecID             uint8
    AVCPacketType       uint8
    CompositionTime     int32
}

type VideoBody struct {
    VideoTagHeader
    Data []byte
}

func (header *VideoTagHeader) Encode() []byte {
    var result []byte

    if header.CodecID == 7 {
        result = make([]byte, 5)
        result[0] = header.FrameType << 4 |
                    header.CodecID
        result[1] = header.AVCPacketType
        bin.PutI24BE(result[2:], header.CompositionTime)
    } else {
        result = make([]byte, 1)
        result[0] = header.FrameType << 4 |
                    header.CodecID
    }

    return result
}