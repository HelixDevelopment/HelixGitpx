plugins {
    kotlin("multiplatform") version "2.1.0" apply false
    id("io.gitlab.arturbosch.detekt") version "1.23.7" apply false
    id("org.jlleitschuh.gradle.ktlint") version "12.1.2" apply false
}
allprojects {
    group = "dev.helixgitpx"
    version = "0.1.0"
}
