package service

import "context"

func (s *service) Hello(ctx context.Context, name string) (resp string, err error) {
	resp = "Hello!"

	if name != "" {
		resp = "Hello, " + name + "!"
	}
	return resp, nil
}
