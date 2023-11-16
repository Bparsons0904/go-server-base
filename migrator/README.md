# Migrations Setup

This document describes the process of setting up and using Gorm migrations. We are using `gormigrate` for handling our database schema migrations.

## Setup

1. **Code Snippets**: We're using VSCode code snippets for creating migration templates. To set this up, follow these steps:

   - Go to `File > Preferences > User Snippets` in VSCode.
   - Select `New Global Snippets file` and give it a relevant name (e.g., `gorm.snippets.json`).
   - Paste the JSON content for the snippets in this file. An example of such content might look like the one we discussed above for creating a new migration file or a new migration with transactions.

2. **Migration Files**: Migrations reside in individual Go files within the `migrations` directory (you can name the directory as you see fit). Each file should register itself to the `RegisteredMigrations` array. The migration files should follow the naming convention like `YYYYMMDDHHMMSS_description.go`.

3. **Migration Registry**: The `registry.go` in the `migrations` directory contains the `RegisteredMigrations` array. All migrations register themselves to this array using the `init` function in their respective files.

## Creating Migrations

To create a new migration, you can use the command `go run ./migrator/migrator.go -create some_description` to generate a new migration file. This command will create a new file in the `migrations` directory with the current timestamp and the provided description. You can then use the `gormmigration` or `gormmigrationtrans` snippet to generate a template for additional migrations.

## Running Migrations

Once your migrations are set up, you can run them using the command-line interface. The exact commands may vary based on your project's setup. Below are sample commands:

- To **apply** migrations: `go run ./migrator/migrator.go -up`

- To **rollback** migrations: `go run ./migrator/migrator.go -down`


These commands should be run from the server directory.

## Notes

- Make sure your migration IDs are unique and chronologically incremental.
- Always test your migrations locally (up and down) before applying them to production databases.
- All changes should be reflected in the models as well.
- The `down` or rollback functionality should reverse the changes made by the `up` migration. Make sure to implement it carefully to avoid unwanted data loss.

## References
- [Gormigrate]("https://github.com/go-gormigrate/gormigrate")
- [Gorm]("https://gorm.io/docs/migration.html")
****