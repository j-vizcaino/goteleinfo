package teleinfo

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadRawFrame(t *testing.T) {
	expectedFrame := []byte("frame content")
	bufferContent := bytes.NewBufferString("some junk prefix\r\n")
	bufferContent.WriteByte(0x2)
	bufferContent.Write(expectedFrame)
	bufferContent.WriteByte(0x3)
	bufferContent.WriteString("trailing junk")

	buf := bytes.NewReader(bufferContent.Bytes())
	reader := bufio.NewReader(buf)

	frame, err := readRawFrame(reader)
	assert.Nil(t, err)
	assert.Equal(t, frame, expectedFrame)
}

func TestReadFrameNoStart(t *testing.T) {
	invalidData := []byte("qwertyuiop")
	reader := NewReader(bytes.NewReader(invalidData), &historic)

	_, err := reader.ReadFrame()
	assert.Error(t, err)
}

func TestReadFrameNoEnd(t *testing.T) {
	invalidData := make([]byte, 0, bufio.MaxScanTokenSize+1)
	// Set frame start
	invalidData = append(invalidData, byte(0x2))
	for i := 1; i < cap(invalidData); i++ {
		invalidData = append(invalidData, byte(i%255))
	}

	reader := NewReader(bytes.NewReader(invalidData), &historic)
	_, err := reader.ReadFrame()
	assert.Error(t, err)
}

func TestReadFrameEOF(t *testing.T) {
	reader := NewReader(bytes.NewReader(nil), &historic)

	_, err := reader.ReadFrame()
	assert.Error(t, err)
	assert.True(t, errors.Is(err, io.EOF))
}

// Test for Historical Frame 
var historic = "historic"
var HistoricRawFramePayload = []byte("\rPAPP 012345  \r\nPTEC HP..  \r\nHCHC 09876543 @\r\nHCHP 1654800 K\n")
var HistoricRawFrame = []byte("\x02\rPAPP 012345  \r\nPTEC HP..  \r\nHCHC 09876543 @\r\nHCHP 1654800 K\n\x03")
var HistoricDecodedFrame = frame{
	"PAPP": "012345",
	"PTEC": "HP..",
	"HCHC": "09876543",
	"HCHP": "1654800",
}

func TestHistoricReadFrame(t *testing.T) {
	reader := NewReader(bytes.NewReader(HistoricRawFrame), &historic)

	frame, err := reader.ReadFrame()
	assert.Nil(t, err)
	assert.Equal(t, HistoricDecodedFrame, frame)
}


// Test for Standard Frame 
var standard = "standard"
var StandardRawFramePayload = []byte("\nADSC\t12345678900\tJ\r\nVTIC\t02\tJ\r\nDATE\tH191203125805\t\tF\r\nNGTF\t HC et Week-End \tU\r\nLTARF\t HEURE  CREUSE  \tK\r\nEAST\t017026175\t,\r\nEASF01\t004432429\t>\r\nEASF02\t009236697\tM\r\nEASF03\t003357049\tC\r\nEASF04\t000000000\t%\r\nEASF05\t000000000\t&\r\nEASF06\t000000000\t'\r\nEASF07\t000000000\t(\r\nEASF08\t000000000\t)\r\nEASF09\t000000000\t*\r\nEASF10\t000000000\t\"\r\nEASD01\t003706174\t<\r\nEASD02\t008070105\t6\r\nEASD03\t001746697\tJ\r\nEASD04\t003503199\tA\r\nIRMS1\t002\t0\r\nURMS1\t232\tA\r\nPREF\t09\tH\r\nPCOUP\t09\t\"\r\nSINSTS\t00519\tU\r\nSMAXSN\tH191203032158\t05706\tB\r\nSMAXSN-1\tH191202171658\t05661\t(\r\nCCASN\tH191203120000\t00521\t6\r\nCCASN-1\tH191203110000\t00415\tU\r\nUMOY1\tH191203125000\t232\t-\r\nSTGE\t003A8000\tA\r\nMSG1\tPAS DE          MESSAGE         \t<\r\nPRM\t01234567890123\t4\r\nRELAIS\t001\tC\r\nNTARF\t01\tN\r\nNJOURF\t00\t&\r\nNJOURF+1\t00\tB\r\nPJOURF+1\t00008002 0256C001 07568002 1156C001 14568002 NONUTILE NONUTILE NONUTILE NONUTILE NONUTILE NONUTILE\t:\r")
var StandardRawFrame = []byte("\x02\nADSC\t12345678900\tJ\r\nVTIC\t02\tJ\r\nDATE\tH191203125805\t\tF\r\nNGTF\t HC et Week-End \tU\r\nLTARF\t HEURE  CREUSE  \tK\r\nEAST\t017026175\t,\r\nEASF01\t004432429\t>\r\nEASF02\t009236697\tM\r\nEASF03\t003357049\tC\r\nEASF04\t000000000\t%\r\nEASF05\t000000000\t&\r\nEASF06\t000000000\t'\r\nEASF07\t000000000\t(\r\nEASF08\t000000000\t)\r\nEASF09\t000000000\t*\r\nEASF10\t000000000\t\"\r\nEASD01\t003706174\t<\r\nEASD02\t008070105\t6\r\nEASD03\t001746697\tJ\r\nEASD04\t003503199\tA\r\nIRMS1\t002\t0\r\nURMS1\t232\tA\r\nPREF\t09\tH\r\nPCOUP\t09\t\"\r\nSINSTS\t00519\tU\r\nSMAXSN\tH191203032158\t05706\tB\r\nSMAXSN-1\tH191202171658\t05661\t(\r\nCCASN\tH191203120000\t00521\t6\r\nCCASN-1\tH191203110000\t00415\tU\r\nUMOY1\tH191203125000\t232\t-\r\nSTGE\t003A8000\tA\r\nMSG1\tPAS DE          MESSAGE         \t<\r\nPRM\t01234567890123\t4\r\nRELAIS\t001\tC\r\nNTARF\t01\tN\r\nNJOURF\t00\t&\r\nNJOURF+1\t00\tB\r\nPJOURF+1\t00008002 0256C001 07568002 1156C001 14568002 NONUTILE NONUTILE NONUTILE NONUTILE NONUTILE NONUTILE\t:\r\x03")
var StandardDecodedFrame = frame{
"ADSC":"12345678900", 
"CCASN":"00521", 
"CCASN-1":"00415", 
"DATE":"", 
"EASD01":"003706174", 
"EASD02":"008070105", 
"EASD03":"001746697", 
"EASD04":"003503199", 
"EASF01":"004432429", 
"EASF02":"009236697", 
"EASF03":"003357049", 
"EASF04":"000000000", 
"EASF05":"000000000", 
"EASF06":"000000000", 
"EASF07":"000000000", 
"EASF08":"000000000", 
"EASF09":"000000000", 
"EASF10":"000000000", 
"EAST":"017026175", 
"IRMS1":"002", 
"LTARF":" HEURE  CREUSE  ", 
"MSG1":"PAS DE          MESSAGE         ", 
"NGTF":" HC et Week-End ", 
"NJOURF":"00", 
"NJOURF+1":"00", 
"NTARF":"01", 
"PCOUP":"09", 
"PJOURF+1":"00008002 0256C001 07568002 1156C001 14568002 NONUTILE NONUTILE NONUTILE NONUTILE NONUTILE NONUTILE", 
"PREF":"09", 
"PRM":"01234567890123", 
"RELAIS":"001", 
"SINSTS":"00519", 
"SMAXSN":"05706", 
"SMAXSN-1":"05661", 
"STGE":"003A8000", 
"UMOY1":"232", 
"URMS1":"232", 
"VTIC":"02",
}

func TestStandardReadFrame(t *testing.T) {
	reader := NewReader(bytes.NewReader(StandardRawFrame), &standard)

	frame, err := reader.ReadFrame()
	assert.Nil(t, err)
	assert.Equal(t, StandardDecodedFrame, frame)
}
