local_resource(
  'webserver-api-build',
  cmd='CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/webserver-api main.go',
  deps=['main.go', 'go.mod', 'go.sum', "http", "mux"],
  labels=['api']
)

docker_build(
  'webserver-api-image',
  '.',
  dockerfile='Dockerfile',
  only=['build/webserver-api'],
  live_update=[
    sync('./build/webserver-api', '/app/webserver-api')
  ]
)

docker_compose('./compose.yml')
dc_resource('webserver-api', labels=["api"])
