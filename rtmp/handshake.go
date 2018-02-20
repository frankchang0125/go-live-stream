package rtmp

import (
	"crypto/rand"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
)

const (
	c0Len = 1
	c1Len = 1536
	c2Len = 1536
	s0Len = 1
	s1Len = 1536
	s2Len = 1536
)

func (conn *Conn) Handshake() error {
	log.Info("Handshaking...")
	var S0S1S2 [(s0Len + s1Len + s2Len)]byte

	S0 := S0S1S2[:s0Len]
	S1 := S0S1S2[s0Len:(s0Len + s1Len)]
	S2 := S0S1S2[(s0Len + s1Len):]

	// <- C0C1
	C0C1 := make([]byte, c0Len+c1Len)
	_, err := io.ReadFull(conn, C0C1)
	if err != nil {
		log.WithField("err", err).Error("Fail to read C0C1 packet.")
		return err
	}

	C0 := C0C1[:c0Len]
	C1 := C0C1[c0Len:]

	if C0[0] != 3 {
		log.WithField("version", C0[0]).Error("Invalid RTMP version.")
		err = fmt.Errorf("Invalid RTMP version %d", C0[0])
		return err
	}

	// Version: 3
	S0[0] = 3

	S1Time := make([]byte, 4)
	S1Zero := make([]byte, 4)
	S1Random := make([]byte, 1528)

	_, err = rand.Read(S1Random)
	if err != nil {
		log.WithField("err", err).Error("Fail to create S1 random.")
		return err
	}

	copy(S1, append(append(S1Time, S1Zero...), S1Random...))

	/*
	 * Send S0S1S2 immediately after receiving C0C1,
	 * thus the epoch timestamp set to be same as C1
	 */
	S2Time1 := C1[0:4]
	S2Time2 := C1[0:4]
	S2Random := C1[8:c1Len]

	copy(S2, append(append(S2Time1, S2Time2...), S2Random...))

	// S0S1S2 ->
	_, err = conn.Write(S0S1S2[:])
	if err != nil {
		log.WithField("err", err).Error("Fail to handshake.")
		return err
	}

	// <- C2
	C2 := make([]byte, c2Len)
	_, err = io.ReadFull(conn, C2)
	if err != nil {
		log.WithField("err", err).Error("Fail to read C2 packet.")
		return err
	}

	log.Info("Handshake complete.")

	return nil
}
