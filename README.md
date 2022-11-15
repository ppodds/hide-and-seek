# Hide And Seek

## Introduction

This is the socket programming homework of Computer Network course at National Central University.
Hide And Seek is a simple party game. Ghost (ball) need to catch other players to win. Player need to run away from the ghost to survive.

## Features

- Multiple clients
- Support multiple games play in the same time
- Non blocking socket
- Game (multi-threading and async IO)
- Server (multi-threading)

## Game Protocol

### TCP RPC

- Request
  - Header (Five bytes)
    - First byte - TCP RPC ID
    - Four bytes - Content length
  - Data
    - Protobuf
    - Binary encoded data
- Response
  - Header (Four bytes)
    - Content length
  - Data
    - Protobuf
    - Binary encoded data

### UDP RPC / Broadcast

- Request
  - Header (Five bytes)
    - First byte - UDP RPC ID
    - Four bytes - Content length
  - Data
    - Protobuf
    - Binary encoded data
- Response / Broadcast
  - Data
    - Protobuf
    - Binary encoded data

## Screenshots

![Game screenshot 1](docs/screenshots/1.jpg)
![Game screenshot 2](docs/screenshots/2.jpg)
![Game screenshot 3](docs/screenshots/3.jpg)
![Game screenshot 4](docs/screenshots/4.jpg)
![Game screenshot 5](docs/screenshots/5.jpg)
![Game screenshot 6](docs/screenshots/6.jpg)
