# Go Live
Simple RTMP-based live stream server written in Go.

#### Supported protocols

- [x] RTMP
- [ ] HLS

#### Supported containers

- [x] FLV

## Install

1. Download source codes: `git clone https://github.com/frankchang0125/go-live-stream.git`
2. `cd go-live-stream`
3. `go build`

## Use

Execute `go-live-stream`

## Publish stream

### FFmpeg

`ffmpeg -re -i <Your FLV file> -c copy -f flv rtmp://localhost:1935/golive/mylive`

### OBS

1. Open preferences
2. Switch to `Stream` tab
3. Change URL to: `rtmp://localhost:1935/golive`
4. Change Stream key to: `mylive`

![OBS settings](https://i.imgur.com/RqiN2sD.png)

## Receive stream

### VLC

Open Network URL: `rtmp://localhost:1935/golive/mylive`

### FFplay

`ffplay rtmp://localhost:1935/golive/mylive`

## Known issues

1. Second and later stream receivers on the same published stream cannot play the stream correctly.
2. When packet buffer (size: **1024** video/audio packets) of a published stream has fulled and starting dropping packets, the later connected stream recevier cannot play the stream correctly.
