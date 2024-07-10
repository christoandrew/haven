```markdown
# Haven Project

Haven is a comprehensive financial management system designed to provide users with detailed insights into their transactions, accounts, and financial health through a variety of data visualizations and summaries.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

Before you begin, ensure you have the following installed:
- Go (version 1.21)
- MySQL (version 5.7 or later)

### Installing

A step-by-step series of examples that tell you how to get a development environment running.

1. **Clone the repository**

```bash
git clone https://github.com/christo-andrew/haven.git
cd haven
```


2. **Set up the database**

If you're using Docker, you can set up a MySQL database using:

Remember to update the `.env`database configuration with your credentials.

```bash
docker run --name haven-db -e MYSQL_ROOT_PASSWORD=yourPassword -e MYSQL_DATABASE=haven -p 3306:3306 -d mysql:latest
```
3. **Install dependencies**

Navigate to the project directory and install the Go dependencies:

```bash
go mod tidy
```



### Running the application

To start the server, run:

```bash
go run main.go
```

The server should now be running and accessible at `http://localhost:8080/api/v1/swagger/index.html`.


## API Documentation

The API documentation is available in the `docs/swagger.json` file. You can view it using any Swagger UI viewer.

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.