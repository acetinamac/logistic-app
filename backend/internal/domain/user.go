package domain

type Role string

const (
	RoleClient Role = "client"
	RoleAdmin  Role = "admin"
)

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Email    string `json:"email" gorm:"uniqueIndex;size:255"`
	Password string `json:"-"`
	Role     Role   `json:"role"`
}
