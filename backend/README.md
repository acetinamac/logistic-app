# Backend (Go + GORM + PostgreSQL)

Este backend implementa una API para gestión de órdenes con autenticación JWT (sin expiración) y roles (cliente y admin) siguiendo una estructura inspirada en Clean Architecture.

## Requisitos
- Go 1.24
- Docker / Docker Compose

## Estructura
- cmd/api: punto de entrada
- internal/domain: entidades (Order, User)
- internal/usecase: casos de uso (OrderService)
- internal/repository: implementación GORM
- internal/infra/db: conexión a Postgres
- internal/delivery/http: handlers HTTP y autenticación

## Endpoints (MVP)
- POST /api/users => registrar usuario (body: {email, password, role?})
- DELETE /api/users/{id} => eliminar usuario (admin o el propio usuario)
- POST /api/login => body: {"user_id": number, "role": "client"|"admin"} devuelve token JWT
- POST /api/orders (cliente|admin)
- GET /api/orders (cliente => solo propias; admin => si ?all=1, todas)
- GET /api/admin/orders (admin)
- PATCH /api/admin/orders/{id}/status (admin)

Reglas de negocio: tamaño del paquete según peso (S ≤5kg, M ≤15kg, L ≤25kg). Si peso>25kg => error solicitando convenio especial.

## Ejecutar con Docker
```
docker compose up --build
```
La API quedará en http://localhost:8080

## Variables de entorno relevantes
- POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB
- JWT_SECRET

## Justificación PostgreSQL
PostgreSQL ofrece integridad ACID, tipos avanzados (jsonb), extensiones geoespaciales (PostGIS) ideales para logística, y es excelente con GORM por su madurez.

## Postman
En la carpeta postman/ se incluye una colección con ejemplos (login cliente, crear orden, listar, actualizar status).
