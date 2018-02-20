package amf

import (
	"bytes"
	"errors"
)

const (
	AMF0 = 0x0
	AMF3 = 0x3
)

const (
	AMF0Number byte = iota
	AFM0Boolean
	AMF0String
	AMF0Object
	_
	AMF0Null
	_
	_
	AMF0ECMAArray
	AMF0ObjectEnd
	AFM0StrictArray
	AMF0Date
	AMF0LongString
	_
	_
	AFM0XMLDocument
	AMF0TypedObject
	AMF0SwitchAMF3
)

const (
	AMF3Undefined 	byte = iota
	AFM3Null
	AMF3BooleanFalse
	AMF3BooleanTrue
	AMF3Integer
	AMF3Double
	AMF3String
	AMF3XML
	AMF3Date
	AMF3Object
	AMF3XMLEnd
	AMF3ByteArrray
	AMF3VectorInt
	AFM3VectorUInt
	AMF3VectorDouble
	AMF3VectorObject
	AFM3Dictionary
)

type Object map[string]interface{}

func DecodeAMF(buf []byte, encoding float64) ([]interface{}, error) {
	reader := bytes.NewReader(buf)

	switch encoding {
	case AMF0:
		return decodeAMF0Body(reader)
	case AMF3:
		return decodeAMF3Body(reader)
	default:
		return nil, errors.New("Unsupported AMF encoding")
	}
}

func EncodeAMF(vals []interface{}, encoding float64) ([]byte, error) {
	writer := new(bytes.Buffer)

	switch encoding {
	case AMF0:
		err := encodeAMF0Body(writer, vals)
		if err != nil {
			return nil, err
		}
	case AMF3:
		err := encodeAMF3Body(writer, vals)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("Unsupported AMF encoding")
	}

	return writer.Bytes(), nil
}
