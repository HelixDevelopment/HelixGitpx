plugins {
    id("org.jetbrains.kotlin.multiplatform")
    id("io.gitlab.arturbosch.detekt")
    id("org.jlleitschuh.gradle.ktlint")
}

kotlin {
    jvmToolchain(21)
    jvm()
    iosX64()
    iosArm64()
    iosSimulatorArm64()
    linuxX64()
    androidTarget()
}

detekt {
    buildUponDefaultConfig = true
    allRules = false
}

tasks.withType<org.jlleitschuh.gradle.ktlint.tasks.KtLintCheckTask>().configureEach {
    exclude { it.file.path.contains("/gen/") }
}
