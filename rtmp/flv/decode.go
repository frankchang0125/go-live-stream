package flv

import (
    bin "github.com/frankchang0125/go-live-stream/binary"
)

func DecodeAudio(data []byte) *AudioBody {
    soundFormat := (data[0] & 0xf0) >> 4
    soundRate := (data[0] & 0xc) >> 2
    soundSize := (data[0] & 0x2) >> 1
    soundType := data[0] & 0x1

    audioTagHeader := AudioTagHeader{
        SoundFormat: soundFormat,
        SoundRate: soundRate,
        SoundSize: soundSize,
        SoundType: soundType,
    }
   
    var audioBody AudioBody

    // If SoundFormat == 10 (UI8)
    if soundFormat == 10 {
        audioTagHeader.AACPacketType = data[1]
        audioData := data[2:]

        audioBody = AudioBody{
            //Tag: flvTag,
            AudioTagHeader: audioTagHeader,
            Data: audioData,
        }
    } else {
        audioData := data[1:]

        audioBody = AudioBody{
            //Tag: flvTag,
            AudioTagHeader: audioTagHeader,
            Data: audioData,
        }
    }

    return &audioBody
}

func DecodeVideo(data []byte) *VideoBody {
    frameType := (data[0] & 0xf0) >> 4
    codecID := data[0] & 0xf

    videoTagHeader := VideoTagHeader{
        FrameType: frameType,
        CodecID: codecID,
    }

    var videoBody VideoBody

    if codecID == 7 {
        videoTagHeader.AVCPacketType = data[1]
        videoTagHeader.CompositionTime = bin.I24BE(data[2:5])
        videoData := data[5:]

        videoBody = VideoBody{
            //Tag: flvTag,
            VideoTagHeader: videoTagHeader,
            Data: videoData,
        }
    } else {
        videoData := data[1:]
        
        videoBody = VideoBody{
            //Tag: flvTag,
            VideoTagHeader: videoTagHeader,
            Data: videoData,
        }
    }

    return &videoBody
}
