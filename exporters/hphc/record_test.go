package hphc

import (
	"goconso/teleinfo"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRecord(t *testing.T) {
	frame := teleinfo.Frame{
		"PAPP": "001234",
		"PTEC": "HP..",
		"HCHC": "0987654",
		"HCHP": "1234567",
	}

	r := newRecord(frame)
	assert.NotNil(t, r)
	assert.False(t, r.Timestamp.IsZero())
	assert.EqualValues(t, 1234, r.PAPP_va)
	assert.EqualValues(t, 987654, r.HC_wh)
	assert.EqualValues(t, 1234567, r.HP_wh)
	assert.True(t, r.IsHP)

	frame["PTEC"] = "HC.."
	r = newRecord(frame)
	assert.NotNil(t, r)
	assert.False(t, r.IsHP)
}
