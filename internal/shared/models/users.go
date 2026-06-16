package models

import "time"

type User struct {
	ID        int64      `bun:"id,pk,autoincrement" json:"id"`
	Status    string     `bun:"status,default:'A'" json:"status"`
	UUID      string     `bun:"uuid,notnull,unique,type:uuid,default:gen_random_uuid()" json:"uuid"`
	CreatedAt time.Time  `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time  `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"deleted_at"`

	FirstName  string `bun:"first_name" json:"first_name"`
	LastName   string `bun:"last_name" json:"last_name"`
	Email      string `bun:"email" json:"email"`
	Username   string `bun:"username" json:"username"`
	Password   string `bun:"password" json:"password"`
	UserType   string `bun:"user_type" json:"user_type"`
	Addr1      string `bun:"addr1" json:"addr1"`
	Addr2      string `bun:"addr2" json:"addr2"`
	City       string `bun:"city" json:"city"`
	Postal     string `bun:"postal" json:"postal"`
	IsLoggedIn bool   `bun:"is_logged_in" json:"is_logged_in"`
}
