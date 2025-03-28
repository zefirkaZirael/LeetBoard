# LeetBoard
 An anonymous imageboard inspired by early internet forums. Users can post text and images, with avatars from the Rick and Morty universe. No user registration is required. Posts and comments are temporary, automatically disappearing after a set time based on activity, fostering dynamic discussions. S3 used to upload Pictures
## To run the project, follow these steps:

Run the Database Server and Triple-s Server:

First, set up the necessary servers by running:
docker-compose -f scripts/docker-compose.yml up --build

This will launch both the database and Triple-s storage server.

Build the Go Application:

Once the servers are running, you can build the application using the following Go command:

```go build -o 1337b04rd ./cmd/1337b04rd```

Run Tests:

To run all tests in the project, use:

```go test -v ./...```

Or to run only the tests in the tests folder, use:

```go test -v ./tests/```

To check the test coverage, use:

```go test -cover ./...```
## Architecture
This project follows Hexagonal Architecture, where services interact with domain interfaces rather than concrete implementations. Key Changes in Codebase:

Before:

Services directly depended on repositories (e.g., CommentRepository, PostRepository).
type CommentService struct { Repo *repository.CommentRepository PostRepo *repository.PostRepository }

After:

Services now depend on domain interfaces, which enhances flexibility and testability.

type CommentService struct {
    Repo     domain.CommentRepoInt
    PostRepo domain.PostRepo
}
Why This Is Better:

Flexibility: The service can work with any repository implementation, making it easy to swap databases or mock repositories for testing.
Testability: The service can easily be tested with mock repositories.
Separation of Concerns: The service logic is decoupled from the details of data storage, improving maintainability.

## Testing
### Business Logic Testing
Various business logic components of the project are tested, such as:

Creating posts (TestCreatePost)
Creating comments (TestCreateComment)
Validating sessions (TestCreateSession, TestIsValidSession)
### Integration Testing
Tests are written for database interactions and storage services, including:

Database repository tests (e.g., post_repository_db_test.go)
Storage service tests (e.g., storage_service_test.go)
Mock data is used to simulate the database and storage operations, ensuring that the business logic is tested in isolation. Example Test:

func TestCreatePost(t *testing.T) {
    // Test logic for creating a post
}
### Mock Testing
Mock repositories are used for testing purposes. For example, MockStorageRepo simulates storage behavior using an in-memory map, which allows for more controlled testing scenarios.

## Log Management
The application uses a logging system to track important actions and errors. Logs are categorized by severity:

Info: General information.
Warning: Warnings about possible issues (e.g., exceeding avatar limit).
Error: Critical errors that need attention.
Log messages are generated for operations such as session validation and error handling. Features


 
## Key Features
Hexagonal Architecture: The application is structured with Hexagonal Architecture to ensure separation of concerns and maintainability. The core logic is independent of external systems like databases, web frameworks, and external APIs.

Anonymous User Identification: Users are tracked through browser sessions with cookies. Upon the first visit, they are assigned unique avatars and names from the Rick and Morty API.

Posts and Comments: Users can create posts (with optional images) and comment on them. Comments can be replies to posts or other comments. Posts are deleted after 10 minutes without comments or 15 minutes after the last comment.

Image Storage: Images attached to posts and comments are stored in S3-compatible storage (MinIO or another service). The project avoids saving images locally.

Session Management: Users are assigned avatars and names via cookies, and sessions expire after 1 week.

Logging: The application logs significant events and errors using Go's log/slog package.

Testing: The project includes unit tests covering core functionalities like post and comment creation, session management, and image storage integration.
Session Expiry: Automatically delete expired sessions.
Data is stored in PostgreSQL, and Triple-s is used for object storage.
## Functionality
Session Management: HTTP cookies are used to manage user sessions with a 1-week expiration time. Each session is associated with a unique avatar and name from the Rick and Morty API.

Post Deletion: Posts with no comments are deleted 10 minutes after creation. Posts with comments are deleted 15 minutes after the last comment.

Comments: Users can comment on posts and other comments. Each comment includes the user’s avatar and a unique comment ID. Users can reply to specific posts or comments by clicking on their IDs.

## Database and Storage
PostgreSQL is used to store posts, comments, user sessions, and related metadata.

S3-Compatible Storage (like MinIO) is used for storing images attached to posts and comments.

