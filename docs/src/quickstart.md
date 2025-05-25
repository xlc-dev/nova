# Quickstart

To get started with Nova, follow these steps:

1. Install Nova:

```sh
go install github.com/xlc-dev/nova
```

2. Create a new project and follow the prompts:

```sh
nova new myproject
```

3. Run the project:

```sh
cd myproject
go build # Or run `make` if you enabled the Makefile in the setup process
./myproject # Or ./myproject api to run the API server at localhost:8080
```

That's it! Now you got a Nova project up and running.

## Next steps

Here are some links to help you get started with Nova:

- 🛠️ [CLI](./cli.md) - Learn about Nova's powerful CLI tooling
- 🏗️ [Scaffolding](./scaffolding.md) - Generate new projects fast from pre-built templates
- 🔀 [Router](./router.md) - Use Nova's simple and powerful router to handle HTTP requests based on [net/http](https://pkg.go.dev/net/http)
- 🗃️ [Database Migrations](./database_migrations.md) - Manage your database migrations
- 🔐 [OpenAPI](./openapi.md) - Generate an OpenAPI specification for your API
- 🗳️ [Middleware](./middleware.md) - Use middleware to add functionality to your Nova API
