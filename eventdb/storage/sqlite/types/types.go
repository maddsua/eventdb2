package types

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/maddsua/eventdb2/storage/model"
)

type Blob sql.RawBytes
type NullBlob sql.Null[Blob]

func NullBlobSlice(val []byte) NullBlob {
	if len(val) == 0 {
		return NullBlob{}
	}
	return NullBlob{V: val, Valid: true}
}

func NullInt(val int64) sql.NullInt64 {
	return sql.NullInt64{Int64: val, Valid: true}
}

func NullIntPtr(val *int64) sql.NullInt64 {
	if val == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *val, Valid: true}
}

func NullTimePtr(val *time.Time) sql.NullInt64 {
	if val == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: val.UnixNano(), Valid: true}
}

func NullUUID(val uuid.NullUUID) NullBlob {
	if !val.Valid {
		return NullBlob{}
	}
	return NullBlob{V: val.UUID[:], Valid: true}
}

func EncodeStringMap(val model.StringMap) (Blob, error) {

	if val == nil {
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

	for key, val := range val {

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

func DecodeStringMap(data NullBlob) (model.StringMap, error) {

	if !data.Valid || len(data.V) == 0 {
		return nil, nil
	}

	val := model.StringMap{}

	reader := bytes.NewReader(data.V)

	for {

		keySizeByte, err := reader.ReadByte()
		if err == io.EOF {
			break
		}

		valSizeBuff := make([]byte, 2)
		if _, err := io.ReadFull(reader, valSizeBuff); err != nil {
			return val, fmt.Errorf("unable to read value size: %v", err)
		}

		keyBuff := make([]byte, int(keySizeByte))
		if _, err := io.ReadFull(reader, keyBuff); err != nil {
			return val, fmt.Errorf("unable to read key: %v", err)
		}

		valSize := int(binary.LittleEndian.Uint16(valSizeBuff))
		if valSize > math.MaxUint16 {
			return val, errors.New("invalid field value size")
		}

		valBuff := make([]byte, valSize)
		if _, err := io.ReadFull(reader, valBuff); err != nil {
			return val, fmt.Errorf("unable to read value: %v", err)
		}

		val[string(keyBuff)] = string(valBuff)
	}

	return val, nil
}
