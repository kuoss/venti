package auth

type User struct {
	ID           int    `gorm:"primaryKey"`
	Username     string `gorm:"index:,unique"`
	Hash         string
	IsAdmin      bool
	Token        string
	TokenExpires time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
