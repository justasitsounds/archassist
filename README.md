# ArchAssist

ArchAssist is short for Architecture Assistant

## Getting Started

### Prerequisites

 - go 1.4+
 - [aws-icons-for-plantuml](https://github.com/awslabs/aws-icons-for-plantuml.git) github repository cloned to a locally accessible file location
 - Docker
 - (optional) [PlantUML extension for vs-code](https://github.com/qjebbs/vscode-plantuml.git)

### Installation

```bash
go install
```


## Usage

1. Set environment variable `AWSICONDIST` to match the filepath to the `/dist` directory location in the local aws-icons-for-plantuml github repository

```bash
export AWSICONDIST=/Users/justasitsounds/development/github.com/awslabs/aws-icons-for-plantuml/dist
```

2. if you wish to edit and make use of the live(ish) preview option available with the vs-code plantuml extension, start the local preview server:

```bash
make start-puml
```

This will download the `plantuml/plantuml-server` docker image from DockerHub and start a new container, listening on port 8080

3. invoke archassist from the command line to start the architecture diagram wizard