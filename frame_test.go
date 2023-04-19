package teleinfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHistoricFrame(t *testing.T) {
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

	assert.Equal(t, "HC..", f.Type())
	assert.Equal(t, "historic", f.Mode())

	assert.Equal(t, map[string]string{
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
	}, f.GetMap())

	s, ok := f.GetStringField("PTEC")
	assert.True(t, ok)
	assert.Equal(t, "HP..", s)

	_, ok = f.GetStringField("missing")
	assert.False(t, ok)

	num, ok := f.GetUIntField("IMAX")
	assert.True(t, ok)
	assert.Equal(t, uint(36), num)

	_, ok = f.GetUIntField("missing")
	assert.False(t, ok)

	_, ok = f.GetUIntField("OPTARIF")
	assert.False(t, ok)
}

func TestDecodeHistoricFrame(t *testing.T) {
	InvalidChecksumValue := []byte("\rHCHC 09876543 e\n")
	InvalidChecksumLength := []byte("\rHCHC 09876543 ea\n")
	InvalidFieldCount := []byte("\rHCHC 09876543\n")

	r, err := decodeHistoricFrame(InvalidChecksumValue)
	assert.Nil(t, r)
	assert.Error(t, err)

	r, err = decodeHistoricFrame(InvalidChecksumLength)
	assert.Nil(t, r)
	assert.Error(t, err)

	r, err = decodeHistoricFrame(InvalidFieldCount)
	assert.Nil(t, r)
	assert.Error(t, err)

	r, err = decodeHistoricFrame(HistoricRawFramePayload)
	assert.Nil(t, err)
	assert.Equal(t, HistoricDecodedFrame, r)
}

func TestHistoricChecksum(t *testing.T) {
	name := []byte("PAPP")
	value := []byte("012345")
	expectedChecksum := byte(0x20)

	assert.Equal(t, expectedChecksum, historicChecksum(name, value))
}

// Test Standard Frame

func TestStandardFrame(t *testing.T) {
	f := frame{
		"ADSC":     "12345678900",
		"CCASN":    "00521",
		"CCASN-1":  "00415",
		"DATE":     "",
		"EASD01":   "003706174",
		"EASD02":   "008070105",
		"EASD03":   "001746697",
		"EASD04":   "003503199",
		"EASF01":   "004432429",
		"EASF02":   "009236697",
		"EASF03":   "003357049",
		"EASF04":   "000000000",
		"EASF05":   "000000000",
		"EASF06":   "000000000",
		"EASF07":   "000000000",
		"EASF08":   "000000000",
		"EASF09":   "000000000",
		"EASF10":   "000000000",
		"EAST":     "017026175",
		"IRMS1":    "002",
		"LTARF":    " HEURE  CREUSE  ",
		"MSG1":     "PAS DE          MESSAGE         ",
		"NGTF":     " HC et Week-End ",
		"NJOURF":   "00",
		"NJOURF+1": "00",
		"NTARF":    "01",
		"PCOUP":    "09",
		"PJOURF+1": "00008002 0256C001 07568002 1156C001 14568002 NONUTILE NONUTILE NONUTILE NONUTILE NONUTILE NONUTILE",
		"PREF":     "09",
		"PRM":      "01234567890123",
		"RELAIS":   "001",
		"SINSTS":   "00519",
		"SMAXSN":   "05706",
		"SMAXSN-1": "05661",
		"STGE":     "003A8000",
		"UMOY1":    "232",
		"URMS1":    "232",
		"VTIC":     "02",
	}

	assert.Equal(t, "standard", f.Mode())

	s, ok := f.GetStringField("NGTF")
	assert.True(t, ok)
	assert.Equal(t, " HC et Week-End ", s)

	_, ok = f.GetStringField("missing")
	assert.False(t, ok)

	num, ok := f.GetUIntField("SINSTS")
	assert.True(t, ok)
	assert.Equal(t, uint(519), num)

	_, ok = f.GetUIntField("missing")
	assert.False(t, ok)

	_, ok = f.GetUIntField("OPTARIF")
	assert.False(t, ok)
}

func TestStandardChecksum(t *testing.T) {
	name := []byte("ADSC")
	value := []byte("12345678900")
	expectedChecksum := byte(0x4A)

	assert.Equal(t, expectedChecksum, standardChecksum(name, []byte(""), value))
}

func TestStandardChecksumWithTimestamp(t *testing.T) {
	name := []byte("CCASN")
	timestamp := []byte("H191203120000")
	value := []byte("00521")
	expectedChecksum := byte(0x36)

	assert.Equal(t, expectedChecksum, standardChecksum(name, timestamp, value))
}
