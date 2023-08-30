## run
```
# go run .
panic: 需要 sm ms token才能运行


# go run . -user <username> -password <password>


# go run . -token <sm ms token>


# SM_MS_TOKEN=<sm ms token> go run .
```

## build

```
# go build -ldflags="-s -w"
# ./go-sm-ms


# GOOS=windows go build -ldflags="-s -w"
```