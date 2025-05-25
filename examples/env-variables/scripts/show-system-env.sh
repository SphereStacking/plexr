#!/bin/bash

echo "Common System Environment Variables"
echo "====================================="
echo ""

echo "USER: $USER"
echo "HOME: $HOME"
echo "PATH: $PATH" | fold -w 60
echo "SHELL: $SHELL"
echo "PWD: $PWD"
echo "LANG: $LANG"
echo ""

echo "Total environment variables: $(env | wc -l)"
echo ""