package model

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

type Date struct {
	time.Time
}

func (d *Date) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)

	*d = Date{Time: nullTime.Time}
	return err
}

func (d Date) Value() (driver.Value, error) {
	return d.Time, nil
}

func (d Date) GormDataType() string {
	return "date"
}

func (d Date) GobEncode() ([]byte, error) {
	return d.Time.GobEncode()
}

func (d *Date) GobDecode(b []byte) error {
	return d.Time.GobDecode(b)
}

func (d *Date) UnmarshalJSON(b []byte) (err error) {
	d.Time, err = time.Parse(`"`+time.DateOnly+`"`, string(b))
	return err
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(d.Format(`"` + time.DateOnly + `"`)), nil
}
