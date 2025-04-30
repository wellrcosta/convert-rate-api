package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type lokiWriter struct {
	url string
}

func NewLokiLogger() zerolog.Logger {
	writer := &lokiWriter{
		url: os.Getenv("LOKI_URL"),
	}

	log.Logger = zerolog.New(writer).With().Timestamp().Logger()
	return log.Logger
}

func (lw *lokiWriter) Write(p []byte) (n int, err error) {
	lokiPayload := map[string]interface{}{
		"streams": []map[string]interface{}{
			{
				"stream": map[string]string{
					"app": "exchange-rate-api",
				},
				"values": [][]string{
					{
						// Timestamp in nanoseconds as string
						time.Now().UTC().Format("20060102150405"), string(p),
					},
				},
			},
		},
	}

	payload, err := json.Marshal(lokiPayload)
	if err != nil {
		return 0, err
	}

	resp, err := http.Post(lw.url, "application/json", bytes.NewReader(payload))
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close response body")
		}
	}(resp.Body)

	return len(p), nil
}
