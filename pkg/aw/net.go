package aw

import (
	"context"
	"fmt"

	"github.com/github/github-mcp-server/pkg/translations"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tmc/langchaingo/llms/openai"
)

// NetConfig содержит параметры для сетевого запроса
// (model, url, apiKey)
type NetConfig struct {
	Model  string
	Url    string
	ApiKey string
}

// Query is a tool to query a net.
func Query(getClient GetClientFn, t translations.TranslationHelperFunc, netCfg NetConfig) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("ask_question",
			mcp.WithDescription(t("TOOL_ASK_QUESTION_DESCRIPTION", "To ask another neural network a question.")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_ASK_QUESTION_USER_TITLE", "Ask a question"),
				ReadOnlyHint: toBoolPtr(true),
			}),
			mcp.WithString("query",
				mcp.Required(),
				mcp.Description("Text query"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			query, err := requiredParam[string](request, "query")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			res, err := createAndSendRequest(ctx, query, netCfg)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("createAndSendRequest failed: %s", err)), nil
			}

			return mcp.NewToolResultText(res), nil
		}
}

func createAndSendRequest(ctx context.Context, query string, netCfg NetConfig) (string, error) {
	llm, err := openai.New(
		openai.WithBaseURL(netCfg.Url),
		openai.WithModel(netCfg.Model),
		openai.WithToken(netCfg.ApiKey),
	)
	if err != nil {
		return "", fmt.Errorf("openai.New failed: %w", err)
	}

	resp, err := llm.Call(ctx, query)
	if err != nil {
		return "", fmt.Errorf("llm.Call failed: %w", err)
	}

	return resp, nil
}
