# Scaffolding

**scaffolding** in Nova automates project bootstrapping. It provides:

- **Multiple Layouts:** Choose between Minimal or Structured template.
- **Name-Aware Paths:** Template filenames and directories can include `{{.Name}}` to personalize paths.
- **Content Templating:** File contents are processed as Go `text/template`, with access to project metadata.
- **Database Adapter Injection:** Select between `sqlite`, `postgres`, or `mysql`; the correct Go import path is injected.
- **Verbose Mode:** See directory- and file-creation logs for debugging.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Core Concepts](#core-concepts)

---

## Getting Started

Have the Nova binary installed and in your `$PATH`, and then run:

```bash
nova new project-name
```

Follow the prompts to customize your project. This will create a `project-name` directory populated by your chosen template,
with Go import paths and filenames adjusted to your project name.

## Core Concepts

When you run `nova new`, Nova creates a new project directory based on a template. The options are:

- **Minimal:** A bare-bones project with a single `main.go` file.
- **Structured:** A more complete project with a `main.go` file, a `cmd` directory, and a `db` directory.

When you select a template, you get the following options:

- Add a database, with the following options:

  - **sqlite:** Use the `modernc.org/sqlite` package.
  - **postgres:** Use the `github.com/lib/pq` package.
  - **mysql:** Use the `github.com/go-sql-driver/mysql` package.

    - Add a `.env` file with the following environment variables when you select a database:
      - `DATABASE_URL`: The database connection string.

- Initialize a git repository.

  - Add a `.gitignore` if you said yes to git

- Add a Make file for making it easier to develop.
