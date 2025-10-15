package main

//go:generate go tool oapi-codegen -generate types,chi-server -o './pkg/server/server.gen.go' -package server ./spec/openapi.yaml
