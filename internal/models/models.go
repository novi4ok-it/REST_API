package models

type User struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"-"` // Хранить хэш пароля
}

type TodoList struct {
	ID     int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Title  string `json:"title"`
	UserID int    `json:"-" gorm:"index;foreignKey:UserID"`
	Tasks  []Task `json:"tasks" gorm:"foreignKey:ListID"`
}

type Task struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	ListID      int    `json:"-"`
}
