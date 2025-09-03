package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDaysCount_ReturnsPositiveNumber(t *testing.T) {
	s := New()
	res := s.DaysCount()
	assert.True(t, res > 0, "DaysCount should return positive int")
}

func TestDaysCount_ReturnsInt64(t *testing.T) {
	s := New()
	res := s.DaysCount()
	assert.IsType(t, int64(0), res, "Should return int64")
}
