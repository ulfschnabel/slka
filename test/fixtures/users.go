package fixtures

// User represents a Slack user for testing
type User struct {
	ID       string
	Name     string
	RealName string
	Email    string
	IsBot    bool
}

// GetTestUsers returns a set of test users
func GetTestUsers() []User {
	return []User{
		{
			ID:       "U001",
			Name:     "alice",
			RealName: "Alice Smith",
			Email:    "alice@example.com",
			IsBot:    false,
		},
		{
			ID:       "U002",
			Name:     "bob",
			RealName: "Bob Jones",
			Email:    "bob@example.com",
			IsBot:    false,
		},
		{
			ID:       "U003",
			Name:     "charlie",
			RealName: "Charlie Brown",
			Email:    "charlie@example.com",
			IsBot:    false,
		},
		{
			ID:       "USLACKBOT",
			Name:     "slackbot",
			RealName: "Slackbot",
			Email:    "",
			IsBot:    true,
		},
	}
}

// GetUserByID returns a test user by ID
func GetUserByID(id string) *User {
	users := GetTestUsers()
	for _, u := range users {
		if u.ID == id {
			return &u
		}
	}
	return nil
}

// GetUserByName returns a test user by name
func GetUserByName(name string) *User {
	users := GetTestUsers()
	for _, u := range users {
		if u.Name == name {
			return &u
		}
	}
	return nil
}

// GetUserByEmail returns a test user by email
func GetUserByEmail(email string) *User {
	users := GetTestUsers()
	for _, u := range users {
		if u.Email == email {
			return &u
		}
	}
	return nil
}
