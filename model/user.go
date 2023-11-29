package model

type UserID string

type User struct {
	UserID        UserID `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// Insert user into db, if already present do nothing
func CreateUserIfNotExist(user *User) error {
	statement := "INSERT INTO users (user_id, email, verified_email, name, given_name, family_name, picture, locale) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (user_id) DO NOTHING"

	_, err := db.Exec(statement, user.UserID, user.Email, user.VerifiedEmail, user.Name, user.GivenName, user.FamilyName, user.Picture, user.Locale)
	if err != nil {
		return err
	}
	return nil
}

func GetUserByID(id string) (*User, error) {
	statement := "SELECT user_id, email, verified_email, name, given_name, family_name, picture, locale FROM users WHERE user_id = $1"

	var user User
	err := db.QueryRow(statement, id).Scan(&user.UserID, &user.Email, &user.VerifiedEmail, &user.Name, &user.GivenName, &user.FamilyName, &user.Picture, &user.Locale)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
