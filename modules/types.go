package modules

type Session struct {
	UserID        string `gorm:"primaryKey"`
	RState        bool
	CurrentImages int
}
