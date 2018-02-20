package amf

import (
	"bytes"
)

func decodeAMF3Body(r *bytes.Reader) (results []interface{}, err error) {
	for {
		r, err := decodeAMF3(r)
		if err != nil {
			return nil, err
		}
		
		results = append(results, r)
	}
}

func decodeAMF3(r *bytes.Reader) (interface{}, error) {
	// TODO:
	return nil, nil
}