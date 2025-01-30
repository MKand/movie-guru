# Use the official PostgreSQL image as a base
FROM pgvector/pgvector:pg17

# Set the working directory inside the container
WORKDIR /docker-entrypoint-initdb.d/

# Copy the init.sql script to the container
COPY init.sql /docker-entrypoint-initdb.d/

# Ensure proper file permissions
RUN chmod 644 /docker-entrypoint-initdb.d/init.sql

# Expose the PostgreSQL port
EXPOSE 5432

# Use the default PostgreSQL entrypoint script
CMD ["postgres"]
