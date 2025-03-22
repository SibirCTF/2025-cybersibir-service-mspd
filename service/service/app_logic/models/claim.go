package models

type Claim struct {
	ID     uint    `gorm:"index"`
	UserID uint    `json:"user_id" gorm:"column:User_id"`
	SusID  uint    `gorm:"column:sus_id"`
	Sus    Suspect `gorm:"foreignKey:SusID;references:ID" json:"sus"`
	Status string  `json:"status"`
	Reward string  `json:"reward"`
}
