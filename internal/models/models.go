package models

// User model
// swagger:model
type User struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement"`
	// Username for login
	// required: true
	// example: john_doe
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"-"` // Хранить хэш пароля
}

// TodoList model
// swagger:model
type TodoList struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement"`
	// Title of the list
	// required: true
	// example: Shopping List
	Title  string `json:"title"`
	UserID int    `json:"-" gorm:"index;foreignKey:UserID"`
	Tasks  []Task `json:"tasks" gorm:"foreignKey:ListID"`
}

// Task model
// swagger:model
type Task struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement"`
	// Title of the task
	// required: true
	// example: Buy milk
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	ListID      int    `json:"-"`
}
