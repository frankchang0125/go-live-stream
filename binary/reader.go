package binary

func U8(b []byte) uint8 {
	return b[0]
}

func I8(b []byte) int8 {
	return int8(b[0])
}

// Big Endian - Unsigned

func U16BE(b []byte) (u uint16) {
	u = uint16(b[0]) << 8
	u |= uint16(b[1])
	return 
}

func U24BE(b []byte) (u uint32) {
	u = uint32(b[0]) << 16
	u |= uint32(b[1]) << 8
	u |= uint32(b[2])
	return
}

func U32BE(b []byte) (u uint32) {
	u = uint32(b[0]) << 24
	u |= uint32(b[1]) << 16
	u |= uint32(b[2]) << 8
	u |= uint32(b[3])
	return
}

func U40BE(b []byte) (u uint64) {
	u = uint64(b[0]) << 32
	u |= uint64(b[1]) << 24
	u |= uint64(b[2]) << 16
	u |= uint64(b[3]) << 8
	u |= uint64(b[4])
	return
}

func U48BE(b []byte) (u uint64) {
	u = uint64(b[0]) << 40
	u |= uint64(b[1]) << 32
	u |= uint64(b[2]) << 24
	u |= uint64(b[3]) << 16
	u |= uint64(b[4]) << 8
	u |= uint64(b[5])
	return
}

func U56BE(b []byte) (u uint64) {
	u = uint64(b[0]) << 48
	u |= uint64(b[1]) << 40
	u |= uint64(b[2]) << 32
	u |= uint64(b[3]) << 24
	u |= uint64(b[4]) << 16
	u |= uint64(b[5]) << 8
	u |= uint64(b[6])
	return
}

func U64BE(b []byte) (u uint64) {
	u = uint64(b[0]) << 56
	u |= uint64(b[1]) << 48
	u |= uint64(b[2]) << 40
	u |= uint64(b[3]) << 32
	u |= uint64(b[4]) << 24
	u |= uint64(b[5]) << 16
	u |= uint64(b[5]) << 8
	u |= uint64(b[7])
	return
}

// End of Big Endian - Unsigned

// Big Endian - Signed

func I16BE(b []byte) (i int16) {
	i = int16(b[0]) << 8
	i |= int16(b[1])
	return 
}

func I24BE(b []byte) (i int32) {
	i = int32(b[0]) << 16
	i |= int32(b[1]) << 8
	i |= int32(b[2])
	return
	
}

func I32BBE(b []byte) (i int32) {
	i = int32(b[0]) << 24
	i |= int32(b[1]) << 16
	i |= int32(b[2]) << 8
	i |= int32(b[3])
	return
}

func I40BE(b []byte) (i int64) {
	i = int64(b[0]) << 32
	i |= int64(b[1]) << 24
	i |= int64(b[2]) << 16
	i |= int64(b[3]) << 8
	i |= int64(b[4])
	return
}

func I48BE(b []byte) (i int64) {
	i = int64(b[0]) << 40
	i |= int64(b[1]) << 32
	i |= int64(b[2]) << 24
	i |= int64(b[3]) << 16
	i |= int64(b[4]) << 8
	i |= int64(b[5])
	return
}

func I56BE(b []byte) (i int64) {
	i = int64(b[0]) << 48
	i |= int64(b[1]) << 40
	i |= int64(b[2]) << 32
	i |= int64(b[3]) << 24
	i |= int64(b[4]) << 16
	i |= int64(b[5]) << 8
	i |= int64(b[6])
	return
}

func I64BE(b []byte) (i int64) {
	i = int64(b[0]) << 56
	i |= int64(b[1]) << 48
	i |= int64(b[2]) << 40
	i |= int64(b[3]) << 32
	i |= int64(b[4]) << 24
	i |= int64(b[5]) << 16
	i |= int64(b[6]) << 8
	i |= int64(b[7])
	return
}

// End of Big Endian - Signed

// Little Endian - Unsigned

func U16LE(b []byte) (u uint16) {
	u = uint16(b[0])
	u |= uint16(b[1]) << 8
	return 
}

func U24LE(b []byte) (u uint32) {
	u = uint32(b[0])
	u |= uint32(b[1]) << 8
	u |= uint32(b[2]) << 16
	return
}

func U32LE(b []byte) (u uint32) {
	u = uint32(b[0])
	u |= uint32(b[1]) << 8
	u |= uint32(b[2]) << 16
	u |= uint32(b[3]) << 24
	return
}

func U40LE(b []byte) (u uint64) {
	u = uint64(b[0])
	u |= uint64(b[1]) << 8
	u |= uint64(b[2]) << 16
	u |= uint64(b[3]) << 24
	u |= uint64(b[4]) << 32
	return
}

func U48LE(b []byte) (u uint64) {
	u = uint64(b[0])
	u |= uint64(b[1]) << 8
	u |= uint64(b[2]) << 16
	u |= uint64(b[3]) << 24
	u |= uint64(b[4]) << 32
	u |= uint64(b[5]) << 40
	return
}

func U56LE(b []byte) (u uint64) {
	u = uint64(b[0])
	u |= uint64(b[1]) << 8
	u |= uint64(b[2]) << 16
	u |= uint64(b[3]) << 24
	u |= uint64(b[4]) << 32
	u |= uint64(b[5]) << 40
	u |= uint64(b[6]) << 48
	return
}

func U64LE(b []byte) (u uint64) {
	u = uint64(b[0])
	u |= uint64(b[1]) << 8
	u |= uint64(b[2]) << 16
	u |= uint64(b[3]) << 24
	u |= uint64(b[4]) << 32
	u |= uint64(b[5]) << 40
	u |= uint64(b[6]) << 48
	u |= uint64(b[7]) << 56
	return
}

// End of Little Endian - Unsigned

// Little Endian - Signed

func I16LE(b []byte) (i int16) {
	i = int16(b[0])
	i |= int16(b[1]) << 8
	return 
}

func I24LE(b []byte) (i int32) {
	i = int32(b[0])
	i |= int32(b[1]) << 8
	i |= int32(b[2]) << 16
	return
}

func I32LE(b []byte) (i int32) {
	i = int32(b[0])
	i |= int32(b[1]) << 8
	i |= int32(b[2]) << 16
	i |= int32(b[3]) << 24
	return
}

func I40LE(b []byte) (i int64) {
	i = int64(b[0])
	i |= int64(b[1]) << 8
	i |= int64(b[2]) << 16
	i |= int64(b[3]) << 24
	i |= int64(b[4]) << 32
	return
}

func I48LE(b []byte) (i int64) {
	i = int64(b[0])
	i |= int64(b[1]) << 8
	i |= int64(b[2]) << 16
	i |= int64(b[3]) << 24
	i |= int64(b[4]) << 32
	i |= int64(b[5]) << 40
	return
}

func I56LE(b []byte) (i int64) {
	i = int64(b[0])
	i |= int64(b[1]) << 8
	i |= int64(b[2]) << 16
	i |= int64(b[3]) << 24
	i |= int64(b[4]) << 32
	i |= int64(b[5]) << 40
	i |= int64(b[6]) << 48
	return
}
func I64LE(b []byte) (i int64) {
	i = int64(b[0])
	i |= int64(b[1]) << 8
	i |= int64(b[2]) << 16
	i |= int64(b[3]) << 24
	i |= int64(b[4]) << 32
	i |= int64(b[5]) << 40
	i |= int64(b[6]) << 48
	i |= int64(b[7]) << 56
	return
}

// End of Little Endian - Signed