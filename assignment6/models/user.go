package models

import "sync"

type User struct {
	ID               int    `json:"id"`
	Email            string `json:"email"`
	Password         string `json:"-"`
	Role             string `json:"role"`        // "user" or "admin"
	IsVerified       bool   `json:"is_verified"`
	VerificationCode string `json:"-"`
}

var (
	Users = make(map[string]*User)
	Mu    sync.RWMutex
	NextID = 1
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)
