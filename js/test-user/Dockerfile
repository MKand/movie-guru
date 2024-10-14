FROM node:lts-alpine

# set the working direction
WORKDIR /app

# install app dependencies
COPY package.json ./
COPY package-lock.json ./

# Fix permissions for packages
# RUN npm config set unsafe-perm true

RUN npm install
RUN npm install -g serve

# Bundle app source
COPY . ./


# Start the production server (adjust the command if necessary)
CMD [ "npm", "run", "dev", "--", "--host",  "0.0.0.0" ]
