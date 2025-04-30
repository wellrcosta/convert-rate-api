# Convert Rate API

A fast and lightweight currency conversion API written in Go, with caching, rate limiting, metrics, and structured logging to Grafana Loki.

## ✅ Features

- 🔄 Currency conversion via [ExchangeRate-API](https://www.exchangerate-api.com/)
- ⚡ Redis caching with configurable TTL
- 📊 Prometheus metrics exposed via `/metrics`
- 🔒 Rate limiting per IP (configurable)
- 📦 Dockerized with `docker-compose`
- 📈 Logs sent directly to Loki in structured format (info/warn/error)

## 🚀 Getting Started

### 1. Clone the project

```bash
git clone https://github.com/wellrcosta/convert-rate-api.git
cd convert-rate-api
```

### 2. Copy and configure `.env`

```bash
cp .env.example .env
```

Edit `.env` and set:

- `EXCHANGE_API_KEY` → your key from [https://www.exchangerate-api.com](https://www.exchangerate-api.com)
- `LOKI_URL` → the URL to your running Loki server (e.g., `http://localhost:3100/loki/api/v1/push`)

### 3. Build and run

```bash
docker-compose up --build
```

API will be available at: [http://localhost:3000](http://localhost:3000)

---

## 📤 Example Usage

### Convert 59 USD to BRL

```bash
curl "http://localhost:3000/convert?from=USD&to=BRL&amount=59"
```

Response:

```json
{
  "from": "USD",
  "to": "BRL",
  "amount": "59.00",
  "rate": "5.123456",
  "convertedAmount": "302.28",
  "cached": false
}
```

### Metrics endpoint

```bash
curl http://localhost:3000/metrics
```

---

## 🧪 Tests

To be implemented (future)

---

## 📂 Project Structure

```bash
.
├── cmd/                  # Entry point
├── internal/             # App logic
│   ├── cache/            # Redis cache layer
│   ├── client/           # Exchange rate API integration
│   ├── handler/          # HTTP handlers
│   ├── metrics/          # Prometheus setup
│   ├── rate/             # Rate limiting
│   └── logger/           # Loki logger
├── go.mod
├── .env.example
├── Dockerfile
├── docker-compose.yml
└── README.md
```

---

## 🛠 Requirements

- Go 1.24+
- Docker + Docker Compose
- Redis
- Grafana Loki (optional but recommended)

---

## 👨‍💻 Author

Made by [@wellrcosta](https://github.com/wellrcosta)
