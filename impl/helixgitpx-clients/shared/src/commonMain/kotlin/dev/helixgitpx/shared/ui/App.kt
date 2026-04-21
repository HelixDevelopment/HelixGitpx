package dev.helixgitpx.shared.ui

import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun App() {
    val primary = Color(0xFF2563EB)
    var selected by remember { mutableStateOf(Tab.Repos) }

    MaterialTheme(
        colorScheme = lightColorScheme(primary = primary),
    ) {
        Scaffold(
            topBar = {
                TopAppBar(
                    title = { Text("HelixGitpx", fontWeight = FontWeight.SemiBold) },
                    colors = TopAppBarDefaults.topAppBarColors(
                        containerColor = primary,
                        titleContentColor = Color.White,
                    ),
                )
            },
            bottomBar = {
                NavigationBar {
                    Tab.entries.forEach { tab ->
                        NavigationBarItem(
                            selected = selected == tab,
                            onClick = { selected = tab },
                            icon = { Text(tab.glyph) },
                            label = { Text(tab.label) },
                        )
                    }
                }
            },
        ) { innerPad ->
            Box(
                modifier = Modifier.fillMaxSize().padding(innerPad),
                contentAlignment = Alignment.Center,
            ) {
                when (selected) {
                    Tab.Repos -> RepoList()
                    Tab.Prs -> Placeholder("Pull requests")
                    Tab.Issues -> Placeholder("Issues")
                    Tab.Conflicts -> Placeholder("Conflicts inbox")
                    Tab.Settings -> Placeholder("Settings")
                }
            }
        }
    }
}

enum class Tab(val glyph: String, val label: String) {
    Repos("■", "Repos"),
    Prs("↻", "PRs"),
    Issues("!", "Issues"),
    Conflicts("⚠", "Conflicts"),
    Settings("⚙", "Settings"),
}
