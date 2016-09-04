package hphc

import (
	"goteleinfo/teleinfo"
	"strconv"
	"time"
)

type record struct {
	Timestamp time.Time
	PAPP_va   uint32
	HC_wh     uint32
	HP_wh     uint32
	IsHP      bool
}

func extractNumber(value string) uint32 {
	num, _ := strconv.ParseUint(value, 10, 32)
	return uint32(num)
}

func newRecord(f teleinfo.Frame) *record {
	return &record{
		Timestamp: time.Now(),
		PAPP_va:   extractNumber(f["PAPP"]),
		HC_wh:     extractNumber(f["HCHC"]),
		HP_wh:     extractNumber(f["HCHP"]),
		IsHP:      f["PTEC"] == "HP..",
	}
}
