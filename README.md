# Nimage

This is a Go-based image server that converts and serves images in WebP format. It provides caching capabilities and supports runtime configuration through a JSON file.

## Features

- Converts JPEG and PNG images to WebP format on-the-fly
- Caches converted images for improved performance
- Configurable quality and server settings for WebP conversion
- Cache clearing functionality with key-based authentication
- Debug logging for easier troubleshooting

## Prerequisites

- Go 1.16 or higher

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/configurable-image-server.git
   cd configurable-image-server
   ```

2. Build the project:
   ```
   make
   ```

## Configuration

Create a `config.json` file in the project root directory with the following structure:

```json
{
  "quality": 90,
  "cacheFolder": "./cache",
  "cacheClearKey": "yourSecretKey",
  "port": "8080",
  "debug": true
}
```

- `quality`: WebP conversion quality (0-100)
- `cacheFolder`: Directory to store cached WebP images
- `cacheClearKey`: Secret key for cache clearing functionality
- `port`: Port number for the server to listen on
- `debug`: Enable or disable debug logging

## Usage

1. Start the server:
   ```
   go run main.go -config config.json
   ```

2. Access images through the server:
   ```
   http://localhost:8080/path/to/your/image.jpg
   ```
   The server will convert the image to WebP format, cache it, and serve it.

3. Clear the cache:
   ```
   http://localhost:8080/clearcache?key=yourSecretKey
   ```
   Replace `yourSecretKey` with the value specified in your `config.json`.

## Debug Logging

When `debugLogging` is set to `true` in the configuration, the server will output detailed log messages to help with troubleshooting and monitoring. These messages include information about request handling, file operations, and cache management.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the AGPLv3 License - see the [LICENSE](LICENSE) file for details.