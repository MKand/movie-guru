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

# Use the official PostgreSQL image as a base
FROM pgvector/pgvector:pg17

# Set the working directory inside the container
WORKDIR /docker-entrypoint-initdb.d/

# Copy the init.sql script to the container
COPY init_substituted.sql /docker-entrypoint-initdb.d/

# Ensure proper file permissions
RUN chmod 644 /docker-entrypoint-initdb.d/init_substituted.sql

# Expose the PostgreSQL port
EXPOSE 5432

# Use the default PostgreSQL entrypoint script
CMD ["postgres"]
