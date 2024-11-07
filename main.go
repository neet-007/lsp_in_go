package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/neet-007/lsp_in_go/analysis"
	"github.com/neet-007/lsp_in_go/lsp"
	"github.com/neet-007/lsp_in_go/rpc"
)

func main() {
	logger := GetLogger("/home/moayed/personal/lsp_in_go/logs.txt")

	logger.Println("Starting...")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)
	state := analysis.NewState()
	writer := os.Stdout

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, content, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Error:%s\n", err)
		}
		handleMessage(logger, writer, state, method, content)
	}
}

func handleMessage(logger *log.Logger, writer io.Writer, state analysis.State, method string, content []byte) {
	logger.Printf("Message with method:%s\n", method)
	switch method {
	case "initialize":
		{
			var request lsp.InitializeRequest
			if err := json.Unmarshal(content, &request); err != nil {
				logger.Printf("Hey, we couldn't parse this: %s", err)
			}

			logger.Printf("Connected to: %s %s",
				request.Params.ClientInfo.Name,
				request.Params.ClientInfo.Version)

			msg := lsp.NewInitializeResponse(request.Id)
			writeResponse(writer, msg)

			logger.Print("Sent the reply")
		}
	case "textDocument/didOpen":
		{
			var didOpenTextDocumentNotification lsp.DidOpenTextDocumentNotification
			if err := json.Unmarshal(content, &didOpenTextDocumentNotification); err != nil {
				logger.Printf("Error did open:%s\n", err)
				return
			}
			logger.Printf("text document with uri:%s\n", didOpenTextDocumentNotification.Params.TextDocument.URI)
			state.OpenDocument(didOpenTextDocumentNotification.Params.TextDocument.URI,
				didOpenTextDocumentNotification.Params.TextDocument.Text)
		}
	case "textDocument/didChange":
		{
			var didChangeTextDocumentNotification lsp.TextDocumentDidChangeNotification
			if err := json.Unmarshal(content, &didChangeTextDocumentNotification); err != nil {
				logger.Printf("textDocument/didChange: %s", err)
				return
			}

			logger.Printf("Changed: %s", didChangeTextDocumentNotification.Params.TextDocument.URI)
			logger.Println(didChangeTextDocumentNotification)
			logger.Printf("text document with uri:%s\n", didChangeTextDocumentNotification.Params.TextDocument.URI)
			for _, change := range didChangeTextDocumentNotification.Params.ContentChanges {
				state.OpenDocument(didChangeTextDocumentNotification.Params.TextDocument.URI, change.Text)
			}
		}
	case "textDocument/definition":
		{
			var request lsp.HoverRequest
			if err := json.Unmarshal(content, &request); err != nil {
				logger.Printf("textDocument/hover: %s", err)
				return
			}

			response := state.Hover(request.Id, request.Params.TextDocument.URI, request.Params.Position)

			writeResponse(writer, response)
		}
	}
}

func writeResponse(writer io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	writer.Write([]byte(reply))

}

func GetLogger(fileName string) *log.Logger {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}

	logger := log.New(file, "[LSP_IN_GO] ", log.Ldate|log.Lshortfile)
	return logger
}
