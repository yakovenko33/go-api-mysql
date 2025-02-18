## Golang API with MySQL 

A Restful API in Golang using Gin-Gonic framework and MySQL.

## Setting up the API

Once you clone the repository, copy the `.env.example` file to `.env` and update the values as per your MySQL configuration. Make sure your MySQL server is running and have the necessary permissions and the database created with the name you have provided in the `.env` file.

Once you are done with the configuration, you can run the following command to start the API:

```sh
go run main.go
```

This will start the API on `http://localhost:3000`.


### Command migration ###

`goose create new_user_table sql` - to create SQL migration;  
`goose -dir ./database/migrations mysql "user:password@tcp(mysql:3306)/my_database?parseTime=true" up` - to execute migration;
`goose -dir ./database/migrations mysql "user:password@tcp(mysql:3306)/my_database?parseTime=true" down` - roll back migration;


### Command rebuild ###
go run main.go


### Command start CLI app ###
`go run main.go cli create-super-admin`
