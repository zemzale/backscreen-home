package main

//go:generate go tool oapi-codegen -generate chi-server,types,strict-server -o ./pkg/server/server.gen.go -package server ./spec/openapi.yaml
