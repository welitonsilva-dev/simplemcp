package server

import (
	"encoding/json"
	"net/http"

	"simplemcp_v0.1/internal/protocol"
	"simplemcp_v0.1/internal/tools"
)

type MCPServer struct {
	tools map[string]tools.Tool
}

func NewServer(toolList []tools.Tool) *MCPServer {
	toolMap := make(map[string]tools.Tool)
	for _, t := range toolList {
		toolMap[t.Name()] = t
	}
	return &MCPServer{tools: toolMap}
}

func (s *MCPServer) TestTools(w http.ResponseWriter, r *http.Request) {
	var req protocol.Request
	json.NewDecoder(r.Body).Decode(&req)

	res := protocol.Response{ID: req.ID}

	tool, exists := s.tools[req.Method]
	if !exists {
		res.Error = "tool not found"
		json.NewEncoder(w).Encode(res)
		return
	}

	result, err := tool.Execute(req.Params)
	if err != nil {
		res.Error = err.Error()
	} else {
		res.Result = result
	}

	json.NewEncoder(w).Encode(res)
}

func (s *MCPServer) Tools() map[string]tools.Tool {
	return s.tools
}
