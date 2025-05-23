package encoders

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

type Metadata map[string]string

func (this Metadata) MarshalBinary() ([]byte, error) {

	if this == nil {
		return nil, nil
	}

	var packString = func(val string, maxBytes int) []byte {

		if maxBytes <= 0 {
			panic("invalid slicing")
		}

		if data := []byte(val); len(data) > maxBytes {

			utfLength := (len(val) * maxBytes) / len(data)
			if utfLength < 1 {
				panic("invalid slicing")
			}

			return []byte(val[:utfLength-1])
		}

		return []byte(val)
	}

	var buff bytes.Buffer

	for key, val := range this {

		keyData := packString(key, math.MaxUint8)
		valData := packString(val, math.MaxUint16)

		if err := buff.WriteByte(byte(len(keyData))); err != nil {
			return nil, err
		}

		valSizeBuff := make([]byte, 2)
		binary.LittleEndian.PutUint16(valSizeBuff, uint16(len(valData)))
		if _, err := buff.Write(valSizeBuff); err != nil {
			return nil, err
		}

		if _, err := buff.Write(keyData); err != nil {
			return nil, err
		}

		if _, err := buff.Write(valData); err != nil {
			return nil, err
		}
	}

	return buff.Bytes(), nil
}

func (this *Metadata) UnmarshalBinary(data []byte) error {

	if len(data) == 0 {
		return nil
	} else if this == nil {
		return errors.New("nil target value")
	} else if *this == nil {
		*this = make(Metadata)
	}

	reader := bytes.NewReader(data)

	for {

		keySizeByte, err := reader.ReadByte()
		if err == io.EOF {
			break
		}

		valSizeBuff := make([]byte, 2)
		if _, err := io.ReadFull(reader, valSizeBuff); err != nil {
			return fmt.Errorf("unable to read value size: %v", err)
		}

		keyBuff := make([]byte, int(keySizeByte))
		if _, err := io.ReadFull(reader, keyBuff); err != nil {
			return fmt.Errorf("unable to read key: %v", err)
		}

		valSize := int(binary.LittleEndian.Uint16(valSizeBuff))
		if valSize > math.MaxUint16 {
			return errors.New("invalid field value size")
		}

		valBuff := make([]byte, valSize)
		if _, err := io.ReadFull(reader, valBuff); err != nil {
			return fmt.Errorf("unable to read value: %v", err)
		}

		(*this)[string(keyBuff)] = string(valBuff)
	}

	return nil
}
