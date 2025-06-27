# Input Setup
```bash
cd ./cli/
source .env

# Start HTTP server in background for HTTP extractor
cd ./e2e/_input/mcp/
npx http-server -p 8081 &
HTTP_SERVER_PID=$!
cd ../../../

# Extract current state
rm -f ./e2e/_input/mcp/current.sqlite
./hyaline --debug extract current --config ./e2e/_input/mcp/config.yml --system mcp-test --output ./e2e/_input/mcp/current.sqlite

# Stop HTTP server
kill $HTTP_SERVER_PID
```