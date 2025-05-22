package aw

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/github/github-mcp-server/pkg/translations"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// CreateOrUpdateFile creates a tool to create or update a file.
func CreateOrUpdateFile(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("create_or_update_file",
			mcp.WithDescription(t("TOOL_CREATE_OR_UPDATE_FILE_DESCRIPTION", "Create or update a single local file.")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_CREATE_OR_UPDATE_FILE_USER_TITLE", "Create or update file"),
				ReadOnlyHint: toBoolPtr(false),
			}),
			mcp.WithString("path",
				mcp.Required(),
				mcp.Description("Path where to create/update the file"),
			),
			mcp.WithString("content",
				mcp.Required(),
				mcp.Description("Content of the file"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			path, err := requiredParam[string](request, "path")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			content, err := requiredParam[string](request, "content")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			// json.Marshal encodes byte arrays with base64, which is required for the API.
			contentBytes := []byte(content)

			// WriteFile creates a new file and writes data to it, or truncates and overwrites
			// existing file. The permissions are set to 0644 by default.
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("failed to get user home directory: %w", err)
			}
			pathToFile := filepath.Join(homeDir, path)
			if err := os.WriteFile(pathToFile, contentBytes, 0644); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to create/update file %q: %s", pathToFile, err)), nil
			}

			return mcp.NewToolResultText(fmt.Sprintf("file %q successfully created.", pathToFile)), nil
		}
}

// GetFileContents creates a tool to get the contents of a file.
func GetFileContents(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_file_contents",
			mcp.WithDescription(t("TOOL_GET_FILE_CONTENTS_DESCRIPTION", "Get the contents of a file from local file system.")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_GET_FILE_CONTENTS_USER_TITLE", "Get file contents"),
				ReadOnlyHint: toBoolPtr(true),
			}),
			mcp.WithString("path",
				mcp.Required(),
				mcp.Description("Path where to create/update the file"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			path, err := requiredParam[string](request, "path")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			// WriteFile creates a new file and writes data to it, or truncates and overwrites
			// existing file. The permissions are set to 0644 by default.
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("failed to get user home directory: %w", err)
			}
			pathToFile := filepath.Join(homeDir, path)
			content, err := os.ReadFile(pathToFile)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to read file %q: %s", pathToFile, err)), nil
			}

			return mcp.NewToolResultText(string(content)), nil
		}
}

// DeleteFile creates a tool to delete a file.
func DeleteFile(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("delete_file",
			mcp.WithDescription(t("TOOL_DELETE_FILE_DESCRIPTION", "Delete a file from a GitHub repository")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:           t("TOOL_DELETE_FILE_USER_TITLE", "Delete file"),
				ReadOnlyHint:    toBoolPtr(false),
				DestructiveHint: toBoolPtr(true),
			}),
			mcp.WithString("path",
				mcp.Required(),
				mcp.Description("Path where to create/update the file"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			path, err := requiredParam[string](request, "path")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			// WriteFile creates a new file and writes data to it, or truncates and overwrites
			// existing file. The permissions are set to 0644 by default.
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("failed to get user home directory: %w", err)
			}
			pathToFile := filepath.Join(homeDir, path)

			if err := os.Remove(pathToFile); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to delete file %q: %s", pathToFile, err)), nil
			}

			msg := fmt.Sprintf("File %q successfully deleted.", pathToFile)

			return mcp.NewToolResultText(msg), nil
		}
}
