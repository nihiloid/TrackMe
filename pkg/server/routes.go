package server

import (
	"github.com/pagpeter/trackme/pkg/types"
)

func apiAll(res types.Response, _ map[string][]string) ([]byte, string, error) {
	return []byte(res.ToJson()), "application/json", nil
}
