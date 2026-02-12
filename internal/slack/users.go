package slack

import (
	"fmt"
	"sort"
	"strings"
)

// UserService handles user-related operations
type UserService struct {
	client Client
}

// NewUserService creates a new user service
func NewUserService(client Client) *UserService {
	return &UserService{client: client}
}

// ListUsersOptions contains options for listing users
type ListUsersOptions struct{
	IncludeBots    bool
	IncludeDeleted bool
	Limit          int
}

// List returns all users in the workspace
func (s *UserService) List(opts ListUsersOptions) ([]UserInfo, error) {
	users, err := s.client.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	result := make([]UserInfo, 0, len(users))
	for _, user := range users {
		// Filter based on options
		if !opts.IncludeBots && user.IsBot {
			continue
		}
		if !opts.IncludeDeleted && user.Deleted {
			continue
		}

		result = append(result, convertUser(user))
	}

	// Sort by Updated timestamp descending (most recently updated first)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Updated > result[j].Updated
	})

	// Apply limit after sorting
	if opts.Limit > 0 && len(result) > opts.Limit {
		result = result[:opts.Limit]
	}

	return result, nil
}

// Lookup finds a user by name or email
func (s *UserService) Lookup(query string, byField string) (*UserInfo, error) {
	if byField == "auto" {
		// Auto-detect: if contains @, search by email, otherwise by name
		if strings.Contains(query, "@") {
			byField = "email"
		} else {
			byField = "name"
		}
	}

	if byField == "email" {
		user, err := s.client.GetUserByEmail(query)
		if err != nil {
			return nil, fmt.Errorf("user not found by email: %w", err)
		}
		info := convertUser(*user)
		return &info, nil
	}

	// Search by name
	users, err := s.client.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	for _, user := range users {
		if user.Name == query || user.RealName == query {
			info := convertUser(user)
			return &info, nil
		}
	}

	return nil, fmt.Errorf("user not found: %s", query)
}

// ResolveUserNames fetches info for a set of user IDs and returns a map.
// Only fetches the users requested, not the entire workspace.
func ResolveUserNames(client Client, userIDs []string) map[string]UserInfo {
	result := make(map[string]UserInfo, len(userIDs))
	for _, id := range userIDs {
		user, err := client.GetUserInfo(id)
		if err == nil {
			result[id] = convertUser(*user)
		}
	}
	return result
}

// ResolveUser resolves a user name or email to a user ID
func (s *UserService) ResolveUser(user string) (string, error) {
	// If it starts with U, it's probably an ID
	if strings.HasPrefix(user, "U") {
		return user, nil
	}

	// Try to look up by name or email
	userInfo, err := s.Lookup(user, "auto")
	if err != nil {
		return "", err
	}

	return userInfo.ID, nil
}
