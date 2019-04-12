package routes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindActor(t *testing.T) {

	ppl := []People{
		People{
			Name: "Haroon",
			Photos: []string{
				"1",
				"2",
			},
		},
	}

	x := []string{"1", "2"}
	assert.Equal(t, ppl[0].Name, FindActor(ppl, x))

}
