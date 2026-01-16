# Chess App

![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)
![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20Windows%20%7C%20macOS-lightgrey?style=for-the-badge)
![Status](https://img.shields.io/badge/Status-Active-success?style=for-the-badge)

A high-performance, dual-interface Chess application engineered in Go. Designed for reliability, speed, and clean architecture, this project offers both a robust Command Line Interface (CLI) and a modern Web Server mode with a RESTful API.

---

## Table of Contents

![Overview](https://img.shields.io/badge/-Overview-gray?style=flat-square&logo=read-the-docs&logoColor=white)
![Features](https://img.shields.io/badge/-Features-gray?style=flat-square&logo=starship&logoColor=white)
![Installation](https://img.shields.io/badge/-Installation-gray?style=flat-square&logo=linux&logoColor=white)
![Usage](https://img.shields.io/badge/-Usage-gray?style=flat-square&logo=gnu-bash&logoColor=white)
![API](https://img.shields.io/badge/-API-gray?style=flat-square&logo=postman&logoColor=white)
![Contributing](https://img.shields.io/badge/-Contributing-gray?style=flat-square&logo=github&logoColor=white)

---

## Features

![Dual Mode](https://img.shields.io/badge/Interface-CLI%20%26%20Web-blueviolet?style=flat-square)
![Mechanics](https://img.shields.io/badge/Engine-Legal%20Move%20Validation-orange?style=flat-square)
![State](https://img.shields.io/badge/Game_State-Check%2FCheckmate%2FStalemate-red?style=flat-square)

- **Dual-Mode Operation**: Seamlessly switch between a retro-style terminal interface and a modern web-based experience.
- **Robust Game Engine**: Fully implemented chess rules including castling, en passant, and pawn promotion.
- **Legal Move Generation**: Advanced algorithms ensure only valid moves are permitted, preventing illegal play.
- **Game State Detection**: Instant checking for Check, Checkmate, and Stalemate conditions.
- **Session Management**: The web server supports multiple concurrent game sessions, isolated by unique Session IDs.
- **RESTful API**: Clean JSON API for integration with other frontends or engines.
- **Cross-Platform**: Runs natively on Linux, Windows, and macOS.

---

## Installation

### Prerequisites

- **Go 1.25** or higher installed on your machine.

### Clone the Repository

```bash
git clone https://github.com/Waleed-Ahmad-dev/Chess-app.git
cd Chess-app
```

### Install Dependencies

```bash
go mod download
```

---

## Usage

### Command Line Interface (CLI)

Experience chess in its purest form directly in your terminal.

```bash
go run cmd/chess/main.go
```

**Controls**:

- Enter moves in algebraic notation (e.g., `e2e4`).
- Type `exit` or `quit` to close the application.

### Web Server Mode

Launch the web server to play via a browser or interact with the API.

```bash
go run cmd/chess/main.go web
```

**Access**:

- Open your browser and navigate to: `http://localhost:8080`

### WebAssembly (Browser) Mode

Run the chess engine entirely in the browser using WebAssembly.

> **Note**: If you see a `could not import syscall/js` error in your IDE, this is normal. That package is only available when compiling with `GOOS=js GOARCH=wasm`. The build commands below handle this for you.

**Build and Run**:

```bash
make serve-wasm
```

Then visit `http://localhost:8080` to play.

## API Documentation

The application exposes a set of REST endpoints for external interaction.

### New Game

**POST** `/api/new-game`
Creates a new game session.

- **Response**: `{"sessionId": "...", "state": {...}}`

### Game State

**GET** `/api/game-state?sessionId=<ID>`
Retrieves the current board state.

- **Response**: JSON object containing board representation, legal moves, and turn info.

### Make Move

**POST** `/api/make-move`
Execute a move on the board.

- **Body**: `{"sessionId": "...", "move": "e2e4"}`
- **Response**: Updated game state or error message.

---

## Project Structure

```text
Chess-app/
├── cmd/
│   └── chess/
│       └── main.go       # Entry point (CLI & Web)
├── internal/
│   ├── game/
│   │   ├── board.go      # Board representation & rendering
│   │   ├── checks.go     # Checkmate logic
│   │   ├── logic.go      # Move validation
│   │   └── state.go      # Game state management
│   └── server/
│       ├── server.go     # HTTP Server & API handlers
│       └── html.go       # Embedded frontend assets
├── go.mod                # Module definition
├── LICENSE               # MIT License
└── README.md             # Project documentation
```

---

## Contributing

We welcome contributions from the community. Please adhere to the following guidelines to maintain project quality.

### Contribution Rules

1.  **Fork & Clone**: Fork the repository to your account and clone it locally.
2.  **Branching**: Create a strictly named branch for your feature or fix.
    - `feature/new-mechanic`
    - `fix/castling-bug`
    - `docs/update-readme`
3.  **Code Style**:
    - All code must be formatted using `gofmt`.
    - Ensure variable names are descriptive and follow Go conventions (CamelCase).
    - Exported functions must have comment documentation.
4.  **Testing**:
    - Write unit tests for any new game logic in `internal/game`.
    - Ensure all existing tests pass before submitting.
5.  **Pull Request**:
    - Submit a Pull Request (PR) to the `main` branch.
    - Provide a detailed description of your changes.
    - Reference any related issues.

### Code of Conduct

- Maintain a professional and respectful tone in all communications.
- Harassment or abuse of any kind will not be tolerated.

---

## License

Distributed under the **MIT License**. See `LICENSE` for more information.

![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)

Copyright (c) 2026 Waleed Ahmad
