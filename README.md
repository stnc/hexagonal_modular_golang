# Modular Hexagonal Go Project

Modular Hexagonal Golang Project <https://docs.gofiber.io/recipes/hexagonal/>


## modular hexagolanl Architecture with golang

The entire system adheres to SOLID principles and Clean Architecture.

The system features both an API and a web interface, which are completely isolated from one another.

By fully isolating the REST API and the web interface—which is built using JavaScript-based frameworks (such as React.js and Vue.js)—they are able to utilize shared services through Dependency Injection.

Two distinct endpoints were implemented—one for the API and one for the web interface—resulting in a completely clean architecture developed using shared libraries.

Within the system, every component has been decoupled through the use of Dependency Injection.

This framework contains two modules and library :

- `user`
- `posts`
- `pongo2 template engine - like django`
- `pagination`
- `jquery datatable and pagination`
- `fiber v3`
- `repository`
- `redis cache`
- `mongo db`
- `Gorm`

Both modules have a hexagonal structure and contain the following folders:

- `adapters`
- `app`
- `domain`

- `ports`

Technologies:

- Fiber
- GORM
- PostgreSQL
- Redis
- MongoDB

## Architecture

- `domain`: pure business objects
- `ports`: interfaces
- `app`: use-case / service layer
- `adapters/inbound`: HTTP/Fiber handlers
- `adapters/outbound`: PostgreSQL, Redis, MongoDB implementations

## Execution

For WEB

``` bash
docker compose up -d
go run ./cmd/web
./air 

```

For API

```bash
docker compose up -d
go run ./cmd/api
# ./air -c .airapi.toml

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
- `GET /web/usersDatatable`
- `GET /web//user/create`
- `GET /web/user/:id`

[@credits for bootstrap.css template](https://github.com/StartBootstrap/startbootstrap-sb-admin)
