package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_generateEmailCode(t *testing.T) {
	for i := 0; i < 1000; i++ {
		assert.Equal(t, len(generateEmailCode()), 6)
	}
}
