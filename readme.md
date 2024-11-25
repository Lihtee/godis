# Godis: In-memory Cache

Godis is a simple implementation of a Redis-like in-memory cache. It provides a key-value storage system with support for strings and dictionaries, along with per-key TTL (Time-To-Live) functionality. The server is designed to handle commands over a TCP connection, offering a Telnet-like/HTTP-like API protocol.

## Features

### Implemented Features
- **Key-Value Storage**: Supports string and dictionary types.
- **Per-Key TTL**: Allows setting a TTL for each key.
- **Operations**:
  - **Get**: Retrieve the value of a key.
  - **Set**: Store a value with a key.
  - **Remove**: Delete a key from the storage.
- **Telnet-like/HTTP-like API Protocol**: Communicate with the server using a simple text-based protocol.

### Planned Features
- **Key-Value Storage**: Support for lists.
- **Golang API Client**: A client library for interacting with the server.
- **Tests**: Unit and integration tests for the server and client.
- **API Specification**: Detailed documentation of the API.
- **Deployment Manual**: Instructions for deploying the server.
- **Client Usage Examples**: Examples demonstrating how to use the client library.
- **Performance Tests**: Tests to measure and improve the server's performance.
- **Additional Operations**: Support for more complex operations like listing all keys.

## Getting Started

### Prerequisites
- Go 1.16 or later

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/godis.git
   cd godis
   ```

2. Build the server:
   ```bash
   go build -o godis-server ./server
   ```

3. Run the server:
   ```bash
   ./godis-server
   ```

The server will start listening on port 12345.

## Usage

### Connecting to the Server
You can connect to the server using a Telnet client or any TCP client that supports sending raw text commands.

### Commands
- **Get**: Retrieve a value
  ```plaintext
  get string <key>
  get dict <key> [field1 field2 ...]
  ```

- **Set**: Store a value
  ```plaintext
  set string <key> <value>
  set dict <key> <field1>:<value1> <field2>:<value2> ...
  ```

- **Delete**: Remove a key
  ```plaintext
  delete <key>
  ```

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes. Make sure to follow the project's coding standards and include tests for any new features or bug fixes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

