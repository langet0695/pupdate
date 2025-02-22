# Pupdate

*There is only one corner of the universe you can be certain of improving, and that's your own self.*

*- Aldous Huxley*

![Dog Image](./tmp/daily_dog.jpg)

## Overview

The goal of this project is simple, to spread joy by delivering puppy images to users digital door daily. 

Subscription is simple. 
- To subscribe email `pupdate.today+subscribe@gmail.com` to opt in to daily pupdates. 
- To opt out send an email to `pupdate.today+unsubscribe@gmail.com`.

At the moment this is the full extent of the interface available to end users.

The application is powered by a small Dockerized GO service run through a that stores data to disk through a mounted volume. The reason for storing data to disk is to keep the system lightweight and cheap to run.

Scheduled operations (e.g. backups and pupdates) are orchestrated through Cron jobs that run shell scripts inside of the host running the app.

User triggered activites, e.g. subscription and unsubscription, are handled through emails sent to a centralized GMail inbox. This inbox leverages filters and Apps Scripts trigger user management activites.

## Deployment

### Configure GMail interface
1. Log in or create the [GMail](https://accounts.google.com/) account that will serve as the frontend of the service. 
2. Inside of the Gmail account configure three filters
    1. Apply the label `MessagesToHandle` to all messages sent to `<your-email>`+subscribe@gmail.com
    2. Apply the label `HandledMessages` to all messages sent to `<your-email>`+unsubscribe@gmail.com
    3. Apply the label `Backups` to all messages sent from `<your-email>`@gmail.com with a `Subscribers Save` as the subject 
3. Create an [app password](https://support.google.com/mail/answer/185833?hl=en#zippy=%2Cwhy-you-may-need-an-app-password) for the pupdate app. Secure this information in your password manager. It will be needed later.
4. Finally create a deployment of xxx in Google Apps Script. This will allow your GMail account to send user requests to your service
    - Navigate to [Apps Script](https://script.google.com/) in your Google Account
    - Create a `New Project` and copy the code from `setup/apps_script_handler.js` into the project.
    - Update the constants at the top of the project with project specific information. Key constants to update are:
        - `SUBSCRIBE`
        - `UNSUBSCRIBE`
        - `DOMAIN` *This is the domain that you will use will configure your nginx server to route requests from. If you don't already have one then go purchase a domain for this project. I like using [namecheap.com](https://www.namecheap.com/)*
        - `API_PASSWORD` *Create a secure passcode here and save it. This will control access to your api and will be needed when deploying the docker container*
    - Now `Deploy` the app in Apps Script
    - Once deployed create a Trigger in apps script to run the `Gmail_Trigger` function hourly.

### Deploy Docker
Depending on the system one deploys to some modifications may be necessary. This guide shows how to deploy to most systems and how I modified it to deploy to a Pi Zero-W.

1. Install [docker](https://docs.docker.com/get-started/get-docker/) 
2. Clone this repo to your local machine
3. Create a `.env` file in the root dir of the cloned repo. In it add the following environment variables
    - MAIL_USER=`<your-email>`@gmail.com
    - MAIL_PASSWORD=`<your-email-app-password>`
    - API_USER=admin
    - API_PASSWORD=`<your-api-password>`
    - JWT_SECRET=`<create-a-custom-jwt-secret-string>`
    - SUBSCRIPTIONS_PATH=/tmp/subscriptions.json
    - DOG_IMG_PATH=/tmp/daily_dog.jpg
4. From the root of the project run `docker build -t pupdate no-cache .` to build the image!
    - NOTE: Depending on the infrastructure you are building the image for you may need to make some modifications to this step. For instance I deployed Pupdate to a Raspberyy Pi Zero-W which is on the ARMv6 architecure and has relatively low specs. As a result I used the following process to build the image:
        1. On a machine with enough memory I ran `docker build -t pupdate:armv6 --platform linux/arm/v6 --no-cache .`
        2. Next I saved the image to a `.tar` with `docker save pupdate:armv6 > ~/Desktop/pi-img-v6.tar `
        3. Once the `.tar` was built I copied the image to my host with `scp <source> <destination-remote>`
        4. Finally I loaded the image to my host with `docker load -i pi-img-v6.tar`
5. The last step is to run the image! The image requires a environment variables to be passed in from the `.env` and that the local `./tmp` folder be mounted as a volume. To do this run the following command: `docker run -d -p 8080:8080 -v ~/pupdate/tmp:/app/tmp --env-file .env pupdate`
    - NOTE: Again for the Pi Zero-W I used a modified command that I ran on the remote host: `docker run -d -p 8080:8080 -v ~/pupdate/tmp:/app/tmp --env-file .env --platform linux/arm/v6 pupdate:armv6`

### Configure Cron Jobs
1. 

### Set up nginx server

## Testing

## TODO's
- Describe Objective of repo and why the project exists. Aka update readme
- Explain available http commands in a table
- Outline design decisions around data handling
- Explain how to deploy
- Add app script code to repo
- Add shell scripts
- Copy app script snippet and explain how to add / configure gmail
- Get multi stage deploy working with envs

#DOCKER
TO build multistage image use following command
`docker build -t pupdate:multistage -f Dockerfile.multistage .`
TO run use
`docker run -d -p 8080:8080 pupdate:multistage`
To run with volume 
`docker run -d -p 8080:8080 -v ~/pupdate_data:/ext_data <image>`

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

to build for pi-zero-w use following
- set image base tag to golang:1.24rc3-alpine3.20
- specify platform in build `docker build -t pupdate:armv6 --platform linux/arm/v6 --no-cache .`
- `docker save pupdate:armv6 > ~/Desktop/pi-img-v6.tar `
- `scp <source> <destination-remote>`
- `docker load -i pi-img-v6.tar`
- `docker run -d -p 8080:8080 -v ~/pupdate/tmp:/app/tmp --env-file .env --platform linux/arm/v6 pupdate:armv6`

hosting config
- set up nginx
- port forward dynamic dns 
- configuring for namecheap domain