# Social Media API in Go

This is a RESTful API for a social media application, developed in Go. The project was created as part of the "API com Go" course taught by Otavio Augusto Gallego.

The API allows users to register, log in, create posts, follow other users, and interact with posts.

🔗 **Original Course Repository:** [DevBook API by OtavioGallego](https://github.com/OtavioGallego/DevBook/tree/master/api)

## ✨ Features

-   **Authentication**: Secure login system using JSON Web Tokens (JWT).
-   **User Management**:
    -   Create, search, update, and delete users.
    -   System to follow and unfollow other users.
    -   Fetch followers and the users a specific user is following.
    -   Securely update a user's password.
-   **Post Management**:
    -   Create, search, update, and delete posts.
    -   Fetch posts from a specific user and the feed of posts from followed users.
    -   System to like and unlike posts.

## 🚀 Getting Started

Follow the instructions below to run the project in your local environment.

### Prerequisites

-   [Go](https://golang.org/) (version 1.15 or higher)
-   [MySQL](https://www.mysql.com/)

### Installation

1.  **Clone the repository:**
    ```sh
    git clone [https://github.com/daniel-q-reis/api-social-media.git](https://github.com/daniel-q-reis/api-social-media.git)
    cd api-social-media
    ```

2.  **Install dependencies:**
    ```sh
    go mod tidy
    ```

3.  **Set up environment variables:**
    Create a file named `.env` in the root of the project and add the following variables, adjusting the values as needed:
    ```env
    API_PORT=5000
    DB_USUARIO=your_mysql_user
    DB_SENHA=your_mysql_password
    DB_NOME=gopher
    SECRET_KEY=your_secret_key_here_for_jwt
    ```
    *The `SECRET_KEY` can be any random and secure string.*

4.  **Set up the Database:**
    Run the SQL script to create the necessary tables in your MySQL database.
    ```sql
    CREATE DATABASE IF NOT EXISTS gopher;
    USE gopher;

    CREATE TABLE usuarios(
        id int auto_increment primary key,
        nome varchar(50) not null,
        nick varchar(50) not null unique,
        email varchar(50) not null unique,
        senha varchar(100) not null,
        criadoEm timestamp default current_timestamp()
    ) ENGINE=INNODB;

    CREATE TABLE seguidores(
        usuario_id int not null,
        FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE,
        seguidor_id int not null,
        FOREIGN KEY (seguidor_id) REFERENCES usuarios(id) ON DELETE CASCADE,
        primary key(usuario_id, seguidor_id)
    ) ENGINE=INNODB;

    CREATE TABLE publicacoes(
        id int auto_increment primary key,
        titulo varchar(50) not null,
        conteudo varchar(300) not null,
        autor_id int not null,
        FOREIGN KEY (autor_id) REFERENCES usuarios(id) ON DELETE CASCADE,
        curtidas int default 0,
        criadoEm timestamp default current_timestamp()
    ) ENGINE=INNODB;
    ```

5.  **Run the API:**
    ```sh
    go run main.go
    ```
    The server will be running at `http://localhost:5000` (or the port you defined in `.env`).

## API Endpoints

Here is the list of available endpoints in the API. All routes that require authentication need a `Bearer` token in the `Authorization` header.

---

### Login

-   `POST /login` - Authenticates a user and returns a JWT token.

### Users

-   `POST /usuarios` - Creates a new user.
-   `GET /usuarios` - Fetches all users (with filtering by name or nick).
-   `GET /usuarios/{usuarioId}` - Fetches a specific user by ID.
-   `PUT /usuarios/{usuarioId}` - Updates a user's information (requires authentication).
-   `DELETE /usuarios/{usuarioId}` - Deletes a user (requires authentication).
-   `POST /usuarios/{usuarioId}/seguir` - Allows the authenticated user to follow another user (requires authentication).
-   `POST /usuarios/{usuarioId}/parar-de-seguir` - Allows the authenticated user to unfollow another user (requires authentication).
-   `GET /usuarios/{usuarioId}/seguidores` - Returns a user's followers.
-   `GET /usuarios/{usuarioId}/seguindo` - Returns the users a specific user is following.
-   `POST /usuarios/{usuarioId}/atualizar-senha` - Updates the authenticated user's password (requires authentication).

### Posts

-   `POST /publicacoes` - Creates a new post (requires authentication).
-   `GET /publicacoes` - Fetches the posts from the authenticated user's feed (requires authentication).
-   `GET /publicacoes/{publicacaoId}` - Fetches a specific post by ID.
-   `PUT /publicacoes/{publicacaoId}` - Updates a post (requires authentication).
-   `DELETE /publicacoes/{publicacaoId}` - Deletes a post (requires authentication).
-   `GET /usuarios/{usuarioId}/publicacoes` - Fetches all posts from a specific user.
-   `POST /publicacoes/{publicacaoId}/curtir` - Adds a like to a post.
-   `POST /publicacoes/{publicacaoId}/descurtir` - Removes a like from a post.

---