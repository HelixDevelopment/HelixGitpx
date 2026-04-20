plugins { `kotlin-dsl` }
repositories { gradlePluginPortal(); google(); mavenCentral() }
dependencies {
    implementation("org.jetbrains.kotlin:kotlin-gradle-plugin:2.1.0")
    implementation("com.android.tools.build:gradle:8.7.2")
    implementation("io.gitlab.arturbosch.detekt:detekt-gradle-plugin:1.23.7")
    implementation("org.jlleitschuh.gradle:ktlint-gradle:12.1.2")
}
