{{ title: Nova - Quickstart }}

{{ include-block: doc.html markdown="true" }}

# Quickstart

To get started with Nova, follow these steps:

1. Install Nova:

```bash
go install github.com/xlc-dev/nova@latest
```

2. Create a new project and follow the prompts:

```bash
nova new myproject
```

3. Run the project:

```bash
cd myproject
go build # Or run `make` if you enabled the Makefile in the setup process
./myproject # Or ./myproject api to run the API server at localhost:8080
```

That's it! Now you got a Nova project up and running at `http://localhost:8080`.

An example project is available at [github.com/xlc-dev/nova-example](https://github.com/xlc-dev/novaexample) to showcase the power of Nova.

## Next steps

Here are some links to help you get started with Nova:

- [Scaffolding](./scaffolding.html) - Generate new projects fast from pre-built templates
- [CLI](./cli.html) - Learn about Nova's powerful CLI tooling
- [Router](./router.html) - Use Nova's simple and powerful router to handle HTTP requests based on [net/http](https://pkg.go.dev/net/http)
- [Database Migrations](./database_migrations.html) - Manage your database migrations
- [OpenAPI](./openapi.html) - Generate an OpenAPI specification for your API
- [Middleware](./middleware.html) - Use middleware to add functionality to your Nova API

{{ endinclude }}
