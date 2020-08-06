#!/bin/bash

cd frontend
yarn build

cd ../http
rice embed-go

cd ..
go build

./filebrowser