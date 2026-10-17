package auth

// User is the auth-relevant projection of a users row.
type User struct {
	ID           int64
	PublicID     string
	Username     string
	DisplayName  string
	HashPassword string
	Email        string
	Phone        *string
	IsVerified   bool
}
