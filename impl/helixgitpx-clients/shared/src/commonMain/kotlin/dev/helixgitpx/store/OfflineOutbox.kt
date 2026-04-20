package dev.helixgitpx.store

/** OfflineOutbox buffers user actions performed offline for replay when
 *  connectivity returns. Backed by SQLDelight in commonMain; actual schema
 *  lives in shared/src/commonMain/sqldelight/ (added per-platform in M6 hardening). */
interface OfflineOutbox {
    suspend fun enqueue(action: OutboxAction)
    suspend fun drainPending(): List<OutboxAction>
    suspend fun markSent(id: String)
}

data class OutboxAction(
    val id: String,
    val kind: String,
    val payload: ByteArray,
    val createdAt: Long,
)
