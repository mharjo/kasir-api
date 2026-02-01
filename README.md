# Kasir API (Go net/http)

Kasir API sederhana menggunakan **Go (net/http tanpa framework)**.  
Project ini dibuat untuk mempelajari dasar Go dengan membangun REST API dari nol.

## Fitur
- Health check
- CRUD Produk
- CRUD Category
- JSON request & response
- In-memory storage (tanpa database)
- Siap di-deploy (Railway / Zeabur)

---

## 🛠️ Requirement
- Go sudah terinstall  
  Cek dengan:
  ```bash
  go version

---

## 🚀 Cara Menjalankan (Local)

1. Clone repository:

   ```bash
   git clone <LINK_GITHUB_KAMU>
   cd kasir-api
   ```

2. Jalankan aplikasi:

   ```bash
   go run main.go
   ```

3. Server akan berjalan di:

   ```
   http://localhost:8080
   ```

---

## 🩺 Health Check

**GET** `/health`

```bash
curl http://localhost:8080/health
```

Response:

```json
{
  "status": "OK",
  "message": "API Running"
}
```

---

# 📦 Produk API

## Get all produk

**GET** `/api/produk`

```bash
curl http://localhost:8080/api/produk
```

---

## Create produk

**POST** `/api/produk`

```bash
curl -X POST http://localhost:8080/api/produk \
  -H "Content-Type: application/json" \
  -d '{
    "nama": "Kopi Kapal Api",
    "harga": 2500,
    "stok": 200
  }'
```

---

## Get produk by ID

**GET** `/api/produk/{id}`

```bash
curl http://localhost:8080/api/produk/1
```

---

## Update produk

**PUT** `/api/produk/{id}`

```bash
curl -X PUT http://localhost:8080/api/produk/1 \
  -H "Content-Type: application/json" \
  -d '{
    "nama": "Indomie Goreng Jumbo",
    "harga": 4000,
    "stok": 150
  }'
```

---

## Delete produk

**DELETE** `/api/produk/{id}`

```bash
curl -X DELETE http://localhost:8080/api/produk/1
```

---

# 🗂️ Category API (Task Session 1)

## Get all categories

**GET** `/api/categories`

```bash
curl http://localhost:8080/api/categories
```

---

## Create category

**POST** `/api/categories`

```bash
curl -X POST http://localhost:8080/api/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Minuman",
    "description": "Produk yang bisa diminum"
  }'
```

---

## Get category by ID

**GET** `/api/categories/{id}`

```bash
curl http://localhost:8080/api/categories/1
```

---

## Update category

**PUT** `/api/categories/{id}`

```bash
curl -X PUT http://localhost:8080/api/categories/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Snack",
    "description": "Makanan ringan"
  }'
```

---

## Delete category

**DELETE** `/api/categories/{id}`

```bash
curl -X DELETE http://localhost:8080/api/categories/1
```

---

## 🏗️ Build Binary

Build executable tanpa runtime tambahan:

```bash
go build -o kasir-api
```

Versi lebih kecil:

```bash
go build -ldflags="-s -w" -o kasir-api
```

---

## ☁️ Deployment

Aplikasi ini bisa di-deploy ke:

* Railway
* Zeabur

Pastikan endpoint berikut bisa diakses publik:

* `/health`
* `/api/categories`

Contoh:

```bash
curl https://<URL-DEPLOY>/health
```

---

## Versions
- Task 1: tag `task-1` (CRUD basic, in-memory)
- Task 2: tag `task-2` (layered architecture, database, config)


