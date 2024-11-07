package main

import (
	"bufio"
	"log"
	"os"

	"github.com/neet-007/lsp_in_go/rpc"
)

func main() {
	logger := GetLogger("/home/moayed/personal/lsp_in_go/logs.txt")

	logger.Println("Starting...")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	for scanner.Scan() {
		msg := scanner.Text()
		handleMessage(logger, msg)
	}
}

func handleMessage(logger *log.Logger, msg any) {
	logger.Println(msg)
}

func GetLogger(fileName string) *log.Logger {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}

	logger := log.New(file, "[LSP_IN_GO] ", log.Ldate|log.Lshortfile)
	return logger
}
