# Migrator

A simple database migration tool written in Go for managing and applying database schema changes in a controlled, versioned manner. It supports running migrations forward and backward to evolve database schemas safely.

## Installation

Clone the repository and build:

```bash
git clone https://github.com/iandaly/migrator.git
cd migrator
go build -o migrator
```

Or install directly:

```bash
go install github.com/iandaly/migrator@latest
```

## Usage

### Initialize a new project

```bash
migrator init
```

This creates a `migrator.yaml` configuration file.

### Create a new migration

```bash
migrator make create_users_table
```

This creates a new folder in the `migrations` directory with `up.sql` and `down.sql` files.

### Run migrations

```bash
migrator migrate
```

### List pending migrations

```bash
migrator pending
```

### Rollback migrations

```bash
migrator rollback
```

## Configuration

Edit `migrator.yaml`:

```yaml
url: postgresql://user:password@localhost:5432/database
folder: migrations
driver: postgresql
```

## Supported Databases

- PostgreSQL

## License

MIT