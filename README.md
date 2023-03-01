# devproxy
Quick and dirty proxy to split frontend / backend for local development, simulating an ALB doing Path based proxying.

## Why?
I was working on a project that had a frontend and backend, and I wanted to be able to run them both locally, but have the frontend proxy to the backend. I didn't want to have to run the backend on a different port, and I didn't want to have to run the frontend on a different port. I also didn't want to have to run the frontend on a different domain. I wanted to be able to run them both on the same port, and have the frontend proxy to the backend.

## Usage
To use it download the current os version from the releases page then:
- Extract the zip file
- Configure the config.yaml file
- Run the executable
- Magic...

