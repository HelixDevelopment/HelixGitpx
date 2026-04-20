plugins {
    id("helix.convention")
    kotlin("plugin.serialization") version "2.1.0"
}

kotlin {
    sourceSets {
        val commonMain by getting {
            kotlin.exclude("gen/**")
            dependencies {
                implementation("io.ktor:ktor-client-core:3.0.2")
                implementation("io.ktor:ktor-client-content-negotiation:3.0.2")
                implementation("io.ktor:ktor-serialization-kotlinx-json:3.0.2")
                implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.9.0")
                implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:1.7.3")
            }
        }
        val commonTest by getting {
            dependencies { implementation(kotlin("test")) }
        }
    }
}
