{{ title: Nova - Documentation }}

{{include-block: doc.html markdown="true"}}

# Nova

Nova is a flexible framework that simplifies creating both RESTful APIs and web UIs.
It extends Go's standard library with sensible defaults and helper utilities for components like routing, middleware, OpenAPI, and HTML templating, minimizing decision fatique.
Making it easier than ever to build powerful web applications in Go.

## Why Choose Nova?

Nova is designed to be a lightweight framework that has as little dependencies as possible.
It is built on top of Go's standard library, making it easy to integrate with existing Go codebases and libraries.

Together with the CLI tool, Nova provides a streamlined development experience for building web applications.

## Key Features:

- **CLI Tooling:** Integrated command-line tooling to build any CLI for your application.
- **Project Scaffolding:** Quickly generate new projects with a sensible default structure using `nova new`.
- **Database Migrations:** Manage database migrations effortlessly with the `nova` binary.
- **Streamlined REST APIs:** Simplified routing, request handling, and response generation.
- **OpenAPI:** Built-in support for request validation and OpenAPI (Swagger) spec generation.
- **Middleware Support:** Easily add and manage middleware for enhanced functionality.
- **Templating Engine:** Built-in support for building HTML templates within Go files.

## Design Philosophy

Nova is designed to be lightweight and flexible. It's built on top of the standard library,
without any external dependencies except [fsnotify](https://github.com/fsnotify/fsnotify)
for file watching and database drivers for migrations and `database/sql` support.

The supported database drivers are:

- [SQLite](https://github.com/mattn/go-sqlite3)
- [PostgreSQL](https://github.com/lib/pq)
- [MySQL/MariaDB](https://github.com/go-sql-driver/mysql)

it is very easy to add a new databse driver
to nova if requested, as long as it follows the design philosophy of the framework.

The goal of the framework is to be modular and extensible, so it is designed to follow as many
standards as possible wihtin the Go community, with the goal being able to plug and play different components within the framework.

## Design Principles

- **Minimal Dependencies:** Nova is designed to have as few dependencies as possible.
- **Simplicity:** Nova is designed to be simple and easy to use. Follow pattenrs that Golang developers are already familiar with.
- **Stable:** Nova is designed to be stable and reliable. That is another benifit of only using the standard library which is known to be stable.
- **Reduce Boilerplate:** Nova is designed to reduce boilerplate code, so the developr can focus on what is the most important: the business logic.
- **Batteries Included:** Nova comes with batteries included. It comes with many components like a CLI and router, removing decision fatique, needing to write boilerplate code or installing external dependencies.
- **Undocumented features are bugs:** Undocmunted code or unwritten documentation is concidered a bug, allowing for _actual_ up to date documentation.

## Getting Started

To get started with Nova, read the [Quickstart](./quickstart.html) guide, or check out the [example project](https://github.com/xlc-dev/novaexample).

## License

This project is licensed under the [MIT License](https://github.com/xlc-dev/nova/blob/main/LICENSE).

{{endinclude}}
