# Code Execution Environment

The first version of kira (can be found at branch `v1`) was not fast enough to meet the requirements of parallel code execution processing.

Therefore, the new v2 architecture of kira is based on a single container and the API will be exposed in this container.
The code will be run in the container but without starting extra containers or using any sort of container management (e.g. Docker).
