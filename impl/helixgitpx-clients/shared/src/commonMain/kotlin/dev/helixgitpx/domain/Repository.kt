package dev.helixgitpx.domain

import kotlinx.serialization.Serializable

@Serializable
data class Repository(
    val id: String,
    val orgId: String,
    val slug: String,
    val defaultBranch: String = "main",
)

@Serializable
data class PullRequest(
    val id: String,
    val repoId: String,
    val title: String,
    val sourceBranch: String,
    val targetBranch: String,
)
