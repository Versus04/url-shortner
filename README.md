# 🔗 URL Shortener

A lightweight, fast URL shortener service built with Go. No external dependencies required—just the Go standard library.

## ✨ Features

- **Fast & Lightweight** – Pure Go with no external dependencies
- **Base62 Encoding** – Generates short, URL-safe codes using alphanumeric characters
- **Persistent Storage** – URLs are saved to disk and survive restarts
- **Thread-Safe** – Uses `sync.RWMutex` for concurrent request handling
- **Simple API** – RESTful JSON API for creating and resolving short URLs

## 🚀 Getting Started

### Prerequisites

- Go 1.18 or higher

### Installation

```bash
git clone https://github.com/yourusername/url-shortner.git
cd url-shortner
go build -o url-shortner
```

### Running the Server

```bash
./url-shortner
```

The server will start on `http://localhost:8080`

## 📖 API Reference

### Shorten a URL

**Endpoint:** `POST /shorten`

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "url": "https://example.com/very/long/url"
}
```

**Response:** `201 Created`
```json
{
  "short_code": "W7E",
  "short_url": "http://localhost:8080/W7E"
}
```

### Redirect to Original URL

**Endpoint:** `GET /r/{short_code}`

**Response:** `302 Found` – Redirects to the original URL

## 💻 Usage Examples

### Create a Short URL

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com"}'
```

**Response:**
```json
{
  "short_code": "W7E",
  "short_url": "http://localhost:8080/W7E"
}
```

### Access a Short URL

Open in browser or use curl:
```bash
curl -L http://localhost:8080/r/W7E
```

## 🗂️ Project Structure

```
url-shortner/
├── main.go        # Application entry point and all handlers
├── store.json     # Persistent storage (auto-generated)
├── go.mod         # Go module file
└── README.md
```

## ⚙️ How It Works

1. **Encoding:** Each URL is assigned a unique counter value, which is encoded using Base62 (a-z, A-Z, 0-9) to create a compact short code.

2. **Storage:** URLs and the counter are persisted to `store.json` after each new URL is created.

3. **Redirection:** When a short URL is accessed via `/r/{code}`, the server looks up the original URL and issues a `302 Found` redirect.

## 🛠️ Configuration

| Constant    | Default                   | Description                |
|-------------|---------------------------|----------------------------|
| `base_url`  | `http://localhost:8080/`  | Base URL for short links   |
| `count`     | `78992`                   | Initial counter value      |

To change these, edit the constants in `main.go`.

## 📝 License

This project is open source and available under the [MIT License](LICENSE).

## 🤝 Contributing


---

Made with ❤️ in Go
