package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/neet-007/lsp_in_go/lsp"
	"github.com/neet-007/lsp_in_go/rpc"
)

func main() {
	logger := GetLogger("/home/moayed/personal/lsp_in_go/logs.txt")

	logger.Println("Starting...")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, content, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Error:%s\n", err)
		}
		handleMessage(logger, method, content)
	}
}

func handleMessage(logger *log.Logger, method string, content []byte) {
	logger.Printf("Message with method:%s\n", method)
	switch method {
	case "initialize":
		{
			var initilialzeRequest lsp.InitializeRequest
			if err := json.Unmarshal(content, &initilialzeRequest); err != nil {
				logger.Printf("Error while parsing:%s\n", err)
				return
			}
			logger.Printf("Message with name:%s version:%s\n", initilialzeRequest.Params.ClientInfo.Name, initilialzeRequest.Params.ClientInfo.Version)
			msg := lsp.NewInitializeResponse(initilialzeRequest.Id)
			encoded := rpc.EncodeMessage(msg)

			writer := os.Stdout
			writer.Write([]byte(encoded))

			logger.Println("Sent message")
		}
	case "textDocument/didOpen":
		{
			var didOpenTextDocumentNotification lsp.DidOpenTextDocumentNotification
			if err := json.Unmarshal(content, &didOpenTextDocumentNotification); err != nil {
				logger.Printf("Error while parsing:%s\n", err)
				return
			}
			logger.Printf("text document with uri:%s content:%s\n", didOpenTextDocumentNotification.Params.TextDocument.URI,
				didOpenTextDocumentNotification.Params.TextDocument.Text)
		}
	}
}

func GetLogger(fileName string) *log.Logger {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}

	logger := log.New(file, "[LSP_IN_GO] ", log.Ldate|log.Lshortfile)
	return logger
}
