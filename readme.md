// create the volume
docker volume create mysql-data
// check if the volume is created
docker volume ls
// inspect detail of the volume
docker volume inspect mysql-data
// create the container
<!-- docker run --name gps-tracking-mysql -v mysql-data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=davidzhu2020 -p 3306:3306 -d mysql:latest -->

docker run --name gps-tracking-mysql -v /Volumes/enhui/gps-tracking-backend/mysql-data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=davidzhu2020 -e MYSQL_ROOT_HOST=% -p 3306:3306 -d mysql:latest

docker run --name myadmin -d --link gps-tracking-mysql:db -p 5000:80 phpmyadmin/phpmyadmin

docker run --rm -it --name go-gps-tracking --link gps-tracking-mysql:db -v $PWD:/go/src/gps-tracking-go-backend -p 8080:8080 golang