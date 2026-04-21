package dev.helixgitpx.shared.ui

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp

data class Repo(val name: String, val upstreams: List<String>)

@Composable
fun RepoList() {
    // Placeholder data until the shared Connect client is wired to
    // RepoService. See shared/src/commonMain/kotlin/dev/helixgitpx/network/.
    val seed = remember {
        listOf(
            Repo("acme/hello", listOf("github", "gitlab")),
            Repo("acme/backend", listOf("github", "gitlab", "gitflic")),
            Repo("personal/dotfiles", listOf("gitverse", "github")),
        )
    }
    LazyColumn(
        modifier = Modifier.fillMaxSize().padding(horizontal = 16.dp),
    ) {
        items(seed) { repo ->
            ListItem(
                headlineContent = { Text(repo.name, fontWeight = FontWeight.SemiBold) },
                supportingContent = { Text(repo.upstreams.joinToString(" · ")) },
            )
            HorizontalDivider()
        }
    }
}

@Composable
fun Placeholder(title: String) {
    Column(
        modifier = Modifier.fillMaxSize().padding(16.dp),
        verticalArrangement = Arrangement.Center,
    ) {
        Text(title, style = MaterialTheme.typography.headlineSmall)
        Spacer(Modifier.height(8.dp))
        Text(
            "This surface is part of the KMP/Compose GA scaffold. " +
                "Wire it to the shared Connect client when the screen ships.",
            style = MaterialTheme.typography.bodyMedium,
        )
    }
}
