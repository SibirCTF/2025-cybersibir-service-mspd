package models

type Suspect struct {
	ID        uint   `gorm:"index"`
	AuthorID  uint   `json:"author_id" gorm:"column:Author_id"`
	SusName   string `json:"sus_name"`
	SusDesc   string `json:"sus_desc"`
	CrimeDesc string `json:"crime_desc"`
	Sbertoken string `json:"sbertoken"`
}
