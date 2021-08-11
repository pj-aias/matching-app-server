# matching-app-server
Server for https://github.com/pj-aias/matching-app .

## Usage
0. Before run, generate a password for MySQL database in any way you like (e.g. `pwgen` command, or password manager software). And write it to `/mysql_password.txt`.
1. Run `docker-compose up --build`.
2. In development environment, port `:8080` will be used for the default. Now you can access the server on `localhost:8080`.