package utility

func NewConfigFromLocal[T any](decode Unmarshal, path string) (*T, error) {
	var conf T
	return &conf, nil
}
