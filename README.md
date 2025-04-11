# Go Web Application

A modern web application built with Go and Templ templating engine.

## Prerequisites

- Go 1.21 or later
- Make

## Getting Started

1. Clone the repository:
```bash
git clone <your-repository-url>
cd <repository-name>
```

2. Install development dependencies:
```bash
make install-deps
```

This will install:
- [Air](https://github.com/cosmtrek/air) - Live reload for Go applications
- [Templ](https://github.com/a-h/templ) - HTML templating for Go

## Development

### Available Commands

- **Run the application in development mode**:
  ```bash
  make dev
  ```
  This will start the application with hot reload enabled.

- **Generate Templ files**:
  ```bash
  make generate
  ```
  Run this when you modify any `.templ` files.

- **Build the application**:
  ```bash
  make build
  ```
  This will generate Templ files and build the Go binary.

- **Clean build artifacts**:
  ```bash
  make clean
  ```

- **Run tests**:
  ```bash
  make test
  ```

For a full list of available commands, run:
```bash
make help
```

## Project Structure

```
.
├── .air.toml          # Air configuration for hot reload
├── Makefile           # Build and development commands
└── tmp/               # Build artifacts (gitignored)
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 