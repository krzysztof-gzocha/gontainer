package compiler

import (
	"fmt"
	"testing"

	"github.com/gomponents/gontainer-helpers/caller"
	"github.com/stretchr/testify/assert"
)

type inputOutputScenario struct {
	input  interface{}
	output interface{}
}

func doTestInputOutput(t *testing.T, fn interface{}, scenarios ...inputOutputScenario) {
	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			output := caller.MustCall(fn, s.input)
			assert.Equal(t, s.output, output[0])
		})
	}
}
