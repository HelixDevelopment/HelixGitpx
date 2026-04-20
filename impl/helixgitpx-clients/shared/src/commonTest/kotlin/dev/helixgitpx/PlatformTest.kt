package dev.helixgitpx

import kotlin.test.Test
import kotlin.test.assertEquals

class GreetingTest {
    @Test
    fun greetingWorld() {
        assertEquals("hello, world", HelixGitpx.greeting("world"))
    }
}
