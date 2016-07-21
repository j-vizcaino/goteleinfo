package teleinfo

type Exporter interface {
	ExportFrame(Frame) error
}
