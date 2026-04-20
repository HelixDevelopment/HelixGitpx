package dev.helixgitpx.network

import io.ktor.client.HttpClient
import io.ktor.client.plugins.contentnegotiation.ContentNegotiation
import io.ktor.serialization.kotlinx.json.json

/** ApiClient wraps Ktor's HttpClient with JSON content negotiation for
 *  Connect-compatible HelixGitpx service endpoints. One instance per service. */
class ApiClient(val baseUrl: String) {
    val http: HttpClient = HttpClient {
        install(ContentNegotiation) { json() }
    }
}
