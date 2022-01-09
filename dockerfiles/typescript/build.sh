#!/bin/bash
docker build --no-cache -t kira-typescript:v0.0.1 .
docker tag kira-typescript:v0.0.1 florianwoelki/kira-typescript
