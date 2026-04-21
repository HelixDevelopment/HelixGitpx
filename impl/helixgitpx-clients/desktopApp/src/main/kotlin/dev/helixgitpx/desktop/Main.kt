package dev.helixgitpx.desktop

import androidx.compose.ui.unit.dp
import androidx.compose.ui.window.Window
import androidx.compose.ui.window.application
import androidx.compose.ui.window.rememberWindowState
import dev.helixgitpx.desktop.ui.App

fun main() = application {
    val state = rememberWindowState(width = 1200.dp, height = 800.dp)
    Window(
        onCloseRequest = ::exitApplication,
        state = state,
        title = "HelixGitpx",
    ) {
        App()
    }
}
