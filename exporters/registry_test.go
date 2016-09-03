package exporters

import (
	"goconso/teleinfo"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	assert.Equal(t, []string{}, List())
	assert.Nil(t, FindByName("foo"))
}

type barExporter struct{}

func (b *barExporter) ExportFrame(teleinfo.Frame) error {
	return nil
}

func TestRegister(t *testing.T) {
	fooFactory := func() (teleinfo.Exporter, error) {
		return nil, nil
	}
	barResult := &barExporter{}
	barFactory := func() (teleinfo.Exporter, error) {
		return barResult, nil
	}

	Register("foo", fooFactory)
	Register("bar", barFactory)

	factories := List()
	assert.Len(t, factories, 2)
	assert.Contains(t, factories, "foo")
	assert.Contains(t, factories, "bar")

	// NOTE: Go cannot compare function reliably.
	//       Call factory function and compare result instead.
	foo := FindByName("foo")
	assert.NotNil(t, foo)
	res, err := foo()
	assert.Nil(t, res)
	assert.Nil(t, err)

	bar := FindByName("bar")
	assert.NotNil(t, bar)
	res, err = bar()
	assert.Equal(t, barResult, res)
	assert.Nil(t, err)

	assert.Nil(t, FindByName("missing"))
}
