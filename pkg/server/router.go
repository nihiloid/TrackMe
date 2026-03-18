package server

import (
	"fmt"
	"strings"
	"time"

	"github.com/pagpeter/trackme/pkg/tls"
	"github.com/pagpeter/trackme/pkg/types"
)

func Log(msg string) {
	t := time.Now()
	formatted := t.Format("2006-01-02 15:04:05")
	fmt.Printf("[%v] %v\n", formatted, msg)
}

func cleanIP(ip string) string {
	return strings.Replace(strings.Replace(ip, "]", "", -1), "[", "", -1)
}

// Router returns bytes, content type, and error that should be sent to the client
func Router(_ string, res types.Response, srv *Server) ([]byte, string, error) {
	if v, ok := srv.GetTCPFingerprints().Load(res.IP); ok {
		res.TCPIP = v.(types.TCPIPDetails)
	}
	res.Donate = "Please consider donating to keep this API running. Visit https://tls.peet.ws"
	if res.TLS != nil {
		// Use QUIC JA4 for HTTP/3 connections
		if res.HTTPVersion == "h3" {
			res.TLS.JA4 = tls.CalculateJa4QUIC(res.TLS)
			res.TLS.JA4_r = tls.CalculateJa4QUIC_r(res.TLS)
		} else {
			res.TLS.JA4 = tls.CalculateJa4(res.TLS)
			res.TLS.JA4_r = tls.CalculateJa4_r(res.TLS)
		}
		Log(fmt.Sprintf("%v %v %v %v %v", cleanIP(res.IP), res.Method, res.HTTPVersion, res.Path, res.TLS.JA3Hash))
	} else {
		Log(fmt.Sprintf("%v %v %v %v %v", cleanIP(res.IP), res.Method, res.HTTPVersion, res.Path, "-"))
	}

	// Force all requests to /api/all
	responseBytes, ctype, err := apiAll(res, nil)
	if err != nil {
		return nil, "", err
	}

	// Log the pretty-printed JSON response asynchronously
	go LogResponse(responseBytes)

	return responseBytes, ctype, nil
}
