# Copyright 2025 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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


# Start the production server (adjust the command if necessary)
CMD ["npm", "start", "--", "--host", "0.0.0.0"]
