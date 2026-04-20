plugins {
    id("org.jetbrains.kotlin.multiplatform")
    id("com.android.library")
    id("io.gitlab.arturbosch.detekt")
    id("org.jlleitschuh.gradle.ktlint")
}

extensions.configure<com.android.build.api.dsl.LibraryExtension> {
    namespace = "dev.helixgitpx.${project.name}"
    compileSdk = 35
    defaultConfig { minSdk = 26 }
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
