# pupdate
TODO:
- Describe Objective of repo and why the project exists
- Explain available http commands in a table
- Outline design decisions around data handling
- Explain how to deploy

#DOCKER
TO build multistage image use following command
`docker build -t pupdate:multistage -f Dockerfile.multistage .`
TO run use
`docker run -d -p 8080:8080 pupdate:multistage`