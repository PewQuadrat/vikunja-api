package models

// Rights defines rights methods
type Rights interface {
	IsAdmin(*User) bool
	CanWrite(*User) bool
	CanRead(*User) bool
	CanDelete(*User) bool
	CanUpdate(*User) bool
	CanCreate(*User) bool
}