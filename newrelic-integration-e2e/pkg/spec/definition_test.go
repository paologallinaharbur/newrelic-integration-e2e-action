package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseSpecFile(t *testing.T){
	var sample = `
		description: This is a demo
	`
	spec,err:= ParseSpecFile([]byte(sample))
	assert.Nil(t, err)
	assert.Equal(t, "This is a demo",spec.Description)
}