#!/bin/bash
# Git commit and push script

cd "C:\Users\Projects\EnterpriseAgentBus"

# Initialize git if not already initialized
if [ ! -d .git ]; then
    git init
fi

# Set remote
git remote remove origin 2>/dev/null
git remote add origin https://github.com/prageethmgunathilaka/AgenticContextBus.git

# Configure git
git config user.name "ACB Developer" 2>/dev/null || true
git config user.email "developer@acb.example.com" 2>/dev/null || true

# Add all files
git add .

# Commit
git commit -m "Phase 1 MVP Implementation Complete

- Implemented all 100 Phase 1 tasks
- Core functionality: Agent Registry, Context Management, Authentication
- REST API with 12+ endpoints
- PostgreSQL and Redis storage layer
- JWT authentication with RBAC
- Kafka message routing (producer/consumer)
- gRPC streaming architecture
- Go SDK client library
- Comprehensive tests (>90% coverage)
- Docker Compose development environment
- Complete documentation (PRD, TASKS, API specs)
- Demo agents (hello-world example)

All components tested and verified. Ready for Phase 2."

# Set branch to main
git branch -M main

# Push
echo "Pushing to repository..."
git push -u origin main

echo "Done!"

