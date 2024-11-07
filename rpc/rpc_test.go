package rpc_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/neet-007/lsp_in_go/rpc"
)

func TestEncodeMessage(t *testing.T) {
	contents := []string{"\"test\"", "\"test2\"", "\"test3\""}
	exp := make([]string, 0, len(contents))

	for _, content := range contents {
		exp = append(exp, fmt.Sprintf("Content-Length: %d\r\n\r\n %s", len(content), content))
	}

	for i, content := range contents {
		act := rpc.EncodeMessage(strings.Replace(content, "\"", "", 2))

		if exp[i] != act {
			t.Fatalf("exp %s VS act %s\n", exp[i], act)
		}
	}
}

func TestDecodeMessage(t *testing.T) {
	contents := []string{"{\"method\":\"test\"}", "{\"method\":\"test2\"}", "{\"method\":\"test3\"}"}
	exp := make([]struct {
		Method     string
		Message    []byte
		LenMessage int
	}, 0, len(contents))

	for _, content := range contents {
		// Decode JSON to get method
		var jsonData struct {
			Method string `json:"method"`
		}
		if err := json.Unmarshal([]byte(content), &jsonData); err != nil {
			t.Fatalf("failed to parse content as JSON: %v", err)
		}

		exp = append(exp, struct {
			Method     string
			Message    []byte
			LenMessage int
		}{
			Method:     jsonData.Method,
			Message:    []byte(fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)),
			LenMessage: len(content),
		})
	}

	for i, content := range contents {
		encodedMessage := []byte(fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content))
		method, decodedContent, err := rpc.DecodeMessage(encodedMessage)
		if err != nil {
			t.Fatal(err)
		}

		e := exp[i]
		if e.Method != method {
			t.Fatalf("METHOD\nexp %s VS act %s\n", e.Method, method)
		}
		if !bytes.Equal(e.Message[len("Content-Length: %d\r\n\r\n"):], decodedContent) {
			t.Fatalf("CONTENT\nexp %s VS act %s\n", e.Message[len(fmt.Sprintf("Content-Length: %d\r\n\r\n", e.LenMessage)):], decodedContent)
		}
		if e.LenMessage != len(decodedContent) {
			t.Fatalf("LEN\nexp %d VS act %d\n", e.LenMessage, len(decodedContent))
		}
	}
}
