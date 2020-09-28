package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/volatiletech/null/v8/convert"
)

// NullBytes is a global byte slice of JSON null
var NullBytes = []byte("null")

// Bytes is a nullable []byte.
type Bytes struct {
	Bytes []byte
	Valid bool
	set   bool
}

// NewBytes creates a new Bytes
func NewBytes(b []byte, valid, set bool) Bytes {
	return Bytes{
		Bytes: b,
		Valid: valid,
		set:   set,
	}
}

// BytesFrom creates a new Bytes that will be invalid if nil.
func BytesFrom(b []byte) Bytes {
	return NewBytes(b, b != nil, true)
}

// BytesFromPtr creates a new Bytes that will be invalid if nil.
func BytesFromPtr(b *[]byte) Bytes {
	if b == nil {
		return NewBytes(nil, false, true)
	}
	n := NewBytes(*b, true, true)
	return n
}

func (b Bytes) IsSet() bool {
	return b.set
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Bytes) UnmarshalJSON(data []byte) error {
	b.set = true

	if bytes.Equal(data, NullBytes) {
		b.Valid = false
		b.Bytes = nil
		return nil
	}

	var bv []byte
	if err := json.Unmarshal(data, &bv); err != nil {
		return err
	}

	b.Bytes = bv
	b.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *Bytes) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		b.Bytes = nil
		b.Valid = false
	} else {
		b.Bytes = append(b.Bytes[0:0], text...)
		b.Valid = true
	}

	return nil
}

// MarshalJSON implements json.Marshaler.
func (b Bytes) MarshalJSON() ([]byte, error) {
	if len(b.Bytes) == 0 {
		return NullBytes, nil
	}
	return json.Marshal(b.Bytes)
}

// MarshalText implements encoding.TextMarshaler.
func (b Bytes) MarshalText() ([]byte, error) {
	if !b.Valid {
		return nil, nil
	}
	return b.Bytes, nil
}

// SetValid changes this Bytes's value and also sets it to be non-null.
func (b *Bytes) SetValid(n []byte) {
	b.Bytes = n
	b.Valid = true
}

// Ptr returns a pointer to this Bytes's value, or a nil pointer if this Bytes is null.
func (b Bytes) Ptr() *[]byte {
	if !b.Valid {
		return nil
	}
	return &b.Bytes
}

// IsZero returns true for null or zero Bytes's, for future omitempty support (Go 1.4?)
func (b Bytes) IsZero() bool {
	return !b.Valid
}

// Scan implements the Scanner interface.
func (b *Bytes) Scan(value interface{}) error {
	if value == nil {
		b.Bytes, b.Valid = []byte{}, false
		return nil
	}
	b.Valid = true
	return convert.ConvertAssign(&b.Bytes, value)
}

// Value implements the driver Valuer interface.
func (b Bytes) Value() (driver.Value, error) {
	if !b.Valid {
		return nil, nil
	}
	return b.Bytes, nil
}

// Randomize for sqlboiler
func (b *Bytes) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	if shouldBeNull {
		b.Bytes = nil
		b.Valid = false
	} else {
		b.Bytes = []byte{byte(nextInt() % 256)}
		b.Valid = true
	}
}
