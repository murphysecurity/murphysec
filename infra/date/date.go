package date

import "time"

type Date time.Time

func (d *Date) UnmarshalText(input []byte) error {
	t, e := time.Parse("2006-01-02", string(input))
	if e != nil {
		return e
	}
	*d = Date(t)
	return nil
}
