# Double Or Nothing Dice Game

## Overview

Welcome to **Double Or Nothing Dice**, a fun, web-based dice game where players can place bets, roll dice, and test their luck. The game supports multiple players, user balances, and integrates with a backend system for real-time gameplay. Players can interact with the system through HTTP APIs for account management (register, login, etc.) and WebSocket connections for real-time gameplay.

## Setup

Follow these steps to set up the project on your local machine. *(fill in specific setup steps here)*

1. **Clone the Repository**
   - Clone this repo using Git:
     ```bash
     git clone https://github.com/Anacardo89/doubleOrNothingDice
     ```

2. **Install Dependencies**
   - Ensure you have [Go](https://golang.org/dl/) installed.
   - Install all necessary Go dependencies:
     ```bash
     go mod tidy
     ```

3. **Config**
   - Change `config.yaml_example` to `config.yaml` and change variables as needed

4. **Database Setup**
   - Ensure you have [Docker](https://www.docker.com/products/docker-desktop/) installed.
   ```bash
   docker-compose up -d
   ```
   - If you don't want to configure the email, you can just activate users manually through [pgAdmin](https://www.pgadmin.org/download/)

5. **Start the Server**
   - To start the backend server:
     ```bash
     go run ./cmd
     ```

6. **Postman Collection**
   - Import the provided Postman collection in the root directory to interact with the HTTP API for authentication, registration, and other functionalities.

## HTTP Authentication

The app provides several HTTP endpoints for managing user authentication, including:

### Register

- **POST** `/register`
  - Register a new user with the provided credentials.
  - **Payload Example**:
    ```json
    {
        "type": "register",
        "payload": {
            "username": "user",
            "email": "email@mail.com",
            "password": "password"
        }
    }
    ```

### Login

- **POST** `/login`
  - Login with registered credentials to receive a JWT token.
  - **Payload Example**:
    ```json
    {
        "type": "login",
        "payload": {
            "username": "user",
            "password": "password"
        }
    }
    ```

### Activate Account

- **POST** `/activate?token=`
  - Activate a new account by confirming the activation token.

### Forgot Password

- **POST** `/forgot-password`
  - Request a password reset link.
  - **Payload Example**:
    ```json
    {
        "type": "forgot-password",
        "payload": {
            "email": "email@mail.com"
        }
    }
    ```

### Recover Password

- **POST** `/recover-password?token=`
  - Submit a new password after following the password recovery flow.
  - **Payload Example**:
    ```json
    {
        "type": "recover-password",
        "payload": {
            "password": "password"
        }
    }
    ```

### Postman Collection

You can find the Postman collection for testing these HTTP routes in the root directory of the repo. Import the collection into Postman to test registration, login, and account management.

## WebSocket Interaction

WebSocket communication allows the game to be played in real-time. Players interact with the server via a WebSocket connection for creating a session, playing rounds, and depositing funds.

### Connecting to the WebSocket Server

Once the backend is running, connect to the WebSocket server at the following endpoint:

ws://localhost:8080/ws?token=

### WebSocket Message Examples
Below are the examples of WebSocket messages for interacting with the system. These include requests to manage the wallet, start plays, end plays, and make deposits. Each section includes both the request and the expected response.

#### 1. **Wallet Request**
To request the current balance of a user:

**Request:**
```json
{ 
    "type": "wallet", 
    "payload": { 
        "client_id": "user123" 
    } 
}
```

**Response:**
```json
{ 
    "type": "wallet_response", 
    "payload": { 
        "balance": 100
    } 
}
```

#### 2. **Play Request**
To start a new game and place a bet:

**Request:**
```json
{ 
    "type": "play", 
    "payload": { 
        "client_id": "user123", 
        "bet_amount": 50, 
        "bet_type": "even" 
    } 
}
```

**Response:**
```json
{ 
    "type": "play_response", 
    "payload": { 
        "rolled_number": 4, 
        "next_bet": 100, 
        "outcome": "win" 
    } 
}
```

#### 3. **End Play Request**
To end the current play session and receive winnings:

**Request:**
```json
{ 
    "type": "end_play", 
    "payload": { 
        "client_id": "user123" 
    }
}
```

**Response:**
```json
{ 
    "type": "end_play_response", 
    "payload": { 
        "winnings": 100, 
        "balance": 200 
    }
}
```

#### 4. **Deposit Request**
To deposit money into the user's account:

**Request:**
```json
{ 
    "type": "deposit", 
    "payload": { 
        "client_id": "user123", 
        "deposit": 50 
    } 
}
```

**Response:**
```json
{ 
    "type": "deposit_response", 
    "payload": { 
        "balance": 150 
    } 
}
```