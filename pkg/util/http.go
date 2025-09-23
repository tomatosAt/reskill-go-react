package util

import (
	"context"
	"fmt"
)

func GetHttpRequestId(ctx context.Context) string {
	requestId, ok := ctx.Value("requestid").(string)
	if ok {
		return requestId
	}
	return ""
}

type HttpSkipper struct {
	Rule map[string]struct{}
}

func NewHttpSkipper() *HttpSkipper {
	return &HttpSkipper{Rule: map[string]struct{}{}}
}

func (s *HttpSkipper) Add(m string, p string) {
	s.Rule[fmt.Sprintf("%s|%s", m, p)] = struct{}{}
}

func (s *HttpSkipper) Has(m string, p string) bool {
	if _, ok := s.Rule[fmt.Sprintf("%s|%s", m, p)]; ok {
		return ok
	}
	return false
}
