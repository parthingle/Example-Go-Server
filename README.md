## GO Server

### What is it:

This code is for a server to store and provide access to trades made a client. This Readme is truncated to protect the indentity of the source and thus only contains the bare minimum information.  

### API Spec:

Paste the src/swagger.yml file at https://editor.swagger.io/ to see API spec this project is implementing.

## How to Run:

0. Install prerequisites: 

- `golang v1.13.6`
- `docker >=19`

1. Clone repo: `git clone https://github.com/parthingle/Example-Go-Server.git`

2. Auto build to run with docker: `make`. This cleans the environment, runs all tests, generates a binary if tests succeed, and deploys on an alpine container with `localhost:8080` exposed. 

2a. Run within the terminal: `cd src && go run main.go`

