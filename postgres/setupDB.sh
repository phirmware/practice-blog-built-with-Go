#! /usr/bin/env bash

docker rmi practiceblog

docker build -f dockerfile -t practiceblog .

docker run --rm -d  -p 5432:5432 -e POSTGRES_DB=practiceblog practiceblog:latest

echo "Run GO file main.go TO CHECK IF DB WAS SUCCESSFULLY CREATED"