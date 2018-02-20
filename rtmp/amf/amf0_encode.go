package amf

import (
    "io"
    "reflect"
    "encoding/binary"
    "errors"
)

func encodeAMF0Body(w io.Writer, vals []interface{}) error {
    for _, val := range vals {
        err := encodeAMF0(w, val)
        if err != nil {
            return err
        }
    }
    
    return nil
}

func encodeAMF0(w io.Writer, val interface{}) error {
    if val == nil {
        // Null
        return encodeAMF0Null(w)
    }

    v := reflect.ValueOf(val)
    
    switch v.Kind() {
    // Number
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        return encodeAMF0Number(w, float64(v.Uint()))
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return encodeAMF0Number(w, float64(v.Int()))
    case reflect.Float32, reflect.Float64:
        return encodeAMF0Number(w, float64(v.Float()))
    // Boolean
    case reflect.Bool:
        return encodeAMF0Boolean(w, v.Bool())
    // String
    case reflect.String:
        return encodeAMF0String(w, v.String(), true)
    // Object
    case reflect.Map:
        return encodeAMF0Object(w, val.(Object))
    default:
        return errors.New("Unsupported AMF type")
    }
}

func encodeAMF0Number(w io.Writer, val float64) error {
    err := binary.Write(w, binary.BigEndian, AMF0Number)
    if err != nil {
        return err
    }
    
    return binary.Write(w, binary.BigEndian, val)
}

func encodeAMF0Boolean(w io.Writer, val bool) error {
    err := binary.Write(w, binary.BigEndian, AFM0Boolean)
    if err != nil {
        return nil
    }
    
    return binary.Write(w, binary.BigEndian, val)
}

func encodeAMF0String(w io.Writer, val string, withMarker bool) error {
    if withMarker {
        err := binary.Write(w, binary.BigEndian, AMF0String)
        if err != nil {
            return err
        }
    }
    
    b := []byte(val)
    
    err := binary.Write(w, binary.BigEndian, uint16(len(b)))
    if err != nil {
        return err
    }
    
    return binary.Write(w, binary.BigEndian, b)
}

func encodeAMF0Object(w io.Writer, obj Object) error {
    err := binary.Write(w, binary.BigEndian, AMF0Object)
    if err != nil {
        return err
    }
    
    for key, val := range obj {
        err = encodeAMF0String(w, key, false)
        if err != nil {
            return err
        }
        
        err = encodeAMF0(w, val)
        if err != nil {
            return err
        }
    }
    
    // Object End: preceded by an empty 16-bit string length)
    binary.Write(w, binary.BigEndian, uint16(0))
    if err != nil {
        return err
    }
    return binary.Write(w, binary.BigEndian, AMF0ObjectEnd)
}

func encodeAMF0Null(w io.Writer) error {
    return binary.Write(w, binary.BigEndian, AMF0Null)
}
