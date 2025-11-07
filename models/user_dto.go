package models

// CreateUserRequest represents the request body for creating a new user
type CreateUserRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email,max=100"`
	Password  string `json:"password" validate:"required,min=8,max=72"` // bcrypt max is 72 bytes
	FirstName string `json:"first_name" validate:"max=100"`
	LastName  string `json:"last_name" validate:"max=100"`
}

// UserResponse represents the user data returned in API responses
// It excludes sensitive fields and includes metadata
type UserResponse struct {
	ID                    uint   `json:"id"`
	Username              string `json:"username"`
	Email                 string `json:"email"`
	FirstName             string `json:"first_name,omitempty"`
	LastName              string `json:"last_name,omitempty"`
	ProfilePictureURL     string `json:"profile_picture_url,omitempty"`     // Presigned URL for profile picture
	ProfilePictureExpires string `json:"profile_picture_expires,omitempty"` // When the presigned URL expires
	CreatedAt             string `json:"created_at"`
	UpdatedAt             string `json:"updated_at"`
}

// UpdateUserRequest represents the request body for updating an existing user
// Only includes fields that can be updated
type UpdateUserRequest struct {
	Username  *string `json:"username" validate:"omitempty,min=3,max=50"`
	Email     *string `json:"email" validate:"omitempty,email,max=100"`
	FirstName *string `json:"first_name" validate:"omitempty,max=100"`
	LastName  *string `json:"last_name" validate:"omitempty,max=100"`
}

// ToUserResponse converts a User model to a UserResponse DTO
// Note: This method does not include profile picture URL. Use ToUserResponseWithPresignedURL for that.
func (u *User) ToUserResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToUserResponseWithPresignedURL converts a User model to a UserResponse DTO with profile picture presigned URL
func (u *User) ToUserResponseWithPresignedURL(profilePictureURL, profilePictureExpires string) *UserResponse {
	return &UserResponse{
		ID:                    u.ID,
		Username:              u.Username,
		Email:                 u.Email,
		FirstName:             u.FirstName,
		LastName:              u.LastName,
		ProfilePictureURL:     profilePictureURL,
		ProfilePictureExpires: profilePictureExpires,
		CreatedAt:             u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:             u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
