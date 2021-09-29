# go-logrotate
go logrotate

## install

```
go get -u github.com/lyuangg/go-logrotate
```


## usage


```go
import github.com/lyuangg/go-logrotate/log

# stdout
log.Println("hello")

# prefix
log.SetPrefix("test")

# new
mylogger := log.New("./test.log", 2, "[test-log]")

mylogger.Printf("hello")
```


