docker build . -t blog:latest 
docker run -d -p 4000:4000 -t blog:latest
