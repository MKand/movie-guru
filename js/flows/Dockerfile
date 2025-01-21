FROM node:lts-alpine

# Set the working directory
WORKDIR /app

# Install app dependencies
COPY package.json ./
COPY package-lock.json ./

# Install dependencies
RUN npm install

# Bundle app source
COPY . ./

# Build the application
RUN npm run build

# Start the production server
CMD ["npm", "start", "--", "--host", "0.0.0.0"]
