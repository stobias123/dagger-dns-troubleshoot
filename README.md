# Dagger Jedis Failure Example Repo

You should see this error.

```shell

```

To reproduce

```shell
go mod tidy
go run ./ci
```

You should see redis cli ping successfully

```shell
71: exec redis-cli -h redis-cluster -p 30000 ping DONE
71: [0.85s] PONG
71: exec redis-cli -h redis-cluster -p 30000 ping DONE
```

But jedis fail
```shell
70: [57.3s] ActiveRideCacheTest > test running redis-cli command() FAILED
70: [57.3s]     java.lang.AssertionError at jedistest.kt:32
70: [61.7s] 
70: [61.7s] ActiveRideCacheTest > test a simple redis ping() FAILED
70: [61.7s]     redis.clients.jedis.exceptions.JedisClusterOperationException at jedistest.kt:23
```
