# Go-Lang Worker for uploading base64 images in s3

## Description
This Go-based worker is designed to efficiently process a database of 10 million user records, each containing a base64-encoded profile picture. The worker converts the base64 text into images and seamlessly uploads them to S3 object storage.

To handle this large-scale operation, the worker utilizes multi-threaded Go routines, ensuring optimal performance and scalability while processing the user records in parallel.

## Technologies
- go 1.18
- Mysql 8
- s3
- Docker
- Docker Compose
- Redis
- Go Air

## Installation (Docker based)

### Step 1:
Clone the **"go-lang-worker-for-s3-uploader"** repo

```sh
# clone via https mode
git clone https://github.com/jasimjuwel/go-lang-worker-for-s3-uploader.git
```

### Step 2:
Create environment files copying from sample example files and make necessary changes.

```sh
cp .env.example .env
```
### Step 4:

- Create a database `testDB` and import QA database from QA environment

    ```sql
    CREATE DATABASE `testDB` COLLATE 'utf8mb4_unicode_ci';
    ```

### Step 5: build the application image

#### with docker compose:
- Go to root direcoty

- Run `docker compose build` to build the container

### Step 6: Run the container

#### [For Local development- ONE TIME ONLY]

- go to root directory
- Run ` docker compose run --rm app sh` to enter image then run `compose install`.
- Run ` docker compose exec -it app bash` to enter container.


#### Run the container with docker compose:
- Run the container in detachable mode with this command: `docker compose up -d app`

- Check container logs: `docker compose logs -f app`

- To stop the container: `docker compose down`.

- For further process, you can directly log into the bash with `docker compose exec app bash` and run specific application configuration related command.

## Health check APIs

- endpoint: `http://localhost:8080/api/v1/health`
