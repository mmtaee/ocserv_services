package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	token1 := Create(1, time.Now().Add(24*time.Hour))
	assert.NotEmpty(t, token1)
}

func BenchmarkGenerateToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Create(uint(i), time.Now().Add(24*time.Hour))
	}
}
