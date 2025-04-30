package logger

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"
)

type lokiPayload struct {
	Streams []lokiStream `json:"streams"`
}

type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

var lokiURL string

func InitLokiLogger(url string) error {
	lokiURL = url
	return nil // nothing to initialize, keep signature for compatibility
}

func LogInfo(message string) {
	sendLoki("info", message)
}

func LogWarn(message string) {
	sendLoki("warn", message)
}

func LogError(message string, err error) {
	sendLoki("error", message+" | err: "+err.Error())
}

func sendLoki(level string, message string) {
	if lokiURL == "" {
		lokiURL = os.Getenv("LOKI_URL")
		if lokiURL == "" {
			return
		}
	}

	payload := lokiPayload{
		Streams: []lokiStream{
			{
				Stream: map[string]string{
					"job":   "convert-rate-api",
					"level": level,
				},
				Values: [][]string{
					{strconv.FormatInt(time.Now().UnixNano(), 10), message},
				},
			},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	http.Post(lokiURL, "application/json", bytes.NewBuffer(data))
}

func ShutdownLoki() {} // no-op
