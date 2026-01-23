package entity

type Employee struct {
	Uuid              uint   `gorm:"column:uuid;size:128;" json:"uuid"`
	EmployeeFirstName string `gorm:"column:employeeFirstName;size:255;" json:"employeeFirstName"`
	EmployeeLastName  string `gorm:"column:employeeLastName;size:255;" json:"employeeLastName"`

	// optional/may remove

	EmployeeUsername string `gorm:"column:employeeUsername;size:255;" json:"employeeUsername"`
	EmployeeEmail    string `gorm:"column:employeeEmail;size:255;" json:"employeeEmail"`
}
