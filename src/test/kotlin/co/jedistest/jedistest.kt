import java.io.IOException
import java.net.InetAddress
import java.net.UnknownHostException
import org.junit.jupiter.api.Test
import redis.clients.jedis.DefaultJedisClientConfig
import redis.clients.jedis.HostAndPort
import redis.clients.jedis.JedisCluster

class JedisTest {

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
  fun `test redis-cluster IP resolution`() {
    try {
      val hostName = "redis-cluster"
      val inetAddresses = InetAddress.getAllByName(hostName)
      println("IP addresses for $hostName:")
      inetAddresses.forEach { println(it.hostAddress) }

      val ports = listOf(30000, 30001, 30002, 30003, 30004, 30005)
      ports.forEach { port ->
        try {
          val result = Runtime.getRuntime().exec("redis-cli -h $hostName -p $port ping").inputStream.bufferedReader().readText()
          println("Ping result for port $port: $result")
          assert(result.contains("PONG"))
        } catch (e: IOException) {
          println("Failed to ping redis-cluster on port $port")
          e.printStackTrace()
        }
      }
    } catch (e: UnknownHostException) {
      println("Failed to resolve host: redis-cluster")
      e.printStackTrace()
    }
  }

  @Test
  fun `test a simple redis ping`() {
    val hosts = "redis-cluster:30000,redis-cluster:30001,redis-cluster:30002,redis-cluster:30003,redis-cluster:30004,redis-cluster:30005"
    val redisInTest = JedisCluster(
        hosts.toHostPorts(),
        DefaultJedisClientConfig
          .builder()
          .build()
      )
   
    try {
        redisInTest.set("test", "test")
        val result = redis.get("test")
    } catch (e: Exception) {
        println(e)
        throw(e)
    }
  }

  @Test
  fun `test running redis-cli command`() {
    val result = Runtime.getRuntime().exec("redis-cli -h redis-cluster -p 30000 ping").inputStream.bufferedReader().readText()
    println(result)
    assert(result.contains("PONG"))
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