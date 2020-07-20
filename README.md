#pghandler
Docker file with Golang and Postgresql 9 deploying a http server to serv commands for backup / restore / and bench. 
###build image:
docker build -t pghandler .

###run image:
docker run -d /<br>
-e PG_HOST=localhost -e PG_PORT=5432 -e PG_SCHEMA=public /<br>
-e PG_DB=tenant01 -e PG_PASSWORD=password /<br>
-p 8080:8080 /<br>
-v ~/workspace/pghandler/logs:/app/logs /<br>
-v ~/workspace/pghandler/pgbackup:/pgbackup /<br>
pghandler


