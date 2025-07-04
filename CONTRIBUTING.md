# Contributing to Qafka

Thank you for your interest in contributing to Qafka! We welcome contributions from everyone.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/shravan20/qafka.git`
3. Create a feature branch: `git checkout -b feature/amazing-feature`
4. Make your changes
5. Run tests: `./scripts/test.sh`
6. Commit your changes: `git commit -m 'Add amazing feature'`
7. Push to the branch: `git push origin feature/amazing-feature`
8. Open a Pull Request

## Development Setup

### Prerequisites
- Go 1.21+
- Node.js 18+
- Docker & Docker Compose
- PostgreSQL (for local development)

### Quick Setup
```bash
# Clone and setup
git clone https://github.com/shravan20/qafka.git
cd qafka
chmod +x scripts/setup.sh
./scripts/setup.sh

# Start development environment
docker-compose -f docker-compose.dev.yml up -d
```

### Backend Development
```bash
cd backend
go mod download
go run cmd/main.go
```

### Frontend Development
```bash
cd frontend
npm install
npm run dev
```

## Code Style

### Go Code
- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions
- Write tests for new functionality

### TypeScript/React Code
- Use TypeScript for all new code
- Follow React best practices
- Use functional components with hooks
- Write unit tests for components

### Commits
- Use conventional commit format: `type(scope): description`
- Types: feat, fix, docs, style, refactor, test, chore

## Testing

### Backend Tests
```bash
cd backend
go test ./...
```

### Frontend Tests
```bash
cd frontend
npm test
```

### Integration Tests
```bash
./scripts/test-integration.sh
```

## Pull Request Process

1. Ensure all tests pass
2. Update documentation if needed
3. Add entry to CHANGELOG.md
4. Get review from maintainers
5. Squash commits before merge

## Code of Conduct

This project follows the [Contributor Covenant](CODE_OF_CONDUCT.md). Please read it before contributing.

## Questions?

Feel free to open an issue or reach out to the maintainers!
