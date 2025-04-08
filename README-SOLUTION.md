# CSV Utility
A tiny project on csv utility

# Dependencies to run

- go 1.23.5

<p>To run the application</p>
go to csvserver folder and use `go run main.go`

# Endpoints

- Echo

```
curl --location 'localhost:8080/api/v1/echo' \
--form 'file=@"/Users/dinesh/Downloads/matrix.csv"'
```

- Sum

```
curl --location 'localhost:8080/api/v1/sum' \
--header 'Content-Type: multipart/form-data' \
--form 'file=@"/Users/dinesh/Downloads/matrix.csv"'
```

- Flatten
```
curl --location 'localhost:8080/api/v1/flatten' \
--header 'Content-Type: multipart/form-data' \
--form 'file=@"/Users/dinesh/Downloads/matrix.csv"'
```

- Invert
```
curl --location 'localhost:8080/api/v1/invert' \
--header 'Content-Type: multipart/form-data' \
--form 'file=@"/Users/dinesh/Downloads/matrix.csv"'
```

- Multiply

```
curl --location 'localhost:8080/api/v1/multiply' \
--header 'Content-Type: multipart/form-data' \
--form 'file=@"/Users/dinesh/Downloads/matrix.csv"'
```