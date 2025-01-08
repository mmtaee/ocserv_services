package password

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatePassword(t *testing.T) {
	hashPassword := NewPassword("random-password")
	assert.NotNil(t, hashPassword)
}

func TestCheck(t *testing.T) {
	pass := NewPassword("random-password")
	ok := Check("random-password", pass.Hash, pass.Salt)
	assert.True(t, ok)
}

func BenchmarkCreatePassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewPassword("random-password")
	}
}

func BenchmarkCheck(b *testing.B) {
	pass := NewPassword("random-password")
	for i := 0; i < b.N; i++ {
		Check("random-password", pass.Hash, pass.Salt)
	}
}
