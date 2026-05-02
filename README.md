# Modular Hexagonal Golang Project

Modular Hexagonal Golang Project <https://docs.gofiber.io/recipes/hexagonal/>

## Modular Hexagonal Architecture with golang

The entire system adheres to SOLID principles and Clean Architecture.

The system features both an API and a WEB interface, which are completely isolated from one another.

By fully isolating the REST API and the WEB interface—which is built using JavaScript-based frameworks (such as React.js and Vue.js)—they are able to utilize shared services through Dependency Injection.

If you wish to run only the REST API or the WEB interface, please modify the `APP` parameter located in the `.env` file.

Two distinct endpoints were implemented—one for the API and one for the web interface—resulting in a completely clean architecture developed using shared libraries.

Within the system, every component has been decoupled through the use of Dependency Injection.

This framework contains two modules and library :

## Architecture

- `domain`: pure business objects
- `ports`: interfaces
- `app`: use-case / service layer
- `adapters`
- `adapters/inbound`: HTTP/Fiber handlers
- `adapters/outbound`: PostgreSQL, Redis, MongoDB implementations
- `http`
- `model`
- `viewmodel`
- `transport`
- `helpers`
- `platform`

Both modules have a hexagonal structure and contain the following folders:

- `user`
- `posts`

## Technologies

- **Hexagonal Architecture**: domain, ports, use cases, adapters, repositories, HTTP layers
- **Fiber v3**: routing, middleware, rendering
- **Pongo2**: via `github.com/gofiber/template/django/v3`
- **Bootstrap 5**: [css engine and template](https://github.com/StartBootstrap/startbootstrap-sb-admin)
- **CSRF**: Fiber CSRF middleware + session store
- **Flash Messages**: session-based
- **PostgreSQL**: primary data source
- **Redis**:  caching and session storage
- **MongoDB**: audit logs
- **Gorm**: ORM for Go
- **Validator v9**: form validation
- **jQuery DataTables and Pagination**: jQuery datatables and pagination
- **Pagination**: pagination

## Docker Compose

Redis Commander, Adminer, and Mongo Express have been included in the Docker container.

## Execution

```bash
docker compose up -d
go run ./cmd
./air 

```

## Endpoints

### Posts

- `POST /api/posts`
- `GET /api/posts`
- `GET /api/posts/:id`
- `GET /api/posts/user/:user_id`

### User

- `POST /api/users`
- `GET /api/users`
- `GET /api/users/:id`

### web users

- `POST /web/list/list_users_with_pagination`
- `GET /web/list/normal_users"`
- `GET /web/list/datatable`
- `GET /web/user/create`
- `GET /web/user/:id`
- `GET /users/:id/edit`

[@credits for bootstrap.css template](https://github.com/StartBootstrap/startbootstrap-sb-admin)

```
├─ cmd/
│     └─ main.go
├─ internal/
│  ├─ modules/
│  │  ├─ posts/
│  │  │  ├─ adapters/
│  │  │  │  ├─ inbound/
│  │  │  │  │  └─ http/
│  │  │  │  └─ outbound/
│  │  │  │     ├─ mongodb/
│  │  │  │     ├─ postgres/
│  │  │  │     └─ redis/
│  │  │  ├─ app/
│  │  │  │  └─ service.go
│  │  │  ├─ domain/
│  │  │  │  └─ post.go
│  │  │  └─ ports/
│  │  │     └─ ports.go
│  │  └─ user/
│  │     ├─ adapters/
│  │     │  ├─ inbound/
│  │     │  │  └─ http/
│  │     │  └─ outbound/
│  │     │     ├─ mongodb/
│  │     │     ├─ postgres/
│  │     │     └─ redis/
│  │     ├─ app/
│  │     │  └─ service.go
│  │     ├─ domain/
│  │     │  ├─ user.go
│  │     │  └─ userJson.go
│  │     └─ ports/
│  │        └─ ports.go
│  ├─ platform/
│  │  ├─ cache/
│  │  │  └─ redis/
│  │  │     └─ redis.go
│  │  ├─ config/
│  │  │  └─ config.go
│  │  ├─ database/
│  │  │  ├─ mongodb/
│  │  │  │  └─ mongodb.go
│  │  │  └─ postgres/
│  │  │     └─ postgres.go
│  │  ├─ helpers/
│  │  │  ├─ stnccollection/
│  │  │  ├─ stncdatetime/
│  │  │  ├─ stnchelper/
│  │  │  ├─ stncsession/
│  │  │  └─ stncupload/
│  │  └─ id/
│  │     └─ id.go
│  └─ transport/
│     ├─ api/
│     │  ├─ app.go
│     ├─ common/
│     │  └─ common.go
│     └─ web/
│        └─ app.go

```
