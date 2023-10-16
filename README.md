docker run --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=mysql -e MYSQL_DATABASE=Backoffice -e MYSQL_USER=mysql -e MYSQL_PASSWORD=mysql -d mysql:latest


 docker run -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=mysql -e MYSQL_DATABASE=Backoffice -e MYSQL_USER=mysql -e MYSQL_PASSWORD=mysql backoffice-db

 docker exec -it 9ba09d1362a9 bash