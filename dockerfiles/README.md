# Dockerfiles for Kira

This directory contains custom images and dockerfiles for custom code execution.

The sub directories always contain three files:

|File|Description|
|---|---|
|`Dockerfile`|Describes the to be build docker image.|
|`build.sh`|Builds the docker image and uses the local `Dockerfile`.|
|`deploy.sh`|Builds and deploys the TypeScript docker image.
