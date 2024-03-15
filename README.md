# Dagger Jedis Failure Example Repo

## The Issue

Jedis seems to get an error.
```
        at worker.org.gradle.process.internal.worker.GradleWorkerMain.main(GradleWorkerMain.java:74)
        Suppressed: java.net.ConnectException: Connection refused (Connection refused)
                at java.net.PlainSocketImpl.socketConnect(Native Method)
                at java.net.AbstractPlainSocketImpl.doConnect(AbstractPlainSocketImpl.java:350)
                at java.net.AbstractPlainSocketImpl.connectToAddress(AbstractPlainSocketImpl.java:206)
                at java.net.AbstractPlainSocketImpl.connect(AbstractPlainSocketImpl.java:188)
                at java.net.SocksSocketImpl.connect(SocksSocketImpl.java:392)
                at java.net.Socket.connect(Socket.java:607)
                at redis.clients.jedis.DefaultJedisSocketFactory.connectToFirstSuccessfulHost(DefaultJedisSocketFactory.java:75)
```

## Reproduction

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
