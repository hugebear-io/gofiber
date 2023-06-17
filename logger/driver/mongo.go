package driver

type mongoWriter struct {
}

func NewMongoWriter() *mongoWriter {
	return &mongoWriter{}
}

func (w mongoWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}
