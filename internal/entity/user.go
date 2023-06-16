// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

// User -.
type User struct {
	Name     string `json:"name"       example:"Igor"`
	Email    string `json:"email"  example:"johndoe@wxample.com"`
	Password string `json:"password"     example:"285fguh2i3u1^%ih"`
	Auth     bool
}
