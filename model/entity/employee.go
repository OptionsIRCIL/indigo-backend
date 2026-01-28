package entity

type Employee struct {
	Identifier
	//Id        uint   `gorm:"column:id;size:128;" json:"id"`
	FirstName string `gorm:"column:firstName;size:255;" json:"firstName"`
	LastName  string `gorm:"column:lastName;size:255;" json:"lastName"`

	// optional/may remove

	Username string `gorm:"column:username;size:255;" json:"username"`
	Email    string `gorm:"column:email;size:255;" json:"email"`
}
