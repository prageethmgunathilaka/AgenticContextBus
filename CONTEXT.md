# ACB Project Context for AI Assistants

## Project Overview
**Agentic Context Bus (ACB)** is an enterprise-grade communication framework for distributed AI agents to share context across remote locations. This document provides complete context for any AI assistant working on this project.

## Key Requirements
- **Scale**: 10-50 agents initially, designed to scale to 1000s
- **Performance**: P99 latency <50ms (small messages), 100+ MB/s streaming throughput
- **Security**: TLS 1.3, JWT authentication, RBAC, comprehensive audit logging
- **Deployment**: Production-ready from day one with Docker/Kubernetes support
- **Developer**: Learning Go while building (provide helpful comments and explanations)

## Core Problem Being Solved
Allow autonomous agents (AI, RPA, automation, microservices) running in different locations to:
- Share context (state, memory, knowledge, events)
- Handle any size data: small messages to GB-scale streams
- Discover and communicate with each other
- Maintain security and observability
- Scale from prototype to enterprise production

[Full content - see CONTEXT.md file for complete 823 lines]