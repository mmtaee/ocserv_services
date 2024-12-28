package password

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatePassword(t *testing.T) {
	hashPassword := Create("random-password")
	assert.NotNil(t, hashPassword)
}

func TestCheck(t *testing.T) {
	hashPassword := Create("random-password")
	ok := Check("random-password", hashPassword)
	assert.True(t, ok)
}

func BenchmarkCreatePassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Create("random-password")
	}
}

func BenchmarkCheck(b *testing.B) {
	hashPassword := Create("random-password")
	for i := 0; i < b.N; i++ {
		Check("random-password", hashPassword)
	}
}
