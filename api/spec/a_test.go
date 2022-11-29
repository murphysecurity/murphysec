package spec

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidate(t *testing.T) {
	assert.NoError(t, GetSpec().Validate(context.TODO()))
}
