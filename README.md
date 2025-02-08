# master-rad

# Microservices Communication Patterns

## Overview

This repository contains a set of microservices demonstrating different communication patterns in a microservices architecture. The project explores synchronous and asynchronous communication, including blocking and non-blocking approaches. It leverages technologies such as gRPC and RabbitMQ for inter-service communication.

This implementation is based on research conducted for a master's thesis on "Communication Patterns in Microservices Architecture," covering different patterns like synchronous blocking, asynchronous non-blocking, event-driven communication, and message brokering.

## Tech Stack

Programming Language: Go (Golang)

Frameworks & Libraries: gRPC, Protocol Buffers (Protobuf)

Message Broker: RabbitMQ

RPC Protocol: gRPC

Databases: PostgreSQL, MongoDB

Containerization: Docker & Docker Compose

## Microservices

This project includes the following services:

Authentication Service - Handles user authentication and authorization.

Broker Service - Acts as a message broker between services.

Frontend - Client-side application for interacting with microservices.

Listener Service - Listens to events and routes them accordingly.

Logger Service - Implements gRPC-based remote procedure calls for logging.

Mail Service - Handles email notifications asynchronously.

Order Service - Manages order creation and processing.

Payment Service - Handles payment transactions.

## Communication Patterns Demonstrated

Synchronous Blocking Communication: gRPC-based service calls with request-response patterns.

Asynchronous Non-Blocking Communication: Implemented using RabbitMQ message queues.

Event-Driven Communication: Decoupling services using message brokering.

Data Consistency Strategies: Two-phase commit (2PC) and Saga pattern.

Installation & Running Locally

## Prerequisites

Docker & Docker Compose

Protobuf Compiler

## Setup

## Clone the repository
git clone https://github.com/TeodoraDamnjanovic/your-repo.git
cd your-repo

## Start services using Docker Compose
docker-compose up --build

## Testing

gRPC calls can be tested using tools like Postman, BloomRPC, or grpcurl.

RabbitMQ management UI is available at http://localhost:15672 (default credentials: guest/guest).

MongoDB logs can be inspected via MongoDB Compass.

