## Running Ollama container

### Setup

#### Components
- Windows 11
- WSL (Ubuntu-24.04)
- Docker
- Docker Compose

#### Prerequisites
- Install Docker and Docker Compose on WSL (Ubuntu-24.04)
    ```bash
    sudo apt-get install docker docker-compose
    sudo usermod -aG docker $USER
    ```

### Running the container

Used the VSCode Docker extension to run the container from within Windsurf on WSL (Ubuntu-24.04).

#### Llam3.2:1b

#### Pull the model

Lookup the model_id on https://ollama.com/library > `https://ollama.com/library/llama3.2:1b`

`curl http://localhost:8008/api/pull -d '{ "model": "llama3.2:1b" }'`

#### Generate text

`curl http://localhost:8008/api/generate -d '{ "model": "llama3.2:1b", "prompt": "What model are you?" }'`

#### Granite3.1-dense:2b

#### Pull the model

Lookup the model_id on https://ollama.com/library > `https://ollama.com/library/granite3.1-dense:2b`

`curl http://localhost:8008/api/pull -d '{ "model": "granite3.1-dense:2b" }'`

#### Generate text

`curl http://localhost:8008/api/generate -d '{ "model": "granite3.1-dense:2b", "prompt": "What model are you? Are you IBM Granite?" }'`