import org.junit.jupiter.api.Test
import redis.clients.jedis.DefaultJedisClientConfig
import redis.clients.jedis.HostAndPort
import redis.clients.jedis.JedisCluster

class ActiveRideCacheTest {

  companion object {
    private val redis by lazy {
      val hosts =
        "redis-cluster:30000,redis-cluster:30001,redis-cluster:30002,redis-cluster:30003,redis-cluster:30004,redis-cluster:30005"
      JedisCluster(
        hosts.toHostPorts(),
        DefaultJedisClientConfig
          .builder()
          .build()
      )
    }
  }

  @Test
  fun `test a simple redis ping`() {
    redis.set("test", "test")
    val result = redis.get("test")
    assert(result == "test")
  }

  @Test
  fun `test running redis-cli command`() {
    val result = Runtime.getRuntime().exec("redis-cli -h redis-cluster -p 30000 ping").inputStream.bufferedReader().readText()
    assert(result == "PONG")
  }
}



fun String.parseHostPorts(defaultPort: Int = 6379): List<Pair<String, Int>> {
  return split(",")
    .asSequence()
    .map { it.split(":") }
    .map {
      val host = it[0]
      if (host.isBlank()) {
        throw IllegalArgumentException("Host cannot be empty: $this.")
      }

      val port = try {
        it.getOrNull(1)?.toInt() ?: defaultPort
      } catch (e: Exception) {
        throw IllegalArgumentException("Port is not a valid number: $this.")
      }
      Pair(host, port)
    }
    .toList()
}
fun String.toHostPorts(defaultPort: Int = 6379): Set<HostAndPort> {
  return parseHostPorts(defaultPort)
    .map { HostAndPort(it.first, it.second) }
    .toSet()
}