# pupdate
TODO:
- Describe Objective of repo and why the project exists
- Explain available http commands in a table
- Outline design decisions around data handling
- Explain how to deploy
- Build a method to send volume to pupdate inbox for storage
- Copy app script snippet and explain how to add / configure gmail
- Build auth
- Set up nginx reverse proxy
- Get multi stage deploy working with envs

#DOCKER
TO build multistage image use following command
`docker build -t pupdate:multistage -f Dockerfile.multistage .`
TO run use
`docker run -d -p 8080:8080 pupdate:multistage`
To run with volume 
`docker run -d -p 8080:8080 -v ~/pupdate_data:/ext_data 00c7792f22a7`

# Deploying with an external volume
Copy your personal list of subs in json to the subscriptions.json folder
navigage to the root dir fo this repo ~/<path-to-pupdate>
build your image with `docker build -t pupdate:<personal_tag> -f Dockerfile .`
execute `docker run -d -p 8080:8080 -v ~/pupdate/src/startup:/app/tmp pupdate:volrouter`