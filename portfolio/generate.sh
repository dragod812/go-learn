#!/bin/bash

protoc portfolio/pb/portfolio.proto --go_out=plugins=grpc:.