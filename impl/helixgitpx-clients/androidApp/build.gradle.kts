plugins {
    id("com.android.application")
    kotlin("android")
    id("org.jetbrains.compose")
}

android {
    namespace = "dev.helixgitpx.android"
    compileSdk = 35

    defaultConfig {
        applicationId = "dev.helixgitpx.android"
        minSdk = 26
        targetSdk = 35
        versionCode = 1
        versionName = "0.1.0"
    }
}

dependencies {
    implementation(project(":shared"))
    implementation("androidx.activity:activity-compose:1.10.1")
    implementation(compose.runtime)
    implementation(compose.ui)
    implementation(compose.material3)
}
