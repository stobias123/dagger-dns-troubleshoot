plugins {
    kotlin("jvm") version "1.9.0"
    application
}

group = "org.example"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

dependencies {
    testImplementation(kotlin("test"))
    implementation("redis.clients:jedis:5.1.2")
    implementation("org.slf4j:slf4j-api:1.7.5")
    implementation("org.slf4j:slf4j-log4j12:1.7.5")
}
tasks.test {
    useJUnitPlatform()
    testLogging {
        outputs.upToDateWhen {false}
        showStandardStreams = true
        showExceptions = true
        showCauses = true
        showStackTraces = true
    }
}

kotlin {
    jvmToolchain(8)
}

application {
    mainClass.set("MainKt")
}