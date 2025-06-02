# 🛒 E-Commerce Backend API (Go)

A robust, production-grade E-Commerce backend built in **Go**, following **SOLID principles**, **modular architecture**, and powered by **JWT authentication**, **Redis caching**, and **Chi router** for blazing fast HTTP routing.

---

## 🚀 Features

✅ JWT-based Authentication & Authorization  
✅ Modular Layered Architecture (Handler → Service → Repository)  
✅ SOLID Principles throughout the codebase  
✅ Redis for Caching User Sessions & Fast Lookups  
✅ Chi Router for Lightweight Routing  
✅ PostgreSQL with GORM ORM  
✅ Dockerized for Easy Deployment  
✅ Role-based Access (Admin/User)  
✅ Cart & Order Checkout with SQL Transactions  

---

## 🧱 Folder Structure

ecommerce-backend/
│
├── cmd/ # Application entrypoint
├── models/ # Data models
├── handlers/ # HTTP handlers (thin controllers)
├── services/ # Business logic layer
├── repo/ # DB access layer (GORM/SQL)
├── routes/ # Route grouping using Chi
├── middleware/ # JWT auth, admin check, logging
├── utils/ # Helper utilities (parsing, validation)
├── config/ # Redis/Postgres configuration
└── main.go # Server bootstrap



---

## 🔐 Authentication

- Uses **JWT tokens** for stateless, secure authentication.
- Middleware auto-extracts and validates tokens.
- Role-based access with `"admin"` and `"user"` flags.
- Authenticated user IDs are injected into request context.

```http
Authorization: Bearer <your_jwt_token>



📦 API Endpoints
🧑‍💼 Auth
Method	Endpoint	Description
POST	/users/register	Register new user
POST	/users/login	Login & get token

🛒 Cart (JWT Required)
Method	Endpoint	Description
POST	/cart	Add product to cart
GET	/cart	View current cart
PATCH	/cart/:id	Update cart item
DELETE	/cart/:id	Remove item from cart

📦 Orders (JWT Required)
Method	Endpoint	Description
POST	/orders	Create order from cart
GET	/orders	List orders
GET	/orders/:id	Get order by ID
PATCH	/orders/:id	Update payment/total
DELETE	/orders/:id	Delete order

📦 Products
Method	Endpoint	Role	Description
GET	/products	Public	View products
POST	/products	Admin	Add new product
PATCH	/products/:id	Admin	Update product
DELETE	/products/:id	Admin	Delete product

⚙️ Tech Stack
Go (Golang) — Language

Chi — Lightweight, idiomatic HTTP router

GORM — ORM for PostgreSQL

Redis — Caching layer

JWT — Stateless auth tokens

Docker — Containerization

🐳 Running Locally
Make sure Docker & Docker Compose are installed.

bash
Copy
Edit
git clone https://github.com/your-username/ecommerce-backend.git
cd ecommerce-backend
cp .env.example .env    # Fill in your DB, Redis, JWT_SECRET
docker-compose up --build
🔐 .env Configuration
env
Copy
Edit
PORT=8080

POSTGRES_HOST=db
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=ecommerce

REDIS_HOST=redis
REDIS_PORT=6379

JWT_SECRET=your_secret_key
📚 Principles & Design
✅ SOLID
Single Responsibility — Each layer does one thing well.

Open/Closed — Extendable services, no edits to core.

Liskov Substitution — Interfaces for repos/services.

Interface Segregation — Clean, purpose-specific contracts.

Dependency Inversion — High-level modules don’t depend on low-level ones.

🧩 Modular Layers
Handlers: Decode requests, encode responses.

Services: Business logic, validations.

Repositories: DB queries, SQL/ORM logic.

