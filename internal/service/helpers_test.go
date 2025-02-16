package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashAndCheckPassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		wrongPwd    string
		expectError bool
	}{
		{
			name:        "Valid password",
			password:    "my_secret_password",
			wrongPwd:    "wrong_password",
			expectError: false,
		},
		{
			name:        "Empty password",
			password:    "",
			wrongPwd:    "something",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := hashPassword(tt.password)

			if tt.expectError {
				assert.Error(t, err, "Should return error for empty password")
				assert.Contains(t, err.Error(), "password cannot be empty")
				return
			}

			assert.NoError(t, err, "Hash generation failed")
			assert.NotEmpty(t, hash, "Hash should not be empty")

			match := comparePassword(hash, tt.password)
			assert.True(t, match, "Password should match the hash")

			match = comparePassword(hash, tt.wrongPwd)
			assert.False(t, match, "Wrong password should not match the hash")
		})
	}
}
