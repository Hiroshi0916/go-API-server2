package entities

type Item struct {
	ID    string `gorm:"column:id" json:"id"`
	Value string `gorm:"column:value" json:"value"`
}
