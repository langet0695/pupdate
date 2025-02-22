# Pupdate :dog:

*There is only one corner of the universe you can be certain of improving, and that's your own self.*

*- Aldous Huxley*

<img src="./tmp/daily_dog.jpg" width="25%" height="25%">

## Overview

The goal of this project is simple, to spread joy by delivering puppy images to users digital door daily. 

Subscription is simple. 
- To subscribe email `pupdate.today+subscribe@gmail.com` to opt in to daily pupdates. 
- To opt out send an email to `pupdate.today+unsubscribe@gmail.com`.

At the moment this is the full extent of the interface available to end users.

The application is powered by a small Dockerized GO service run through a that stores data to disk through a mounted volume. The reason for storing data to disk is to keep the system lightweight and cheap to run.

Scheduled operations (e.g. backups and pupdates) are orchestrated through Cron jobs that run shell scripts inside of the host running the app.

User triggered activites, e.g. subscription and unsubscription, are handled through emails sent to a centralized GMail inbox. This inbox leverages filters and Apps Scripts trigger user management activites.

## TODO's / Future Improvements
- Get multi stage deploy working with .env file
- Improve error response handling on by adding messges to return to the user based on error type
- Add automated unit testing
- Find a way to enable automated triggers from inbound emails. Will allow immediate response to user activites
- Build customization interface so that users can define parameters that customize their pupdate experience

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
1. Modify the `./setup/CronScripts/pupdate-save.sh` and `./setup/CronScripts/pupdate-trigger.sh` scripts to use your api-password.
2. In your console run `crontab -e` to edit the cron file that you can use to schedule your scripts. In it create two cron routines to trigger the scripts daily. Below are the schedules I used:
    - 45 15 * * * `<path-to-script>`/pupdate-trigger.sh
    - 45 16 * * * `<path-to-script>`/pupdate-save.sh
3. Save the cron file and your routines will now run on a schedule!

### Set up a reverse proxy
The last step is to set up a reverse proxy to route requests from you domain to your local host. I used a nginx server, but there are other great options available. I won't describe this in detail here as there is a lot of documentation on how to do this already on the web. If you haven't done this before then I recommend checking out the [nginx docs](https://nginx.org/en/docs/beginners_guide.html).

The only note I have here is that it's 2025. Configure ssl for your reverse proxy and only configure connections through port 443.


### Completion

Once the above is configured you've deployed pupdate! 

Feel free to re-configure and alter the code to build a custome experience for your users. 

Good luck and happy trails :paw_prints:
