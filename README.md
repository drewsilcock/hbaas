# HBaaS Server

This is a fun demo API for saying happy birthday to people for the purposes of demonstrating server containerisation and deployment.

## Software Stack

This project uses the [Go programming language](https://go.dev/) for rapid prototype development and easy containerisation.
 
For the RESTful API, we are using the [Echo](https://echo.labstack.com/) framework due to its ease-of-use, minimalism, extensability but also performance. 

For persistence, GORM is used to interact with a PostgreSQL database.

## Running

There are multiple commands available to the application for seeding the database with test values, auto-migrating the database schema and running the main server.

The main command for the server is:

```sh
./hbaas-server run-server
```

For information on all command-line options, run:

```sh
./hbaas-server --help
```

In order to connect to your local PostgreSQL database, you must specify the `POSTGRES_URL` environment variable. See the `.env.example` file for what this should look like.
