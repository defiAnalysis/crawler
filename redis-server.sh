#! /bin/bash


docker run -p 6379:6379 -v $HOME/dev/data/redis:/data  -d redis:latest redis-server --appendonly yes --requirepass "QWER1234"


