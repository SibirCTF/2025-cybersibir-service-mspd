package models

type User struct {
	ID       uint      `gorm:"index"`
	Username string    `gorm:"unique" json:"username"`
	Password string    `json:"password"`
	Claims   []Claim   `gorm:"foreignKey:UserID" json:"claims"`
	Suspects []Suspect `gorm:"foreignKey:AuthorID" json:"suspects"`
}
