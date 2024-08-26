FROM  python:3.11

# Set working directory
WORKDIR /app

# Copy requirements file
COPY requirements.txt ./

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY src .

# Expose port
EXPOSE 5000

# Run the application
CMD ["python", "chat_http_service.py"]
