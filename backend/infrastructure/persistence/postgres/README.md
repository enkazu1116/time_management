# Local PostgreSQL and sqlc

## Start PostgreSQL

```sh
docker compose up -d postgres
```

The local database is initialized from `infrastructure/persistence/postgres/db/schema.sql`.

```text
host: localhost
port: 5432
database: time_management
user: time_management
password: time_management
sslmode: disable
```

Connection string:

```text
postgres://time_management:time_management@localhost:5432/time_management?sslmode=disable
```

## Generate sqlc code

Install sqlc locally:

```sh
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

Then generate typed query code:

```sh
sqlc generate
```

Generated code is written to:

```text
infrastructure/persistence/postgres/sqlc
```

## Query ownership

SQL files are split by repository model under `infrastructure/persistence/postgres/db/queries`.

- `users.sql`
- `roles.sql`
- `work_types.sql`
- `work_logs.sql`
- `break_logs.sql`
- `weekly_tasks.sql`
- `daily_tasks.sql`
- `daily_logs.sql`
- `weekly_logs.sql`
- `achievements.sql`

Repository interfaces live under each domain's `interfaces/repositories` directory.
PostgreSQL implementations should be added under `infrastructure/persistence/postgres` and call the generated `sqlc` package.
