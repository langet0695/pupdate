# pupdate
TODO:
- Describe Objective of repo and why the project exists
- Explain available http commands in a table
- Outline design decisions around data handling
- Explain how to deploy
- Copy app script snippet and explain how to add / configure gmail
- Set up nginx reverse proxy
- Get multi stage deploy working with envs

- Build a method to send volume to pupdate inbox for storage

- Set up authorization for email reg 
1. Wrap createsubscriber with jwt auth middle ware and change parameter fetch from url param to a jwt key pair
2. Create a new route that creates a token and sends it in a formatted url that new subs need to ack

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
execute `docker run -d -p 8080:8080 -v ~/pupdate/<path_to_data>:/app/tmp pupdate:<personal_tag>`
<!-- execute `docker run -d -p 8080:8080 -v ~/pupdate/src/startup:/app/tmp pupdate:volrouter` -->

make sure to set your environment user and password to something that will then be used to generate your jwt
to get your jwt use a command as follows
`curl -X "POST" http://admin:aPassword@localhost:8080/createToken`
// TODO IS make sure this is configured to only work with https