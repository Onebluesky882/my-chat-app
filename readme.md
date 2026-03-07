# Go + ScyllaDB Real-Time Chat Server

High-performance real-time chat backend built with Go, WebSocket, and ScyllaDB.

Designed for scalable messaging systems capable of handling millions of users.

---

## Architecture

Client
│
▼
Load Balancer
│
▼
WebSocket Server (Go)
│
▼
Pub/Sub (Redis / NATS)
│
▼
ScyllaDB Cluster

---

## Tech Stack

- Go (WebSocket server)
- ScyllaDB (message storage)
- Redis / NATS (Pub/Sub for message broadcast)
- WebSocket protocol for real-time communication

---

## Features

- Real-time chat via WebSocket
- Message persistence
- Chat rooms
- Horizontal scaling
- Low latency messaging
- Scalable architecture

---

## Project Structure
