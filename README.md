# Go/Gin Server

- API routes secured by JWT authentication
- User data stored in Postgres DB
- File storage is a local Docker managed volumn
- OneTime file access managed with JWT with default expiration set to 15m
    + download checks memory fileStore and valid expiration before serving file
    + after download fileStore record is removed to invalidate the download token

# Node/Express Server

- in progress

# nginx

- http://localhost/... Go server open endpoints
- http://localhost/api/... Go server secure endpoints
- http://localhost/node/... Node server open endpoints

# installation

docker pull peterjbishop/cautious-dollop:latest 
docker-compose build --no-cache 
docker-compose up

## Store files in a Docker Container.

POST   /upload     
GET    /files              
GET    /download/:filename       
DELETE /delete/:filename  

## Download a file to your local device with CURL

example:
 cd {your_download_directory}
 curl -o "The Home Depot - Cart.png" "http://localhost:8080/download/The%20Home%20Depot%20-%20Cart.png"

## build

docker build -t peterjbishop/cautious-dollop:latest . 
docker push peterjbishop/cautious-dollop:latest 
docker pull peterjbishop/cautious-dollop:latest 
docker-compose down 
docker-compose build --no-cache 
docker-compose up

