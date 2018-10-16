package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	u := NewUser("1", "sam42@outlook.com")
	assert.Equal(t, "1", u.GetID())
	assert.Equal(t, "sam42@outlook.com", u.GetEmail())
}
