package com.helixgitpx.gitingress.v1

import com.helixgitpx.gitingress.v1.GitIngressServiceGrpc.getServiceDescriptor
import io.grpc.CallOptions
import io.grpc.CallOptions.DEFAULT
import io.grpc.Channel
import io.grpc.Metadata
import io.grpc.MethodDescriptor
import io.grpc.ServerServiceDefinition
import io.grpc.ServerServiceDefinition.builder
import io.grpc.ServiceDescriptor
import io.grpc.Status.UNIMPLEMENTED
import io.grpc.StatusException
import io.grpc.kotlin.AbstractCoroutineServerImpl
import io.grpc.kotlin.AbstractCoroutineStub
import io.grpc.kotlin.ClientCalls.unaryRpc
import io.grpc.kotlin.ServerCalls.unaryServerMethodDefinition
import io.grpc.kotlin.StubFor
import kotlin.String
import kotlin.coroutines.CoroutineContext
import kotlin.coroutines.EmptyCoroutineContext
import kotlin.jvm.JvmOverloads
import kotlin.jvm.JvmStatic

/**
 * Holder for Kotlin coroutine-based client and server APIs for helixgitpx.gitingress.v1.GitIngressService.
 */
public object GitIngressServiceGrpcKt {
  public const val SERVICE_NAME: String = GitIngressServiceGrpc.SERVICE_NAME

  @JvmStatic
  public val serviceDescriptor: ServiceDescriptor
    get() = getServiceDescriptor()

  public val listRecentPushesMethod:
      MethodDescriptor<ListRecentPushesRequest, ListRecentPushesResponse>
    @JvmStatic
    get() = GitIngressServiceGrpc.getListRecentPushesMethod()

  public val getRepoStatsMethod: MethodDescriptor<GetRepoStatsRequest, RepoStats>
    @JvmStatic
    get() = GitIngressServiceGrpc.getGetRepoStatsMethod()

  public val getQuotaStatusMethod: MethodDescriptor<QuotaStatusRequest, QuotaStatus>
    @JvmStatic
    get() = GitIngressServiceGrpc.getGetQuotaStatusMethod()

  /**
   * A stub for issuing RPCs to a(n) helixgitpx.gitingress.v1.GitIngressService service as suspending coroutines.
   */
  @StubFor(GitIngressServiceGrpc::class)
  public class GitIngressServiceCoroutineStub @JvmOverloads constructor(
    channel: Channel,
    callOptions: CallOptions = DEFAULT,
  ) : AbstractCoroutineStub<GitIngressServiceCoroutineStub>(channel, callOptions) {
    override fun build(channel: Channel, callOptions: CallOptions): GitIngressServiceCoroutineStub = GitIngressServiceCoroutineStub(channel, callOptions)

    /**
     * Executes this RPC and returns the response message, suspending until the RPC completes
     * with [`Status.OK`][io.grpc.Status].  If the RPC completes with another status, a corresponding
     * [StatusException] is thrown.  If this coroutine is cancelled, the RPC is also cancelled
     * with the corresponding exception as a cause.
     *
     * @param request The request message to send to the server.
     *
     * @param headers Metadata to attach to the request.  Most users will not need this.
     *
     * @return The single response from the server.
     */
    public suspend fun listRecentPushes(request: ListRecentPushesRequest, headers: Metadata = Metadata()): ListRecentPushesResponse = unaryRpc(
      channel,
      GitIngressServiceGrpc.getListRecentPushesMethod(),
      request,
      callOptions,
      headers
    )

    /**
     * Executes this RPC and returns the response message, suspending until the RPC completes
     * with [`Status.OK`][io.grpc.Status].  If the RPC completes with another status, a corresponding
     * [StatusException] is thrown.  If this coroutine is cancelled, the RPC is also cancelled
     * with the corresponding exception as a cause.
     *
     * @param request The request message to send to the server.
     *
     * @param headers Metadata to attach to the request.  Most users will not need this.
     *
     * @return The single response from the server.
     */
    public suspend fun getRepoStats(request: GetRepoStatsRequest, headers: Metadata = Metadata()): RepoStats = unaryRpc(
      channel,
      GitIngressServiceGrpc.getGetRepoStatsMethod(),
      request,
      callOptions,
      headers
    )

    /**
     * Executes this RPC and returns the response message, suspending until the RPC completes
     * with [`Status.OK`][io.grpc.Status].  If the RPC completes with another status, a corresponding
     * [StatusException] is thrown.  If this coroutine is cancelled, the RPC is also cancelled
     * with the corresponding exception as a cause.
     *
     * @param request The request message to send to the server.
     *
     * @param headers Metadata to attach to the request.  Most users will not need this.
     *
     * @return The single response from the server.
     */
    public suspend fun getQuotaStatus(request: QuotaStatusRequest, headers: Metadata = Metadata()): QuotaStatus = unaryRpc(
      channel,
      GitIngressServiceGrpc.getGetQuotaStatusMethod(),
      request,
      callOptions,
      headers
    )
  }

  /**
   * Skeletal implementation of the helixgitpx.gitingress.v1.GitIngressService service based on Kotlin coroutines.
   */
  public abstract class GitIngressServiceCoroutineImplBase(
    coroutineContext: CoroutineContext = EmptyCoroutineContext,
  ) : AbstractCoroutineServerImpl(coroutineContext) {
    /**
     * Returns the response to an RPC for helixgitpx.gitingress.v1.GitIngressService.ListRecentPushes.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun listRecentPushes(request: ListRecentPushesRequest): ListRecentPushesResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.gitingress.v1.GitIngressService.ListRecentPushes is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.gitingress.v1.GitIngressService.GetRepoStats.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun getRepoStats(request: GetRepoStatsRequest): RepoStats = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.gitingress.v1.GitIngressService.GetRepoStats is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.gitingress.v1.GitIngressService.GetQuotaStatus.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun getQuotaStatus(request: QuotaStatusRequest): QuotaStatus = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.gitingress.v1.GitIngressService.GetQuotaStatus is unimplemented"))

    final override fun bindService(): ServerServiceDefinition = builder(getServiceDescriptor())
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = GitIngressServiceGrpc.getListRecentPushesMethod(),
      implementation = ::listRecentPushes
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = GitIngressServiceGrpc.getGetRepoStatsMethod(),
      implementation = ::getRepoStats
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = GitIngressServiceGrpc.getGetQuotaStatusMethod(),
      implementation = ::getQuotaStatus
    )).build()
  }
}
