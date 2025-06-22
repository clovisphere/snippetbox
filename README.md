# 📇 Snippet-Box

Share your text snippets effortlessly! 😊

> This is a hands-on project created while working through the excellent [Let's Go](https://lets-go.alexedwards.net/) book by Alex Edwards.

## 🛠️ Getting Started

Before you can use snippetbox, ensure you have the following prerequisites installed on your system:

- **Go**: version 1.24 or higher ([Installation Guide](https://go.dev/dl/))

## ⚙️ Installation and Usage

Follow these steps to install and run snippetbox:

1. Clone the repository:
   ```sh
   git clone https://github.com/clovisphere/snippetbox.git
   cd snippetbox
   ```

2. Install dependencies:
   ```sh
   go mod download  # Install dependencies
   ```

3. Build and run the application:
   ```sh
   go run main.go   # Build and run the application
   ```

Or just use the Makefile:
   ```sh
   make run  # or just make
   ```

By default, the application will listen on port `4000`. You can use a different port
by passing the `addr` (or `PORT` when using `make`) flag to the `run` command,
i.e. `go run main.go addr=":6969"` or `make run PORT=6969`.
