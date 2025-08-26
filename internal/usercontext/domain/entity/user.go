package entity

import "time"

type UserID string
type Email string

type User struct {
    ID        UserID
    Email     Email
    CreatedAt time.Time
}

