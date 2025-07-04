#!/bin/bash

# Qafka Development Setup Script
# This script initializes the entire Qafka development environment

set -e

echo "🚀 Setting up Qafka - Universal Queue Management Platform"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check prerequisites
check_prerequisites() {
    echo -e "${YELLOW}📋 Checking prerequisites...${NC}"
    
    command -v go >/dev/null 2>&1 || { echo -e "${RED}❌ Go is required but not installed. Please install Go 1.21+${NC}"; exit 1; }
    command -v node >/dev/null 2>&1 || { echo -e "${RED}❌ Node.js is required but not installed. Please install Node.js 18+${NC}"; exit 1; }
    command -v docker >/dev/null 2>&1 || { echo -e "${RED}❌ Docker is required but not installed.${NC}"; exit 1; }
    command -v docker-compose >/dev/null 2>&1 || { echo -e "${RED}❌ Docker Compose is required but not installed.${NC}"; exit 1; }
    
    echo -e "${GREEN}✅ All prerequisites found${NC}"
}

# Setup Go backend
setup_backend() {
    echo -e "${YELLOW}🏗️  Setting up Go backend...${NC}"
    cd backend
    
    # Initialize Go module if not exists
    if [ ! -f "go.mod" ]; then
        go mod init github.com/shravan20/qafka
    fi
    
    # Install dependencies
    go get github.com/go-fuego/fuego
    go get github.com/uptrace/bun
    go get github.com/uptrace/bun/driver/pgdriver
    go get github.com/uptrace/bun/extra/bundebug
    go get github.com/joho/godotenv
    go get github.com/prometheus/client_golang/prometheus
    go get github.com/prometheus/client_golang/prometheus/promhttp
    go get github.com/swaggo/swag/cmd/swag
    go get github.com/swaggo/http-swagger
    go get github.com/swaggo/files
    
    # Install swag for API docs
    go install github.com/swaggo/swag/cmd/swag@latest
    
    cd ..
    echo -e "${GREEN}✅ Backend setup complete${NC}"
}

# Setup React frontend
setup_frontend() {
    echo -e "${YELLOW}🎨 Setting up React frontend...${NC}"
    
    # Create Vite React TypeScript project
    if [ ! -d "frontend/src" ]; then
        npm create vite@latest frontend -- --template react-ts
        cd frontend
        
        # Install additional dependencies
        npm install @tanstack/react-query
        npm install @radix-ui/react-slot
        npm install class-variance-authority
        npm install clsx
        npm install tailwind-merge
        npm install lucide-react
        npm install axios
        
        # Install dev dependencies
        npm install -D tailwindcss postcss autoprefixer
        npm install -D @types/node
        
        # Initialize Tailwind CSS
        npx tailwindcss init -p
        
        cd ..
    else
        cd frontend
        npm install
        cd ..
    fi
    
    echo -e "${GREEN}✅ Frontend setup complete${NC}"
}

# Setup documentation
setup_docs() {
    echo -e "${YELLOW}📚 Setting up documentation...${NC}"
    cd docs
    
    if [ ! -f "package.json" ]; then
        npm init -y
        npm install -g @mintlify/cli
    fi
    
    cd ..
    echo -e "${GREEN}✅ Documentation setup complete${NC}"
}

# Create environment files
create_env_files() {
    echo -e "${YELLOW}⚙️  Creating environment files...${NC}"
    
    # Backend .env
    if [ ! -f "backend/.env" ]; then
        cat > backend/.env << EOF
# Database Configuration
DATABASE_URL=postgres://qafka_user:qafka_password@localhost:5432/qafka?sslmode=disable

# API Configuration
API_PORT=8080
ENVIRONMENT=development

# CORS Configuration
CORS_ORIGINS=http://localhost:5173

# Monitoring
PROMETHEUS_ENABLED=true
PROMETHEUS_PORT=2112
EOF
    fi
    
    # Frontend .env
    if [ ! -f "frontend/.env" ]; then
        cat > frontend/.env << EOF
VITE_API_URL=http://localhost:8080
VITE_APP_NAME=Qafka
VITE_APP_VERSION=1.0.0
EOF
    fi
    
    echo -e "${GREEN}✅ Environment files created${NC}"
}

# Main execution
main() {
    check_prerequisites
    setup_backend
    setup_frontend
    setup_docs
    create_env_files
    
    echo -e "${GREEN}🎉 Qafka setup complete!${NC}"
    echo -e "${YELLOW}Next steps:${NC}"
    echo "1. Start the database: docker-compose -f docker-compose.dev.yml up -d"
    echo "2. Run the backend: cd backend && go run cmd/main.go"
    echo "3. Run the frontend: cd frontend && npm run dev"
    echo "4. Visit http://localhost:5173 to see your app!"
}

# Run main function
main
