package model

import "encoding/json"

type IsOnline struct {
	Valid bool
	Value bool
}

func (i *IsOnline) UnmarshalJSON(bytes []byte) error {
	i.Valid = true
	return json.Unmarshal(bytes, &i.Value)
}

func (i *IsOnline) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("true"), nil
	}
	return json.Marshal(i.Value)
}

func (i *IsOnline) SetOnline(b bool) {
	i.Valid = true
	i.Value = b
}

var _ json.Marshaler = (*IsOnline)(nil)
var _ json.Unmarshaler = (*IsOnline)(nil)

func IsOnlineTrue() IsOnline {
	return IsOnline{
		Value: true,
		Valid: true,
	}
}

func IsOnlineFalse() IsOnline {
	return IsOnline{
		Value: false,
		Valid: true,
	}
}
