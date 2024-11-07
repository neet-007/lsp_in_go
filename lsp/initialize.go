package lsp

type InitializeRequest struct {
	Request
	Params InitializeRequestParams `json:"params"`
}

type InitializeRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResponse struct {
	Respone
	Result InitializeResult `json:"result"`
}

type InitializeResult struct {
	ServerCapabilities ServerCapabilities `json:"capabilities"`
	ServerInfo         ServerInfo         `json:"serverInfo"`
}

type ServerCapabilities struct {
	TextDocumentSync int `json:"textDocumentSync"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Respone: Respone{
			RPC: "2.0",
			Id:  &id,
		},
		Result: InitializeResult{
			ServerCapabilities: ServerCapabilities{
				TextDocumentSync: 1,
			},
			ServerInfo: ServerInfo{
				Name:    "lsp_in_go",
				Version: "0.0.1",
			},
		},
	}
}