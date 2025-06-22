GATOR - Go CLI ApplicationGATOR is a command-line interface (CLI) application built with Go for managing users, posts, and a dynamic feed system. It uses SQL for data persistence and sqlc for type-safe database interactions, focusing on a robust and easy-to-use terminal tool.Table of ContentsFeaturesProject StructureGetting StartedCLI CommandsContributingLicenseFeaturesUser Management: CLI-based creation, retrieval, and management of user accounts.Feed & Post System: CLI interaction for managing user follows, posts, and personalized feeds.SQL Database: Reliable data storage using SQL.Type-Safe Queries: sqlc integration for secure, type-safe Go-database interactions.Modular Design: Clean, maintainable codebase.Project StructureGATOR/
├── internal/
│   ├── config/             # Application configurations
│   └── database/           # Database connection & generated Go code (from sqlc)
├── sql/
│   ├── queries/            # Raw SQL queries
│   └── schema/             # Database schema definitions
├── commands.go             # CLI command definitions
├── go.mod                  # Go module dependencies
├── handler_user.go         # User-related business logic
└── main.go                 # Application entry point
Getting StartedPrerequisitesGo (1.18+)PostgreSQLsqlcInstallationClone: git clone https://github.com/bilalhachim/gator.git && cd gatorDependencies: go mod tidyDatabase Setup: Create gator_db and apply schema:psql -U your_user -d gator_db -f sql/schema/001_users.sql
Generate sqlc code: sqlc generateRunning the ApplicationEnvironment Variable: Set DATABASE_URL (e.g., in a .env file):DATABASE_URL="postgres://user:password@localhost:5432/gator_db?sslmode=disable"
Build (optional): go build -o gatorRun Commands:Using executable: ./gator [command] [arguments]Without build: go run main.go [command] [arguments]CLI Commands(This section is a placeholder. Fill this with your actual commands and examples.)gator users create --name "John Doe" --email "john@example.com": Create user.gator posts new --user-id 123 --content "Hello GATOR!": Create post.gator feeds fetch --user-id 456: Fetch user's feed.gator help: Display commands.ContributingContributions are welcome!Fork, branch (feature/your-feature), make changes.Commit (git commit -m 'feat: Add feature X').Push (git push origin feature/your-feature).Open a Pull Request.LicenseThis project is licensed under the MIT License. See the LICENSE file.
