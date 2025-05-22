package aw

import (
	"context"

	"github.com/awsoremod/mcp/pkg/toolsets"
	"github.com/github/github-mcp-server/pkg/translations"
	"github.com/google/go-github/v69/github"
)

type GetClientFn func(context.Context) (*github.Client, error)

var DefaultTools = []string{"all"}

func InitToolsets(passedToolsets []string, readOnly bool, getClient GetClientFn, t translations.TranslationHelperFunc, netCfg NetConfig) (*toolsets.ToolsetGroup, error) {
	// Create a new toolset group
	tsg := toolsets.NewToolsetGroup(readOnly)

	// Define all available features with their default state (disabled)
	// Create toolsets
	files := toolsets.NewToolset("file", "File related tools").
		AddReadTools(
			toolsets.NewServerTool(GetFileContents(getClient, t)),
		).
		AddWriteTools(
			toolsets.NewServerTool(CreateOrUpdateFile(getClient, t)),
			toolsets.NewServerTool(DeleteFile(getClient, t)),
		)
	net := toolsets.NewToolset("net", "Network related tools").
		AddReadTools(
			toolsets.NewServerTool(Query(getClient, t, netCfg)),
		)

	// Keep experiments alive so the system doesn't error out when it's always enabled
	experiments := toolsets.NewToolset("experiments", "Experimental features that are not considered stable yet")

	// Add toolsets to the group
	tsg.AddToolset(files)
	tsg.AddToolset(net)
	tsg.AddToolset(experiments)
	// Enable the requested features

	if err := tsg.EnableToolsets(passedToolsets); err != nil {
		return nil, err
	}

	return tsg, nil
}

func toBoolPtr(b bool) *bool {
	return &b
}
