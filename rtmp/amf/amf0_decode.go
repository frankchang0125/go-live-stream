package amf

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	log "github.com/sirupsen/logrus"
)

func decodeAMF0Body(r *bytes.Reader) (results []interface{}, err error) {
	for {
		r, err := decodeAMF0(r)
		if err != nil {
			if err == io.EOF {
				return results, nil
			}

			return nil, err
		}

		results = append(results, r)
	}
}

func decodeAMF0(r *bytes.Reader) (interface{}, error) {
	amf0Type, err := decodeAMF0Type(r)
	if err != nil {
		return nil, err
	}

	switch amf0Type {
	case AMF0Number:
		return decodeAMF0Number(r)
	case AFM0Boolean:
		return decodeAMF0Boolean(r)
	case AMF0String:
		return decodeAMF0String(r)
	case AMF0Object:
		return decodeAMF0Object(r)
	case AMF0ECMAArray:
		return decodeAMF0ECMAArray(r)
	case AMF0Null:
		return nil, nil
	default:
		log.WithField("type", amf0Type).Error("Unsupported AMF0 type.")
		return nil, errors.New("unsupported AMF0 type")
	}
}

func decodeAMF0Type(r *bytes.Reader) (amf0Type uint8, err error) {
	err = binary.Read(r, binary.BigEndian, &amf0Type)
	return
}

func decodeAMF0Number(r *bytes.Reader) (number float64, err error) {
	err = binary.Read(r, binary.BigEndian, &number)
	return
}

func decodeAMF0Boolean(r *bytes.Reader) (bool, error) {
	var boolean uint8
	err := binary.Read(r, binary.BigEndian, &boolean)
	if err != nil {
		return false, err
	}

	return boolean == 1, err
}

func decodeAMF0String(r *bytes.Reader) (string, error) {
	var length uint16
	err := binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return "", err
	}

	if length == 0 {
		return "", nil
	}

	buf := make([]byte, length)
	_, err = r.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func decodeAMF0Object(r *bytes.Reader) (Object, error) {
	obj := Object{}

	for {
		key, err := decodeAMF0String(r)
		if err != nil {
			return nil, err
		}

		if key == "" {
			// Check if we have reached Object End
			var objEnd uint8
			err := binary.Read(r, binary.BigEndian, &objEnd)
			if err != nil {
				return nil, err
			}

			if objEnd != AMF0ObjectEnd {
				// Oops, something goes wrong
				return nil, errors.New("invalid AMF0 object")
			}

			return obj, nil
		}

		value, err := decodeAMF0(r)
		if err != nil {
			return nil, err
		}

		obj[key] = value
	}
}

func decodeAMF0ECMAArray(r *bytes.Reader) (Object, error) {
	var length uint32
	err := binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}

	return decodeAMF0Object(r)
}
