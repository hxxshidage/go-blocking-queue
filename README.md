# go-blocking-queue
A blocking queue based on go channel, supporting timeout and thread safe

## Features
- Blocking queue with timeout control
- Minimalist API, based on channel, no other dependencies
- Timeout control, reusing timer, reducing GC pressure under high concurrency
- thread safe
## Usage
Reference test cases.
## Defects

- Capacity must be specified, Unable to cope with scenarios where capacity cannot be estimated. [ Can be parameterized](https://github.com/emirpasic/gods)