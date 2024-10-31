# database
Launch postgres in docker using
```
docker run --name gator-db -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=gator -p 5432:5432 postgres:16.4
```