package connection

type Base struct {
	DriverName string
	Delimiter  string
}

// GetDelimiter implements the method Connection.GetDelimiter.
func (base *Base) GetDelimiter() string {
	return base.Delimiter
}

// Name implements the method Connection.Name.
func (base *Base) Name() string {
	return base.DriverName
}
