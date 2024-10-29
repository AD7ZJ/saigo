package corpus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorpus(t *testing.T) {
	result := Analysis("Are you serious? I knew you were.")
	assert.Equal(t, 6, len(result))
	assert.Equal(t, "you", result[0].Word)
	assert.Equal(t, 2, result[0].Count)
}
