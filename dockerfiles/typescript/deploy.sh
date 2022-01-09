#!/bin/bash
sh ./build.sh
docker login
docker push florianwoelki/kira-typescript
