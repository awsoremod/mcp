# MCP Server

tools:
- ask_question (Ask a question to the neural network)
- get_file_contents (Read information from a file, path relative to the home directory)
- create_or_update_file
- delete_file

## Getting Started

For MCP: set AI_API_CODE in the env in the mcp.json file  
For local: create a `.env` file based on the `.env.example` template.

```bash
make build
make run
```

Give input:
```json
{"jsonrpc":"2.0","method":"tools/call","params":{"name":"ask_question","arguments":{"query":"What is your name? Answer very briefly."}},"id":1}
```

Output:
```json
{"jsonrpc":"2.0","id":1,"result":{"content":[{"type":"text","text":"My name is DeepSeek-R1."}]}}
```