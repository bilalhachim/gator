GATOR - Go CLI Application
GATOR is a command-line interface (CLI) application built with Go for managing users,
posts, and a dynamic feed system. It uses SQL for data persistence and sqlc for
type-safe database interactions, focusing on a robust and easy-to-use terminal tool.
Table of Contents
● Features
● Project Structure
● Geing Started
● CLI Commands
● Contributing
● License
Features
● User Management: CLI-based creation, retrieval, and management of user
accounts.
● Feed & Post System: CLI interaction for managing user follows, posts, and
personalized feeds.
● SQL Database: Reliable data storage using SQL.
● Type-Safe Queries: sqlc integration for secure, type-safe Go-database
interactions.
● Modular Design: Clean, maintainable codebase.
Project Structure
GATOR/
├── internal/
│ ├── cong/ # Application congurations
│ └── database/ # Database connection & generated Go code (from sqlc)
├── sql/
│ ├── queries/ # Raw SQL queries
│ └── schema/ # Database schema denitions
├── commands.go # CLI command denitions
├── go.mod # Go module dependencies
├── handler_user.go # User-related business logic
└── main.go # Application entry point
Geing Started
Prerequisites
● Go (1.18+)
● PostgreSQL
● sqlc
Installation
1. Clone: git clone hps://github.com/bilalhachim/gator.git && cd gator
2. Dependencies: go mod tidy
3. Database Setup: Create gator_db and apply schema:
psql -U your_user -d gator_db -f sql/schema/001_users.sql
4. Generate sqlc code: sqlc generate
Running the Application
1. Environment Variable: Set DATABASE_URL (e.g., in a .env le):
DATABASE_URL="postgres://user:password@localhost:5432/gator_db?sslmode=d
isable"
2. Build (optional): go build -o gator
3. Run Commands:
○ Using executable: ./gator [command] [arguments]
○ Without build: go run main.go [command] [arguments]
CLI Commands
(This section is a placeholder. Fill this with your actual commands and examples.)
● gator users create --name "John Doe" --email "john@example.com": Create user.
● gator posts new --user-id 123 --content "Hello GATOR!": Create post.
● gator feeds fetch --user-id 456: Fetch user's feed.
● gator help: Display commands.
Contributing
Contributions are welcome!
1. Fork, branch (feature/your-feature), make changes.
2. Commit (git commit -m 'feat: Add feature X').
3. Push (git push origin feature/your-feature).
4. Open a Pull Request.
License
This project is licensed under the MIT License. See the LICENSE le
