# basic paxos

## code

```
- cmd: main.go
- pkg
  - core: testable basic paxos
  - grpc: grpc transport
  - http: http transport
  - util: utility
```

## example

- kind=http, num=3, prepareNum=2, voteNum=2

```
go run cmd/main.go -kind=http -self=:5001 -others=:5002,:5003 -path=storage/1 -prepareNum=2 -voteNum=2 -first=1 -timeout=3
go run cmd/main.go -kind=http -self=:5002 -others=:5003,:5001 -path=storage/2 -prepareNum=2 -voteNum=2 -first=2 -timeout=3
go run cmd/main.go -kind=http -self=:5003 -others=:5001,:5002 -path=storage/3 -prepareNum=2 -voteNum=2 -first=3 -timeout=3
```

- kind=http, num=5, prepareNum=3, voteNum=3

```
go run cmd/main.go -kind=http -self=:5001 -others=:5002,:5003,:5004,:5005 -path=storage/1 -prepareNum=3 -voteNum=3 -first=1 -timeout=3
go run cmd/main.go -kind=http -self=:5002 -others=:5003,:5004,:5005,:5001 -path=storage/2 -prepareNum=3 -voteNum=3 -first=2 -timeout=3
go run cmd/main.go -kind=http -self=:5003 -others=:5004,:5005,:5001,:5002 -path=storage/3 -prepareNum=3 -voteNum=3 -first=3 -timeout=3
go run cmd/main.go -kind=http -self=:5004 -others=:5005,:5001,:5002,:5003 -path=storage/4 -prepareNum=3 -voteNum=3 -first=4 -timeout=3
go run cmd/main.go -kind=http -self=:5005 -others=:5001,:5002,:5003,:5004 -path=storage/5 -prepareNum=3 -voteNum=3 -first=5 -timeout=3
```

- kind=http, num=5, prepareNum=2, voteNum=4

```
go run cmd/main.go -kind=http -self=:5001 -others=:5002,:5003,:5004,:5005 -path=storage/1 -prepareNum=2 -voteNum=4 -first=1 -timeout=3
go run cmd/main.go -kind=http -self=:5002 -others=:5003,:5004,:5005,:5001 -path=storage/2 -prepareNum=2 -voteNum=4 -first=2 -timeout=3
go run cmd/main.go -kind=http -self=:5003 -others=:5004,:5005,:5001,:5002 -path=storage/3 -prepareNum=2 -voteNum=4 -first=3 -timeout=3
go run cmd/main.go -kind=http -self=:5004 -others=:5005,:5001,:5002,:5003 -path=storage/4 -prepareNum=2 -voteNum=4 -first=4 -timeout=3
go run cmd/main.go -kind=http -self=:5005 -others=:5001,:5002,:5003,:5004 -path=storage/5 -prepareNum=2 -voteNum=4 -first=5 -timeout=3
```

- kind=http, num=5, prepareNum=4, voteNum=2

```
go run cmd/main.go -kind=http -self=:5001 -others=:5002,:5003,:5004,:5005 -path=storage/1 -prepareNum=4 -voteNum=2 -first=1 -timeout=3
go run cmd/main.go -kind=http -self=:5002 -others=:5003,:5004,:5005,:5001 -path=storage/2 -prepareNum=4 -voteNum=2 -first=2 -timeout=3
go run cmd/main.go -kind=http -self=:5003 -others=:5004,:5005,:5001,:5002 -path=storage/3 -prepareNum=4 -voteNum=2 -first=3 -timeout=3
go run cmd/main.go -kind=http -self=:5004 -others=:5005,:5001,:5002,:5003 -path=storage/4 -prepareNum=4 -voteNum=2 -first=4 -timeout=3
go run cmd/main.go -kind=http -self=:5005 -others=:5001,:5002,:5003,:5004 -path=storage/5 -prepareNum=4 -voteNum=2 -first=5 -timeout=3
```

- kind=grpc, num=3, prepareNum=2, voteNum=2

```
go run cmd/main.go -kind=grpc -self=:5001 -others=:5002,:5003 -path=storage/1 -prepareNum=2 -voteNum=2 -first=1 -timeout=3
go run cmd/main.go -kind=grpc -self=:5002 -others=:5003,:5001 -path=storage/2 -prepareNum=2 -voteNum=2 -first=2 -timeout=3
go run cmd/main.go -kind=grpc -self=:5003 -others=:5001,:5002 -path=storage/3 -prepareNum=2 -voteNum=2 -first=3 -timeout=3
```

- kind=grpc, num=5, prepareNum=3, voteNum=3

```
go run cmd/main.go -kind=grpc -self=:5001 -others=:5002,:5003,:5004,:5005 -path=storage/1 -prepareNum=3 -voteNum=3 -first=1 -timeout=3
go run cmd/main.go -kind=grpc -self=:5002 -others=:5003,:5004,:5005,:5001 -path=storage/2 -prepareNum=3 -voteNum=3 -first=2 -timeout=3
go run cmd/main.go -kind=grpc -self=:5003 -others=:5004,:5005,:5001,:5002 -path=storage/3 -prepareNum=3 -voteNum=3 -first=3 -timeout=3
go run cmd/main.go -kind=grpc -self=:5004 -others=:5005,:5001,:5002,:5003 -path=storage/4 -prepareNum=3 -voteNum=3 -first=4 -timeout=3
go run cmd/main.go -kind=grpc -self=:5005 -others=:5001,:5002,:5003,:5004 -path=storage/5 -prepareNum=3 -voteNum=3 -first=5 -timeout=3
```

- kind=grpc, num=5, prepareNum=2, voteNum=4

```
go run cmd/main.go -kind=grpc -self=:5001 -others=:5002,:5003,:5004,:5005 -path=storage/1 -prepareNum=2 -voteNum=4 -first=1 -timeout=3
go run cmd/main.go -kind=grpc -self=:5002 -others=:5003,:5004,:5005,:5001 -path=storage/2 -prepareNum=2 -voteNum=4 -first=2 -timeout=3
go run cmd/main.go -kind=grpc -self=:5003 -others=:5004,:5005,:5001,:5002 -path=storage/3 -prepareNum=2 -voteNum=4 -first=3 -timeout=3
go run cmd/main.go -kind=grpc -self=:5004 -others=:5005,:5001,:5002,:5003 -path=storage/4 -prepareNum=2 -voteNum=4 -first=4 -timeout=3
go run cmd/main.go -kind=grpc -self=:5005 -others=:5001,:5002,:5003,:5004 -path=storage/5 -prepareNum=2 -voteNum=4 -first=5 -timeout=3
```

- kind=grpc, num=5, prepareNum=4, voteNum=2

```
go run cmd/main.go -kind=grpc -self=:5001 -others=:5002,:5003,:5004,:5005 -path=storage/1 -prepareNum=4 -voteNum=2 -first=1 -timeout=3
go run cmd/main.go -kind=grpc -self=:5002 -others=:5003,:5004,:5005,:5001 -path=storage/2 -prepareNum=4 -voteNum=2 -first=2 -timeout=3
go run cmd/main.go -kind=grpc -self=:5003 -others=:5004,:5005,:5001,:5002 -path=storage/3 -prepareNum=4 -voteNum=2 -first=3 -timeout=3
go run cmd/main.go -kind=grpc -self=:5004 -others=:5005,:5001,:5002,:5003 -path=storage/4 -prepareNum=4 -voteNum=2 -first=4 -timeout=3
go run cmd/main.go -kind=grpc -self=:5005 -others=:5001,:5002,:5003,:5004 -path=storage/5 -prepareNum=4 -voteNum=2 -first=5 -timeout=3
```

## reference

- [paxos](https://github.com/gaoxinge/distributed-system/blob/master/summary/paxos/README.md)