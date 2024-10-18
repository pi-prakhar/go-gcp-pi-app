### Desired structure

# Service Management Commands

Simple guide for managing Docker services using Make commands.

## Commands & Examples

### Start Services
```bash
# Local environment
make up-local FLAGS="-d --build"  # -d flag runs in background

# Debug environment
make up-debug FLAGS="-d"

# Production environment
make up-production FLAGS="-d --build"
```

### Stop Services
```bash
make down  # Stops all running services
```

### Clean Data related to Services
```bash
make clean  # Stops all running services
```

### Build Images
```bash
# Basic build
make build-images

# Build with specific settings
make build-images \
  DOCKER_ACCOUNT=myaccount \
  SERVICES="auth server user" \
  TAGS="latest v1.0.0" \
  DOCKER_FILE_PATH=./docker \
  BUILD_CONTEXT=.
```

### Push Images
```bash
# Push to HUB
make push-images
```

### View Logs
```bash
make logs  # Shows logs from all services
```

### Restart Services
```bash
make restart  # Equivalent to make down + up
```