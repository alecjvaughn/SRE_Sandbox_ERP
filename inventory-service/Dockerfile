# Stage 1: Build dependencies
FROM python:3.11-slim-bookworm AS builder

WORKDIR /app
COPY requirements.txt .

# Install dependencies into local user directory
RUN pip install --user --no-cache-dir -r requirements.txt

# Stage 2: Final runtime image
FROM python:3.11-slim-bookworm

WORKDIR /app

# Copy the local user dependencies
COPY --from=builder /root/.local /root/.local

# Ensure the local bin is on PATH
ENV PATH=/root/.local/bin:$PATH

# Copy application source code
COPY . .

# Expose metrics and health check ports
EXPOSE 8000 8081

# Set the default entrypoint
CMD ["python", "app.py"]
