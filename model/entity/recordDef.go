package entity

type RecordDef struct {
	Identifier
	//Id    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Type  string `gorm:"column:type;size:50;not null" json:"type"`
	Value string `gorm:"column:value;size:255;not null" json:"value"`

	PersonId uint `gorm:"column:personId;not null" json:"personId"`
}
