package mocks

import "errors"

type Reader string

func (r Reader) Read([]byte) (int, error) {
	return 0, errors.New(string(r))
}
