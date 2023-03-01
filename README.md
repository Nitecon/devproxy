# devproxy
Quick and dirty proxy to split frontend / backend for local development, simulating an ALB doing Path based proxying.

## Why?
I was working on a project that had a frontend (svelte) and backend (golang), and I wanted to be able to run them both locally, as they would be in production with an ALB doing path based routing.

## Usage
To use it download the current os version from the releases page then:
- Extract the zip file
- Configure the config.yaml file
- Run the executable
- Magic...

