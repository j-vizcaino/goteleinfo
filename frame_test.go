package teleinfo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFrame(t *testing.T) {
	f := frame{
		"PAPP":     "00340",
		"HHPHC":    "D",
		"ADCO":     "031028217014",
		"OPTARIF":  "HC..",
		"PTEC":     "HP..",
		"IINST":    "001",
		"MOTDETAT": "000000",
		"ISOUSC":   "45",
		"HCHC":     "016771964",
		"HCHP":     "020267321",
		"IMAX":     "036",
	}

	assert.Equal(t,"HC..", f.Type())

	s, ok := f.GetStringField("PTEC")
	assert.True(t, ok)
	assert.Equal(t, "HP..", s)

	s, ok = f.GetStringField("missing")
	assert.False(t, ok)

	num, ok := f.GetUIntField("IMAX")
	assert.True(t, ok)
	assert.Equal(t, uint(36), num)

	num, ok = f.GetUIntField("missing")
	assert.False(t, ok)

	num, ok = f.GetUIntField("OPTARIF")
	assert.False(t, ok)
}

func TestDecodeFrame(t *testing.T) {
	InvalidChecksumValue := []byte("\rHCHC 09876543 e\n")
	InvalidChecksumLength := []byte("\rHCHC 09876543 ea\n")
	InvalidFieldCount := []byte("\rHCHC 09876543\n")

	r, err := decodeFrame(InvalidChecksumValue)
	assert.Nil(t, r)
	assert.Error(t, err)

	r, err = decodeFrame(InvalidChecksumLength)
	assert.Nil(t, r)
	assert.Error(t, err)

	r, err = decodeFrame(InvalidFieldCount)
	assert.Nil(t, r)
	assert.Error(t, err)

	r, err = decodeFrame(RawFramePayload)
	assert.Nil(t, err)
	assert.Equal(t, DecodedFrame, r)
}

func TestComputeChecksum(t *testing.T) {
	name := []byte("PAPP")
	value := []byte("012345")
	expectedChecksum := byte(0x20)

	assert.Equal(t, expectedChecksum, computeChecksum(name, value))
}
