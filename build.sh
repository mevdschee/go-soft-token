#!/bin/bash
env GOOS=windows GOARCH=386 go build
go build
