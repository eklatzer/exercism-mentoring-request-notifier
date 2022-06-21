#!/bin/bash

if [ "$(gofmt -l . | wc -l)" -gt 0 ]; then
  echo "format errors found, check code"
  exit 1
fi
