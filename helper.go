package sqlparser

// get pointer of value
// https://github.com/golang/go/issues/45624#issuecomment-1843832599
//
//lint:ignore U1000 for testing
func addr[T any](v T) *T { return &v }
