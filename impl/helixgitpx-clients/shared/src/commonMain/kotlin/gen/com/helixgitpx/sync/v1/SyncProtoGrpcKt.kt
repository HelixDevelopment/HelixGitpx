package com.helixgitpx.sync.v1

import com.google.protobuf.Empty
import com.helixgitpx.sync.v1.SyncServiceGrpc.getServiceDescriptor
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
 * Holder for Kotlin coroutine-based client and server APIs for helixgitpx.sync.v1.SyncService.
 */
public object SyncServiceGrpcKt {
  public const val SERVICE_NAME: String = SyncServiceGrpc.SERVICE_NAME

  @JvmStatic
  public val serviceDescriptor: ServiceDescriptor
    get() = getServiceDescriptor()

  public val listJobsMethod: MethodDescriptor<ListJobsRequest, ListJobsResponse>
    @JvmStatic
    get() = SyncServiceGrpc.getListJobsMethod()

  public val getJobMethod: MethodDescriptor<GetJobRequest, SyncJob>
    @JvmStatic
    get() = SyncServiceGrpc.getGetJobMethod()

  public val retryJobMethod: MethodDescriptor<RetryJobRequest, SyncJob>
    @JvmStatic
    get() = SyncServiceGrpc.getRetryJobMethod()

  public val cancelJobMethod: MethodDescriptor<CancelJobRequest, Empty>
    @JvmStatic
    get() = SyncServiceGrpc.getCancelJobMethod()

  public val listDLQMethod: MethodDescriptor<ListDLQRequest, ListDLQResponse>
    @JvmStatic
    get() = SyncServiceGrpc.getListDLQMethod()

  public val requeueFromDLQMethod: MethodDescriptor<RequeueFromDLQRequest, SyncJob>
    @JvmStatic
    get() = SyncServiceGrpc.getRequeueFromDLQMethod()

  /**
   * A stub for issuing RPCs to a(n) helixgitpx.sync.v1.SyncService service as suspending coroutines.
   */
  @StubFor(SyncServiceGrpc::class)
  public class SyncServiceCoroutineStub @JvmOverloads constructor(
    channel: Channel,
    callOptions: CallOptions = DEFAULT,
  ) : AbstractCoroutineStub<SyncServiceCoroutineStub>(channel, callOptions) {
    override fun build(channel: Channel, callOptions: CallOptions): SyncServiceCoroutineStub = SyncServiceCoroutineStub(channel, callOptions)

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
    public suspend fun listJobs(request: ListJobsRequest, headers: Metadata = Metadata()): ListJobsResponse = unaryRpc(
      channel,
      SyncServiceGrpc.getListJobsMethod(),
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
    public suspend fun getJob(request: GetJobRequest, headers: Metadata = Metadata()): SyncJob = unaryRpc(
      channel,
      SyncServiceGrpc.getGetJobMethod(),
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
    public suspend fun retryJob(request: RetryJobRequest, headers: Metadata = Metadata()): SyncJob = unaryRpc(
      channel,
      SyncServiceGrpc.getRetryJobMethod(),
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
    public suspend fun cancelJob(request: CancelJobRequest, headers: Metadata = Metadata()): Empty = unaryRpc(
      channel,
      SyncServiceGrpc.getCancelJobMethod(),
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
    public suspend fun listDLQ(request: ListDLQRequest, headers: Metadata = Metadata()): ListDLQResponse = unaryRpc(
      channel,
      SyncServiceGrpc.getListDLQMethod(),
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
    public suspend fun requeueFromDLQ(request: RequeueFromDLQRequest, headers: Metadata = Metadata()): SyncJob = unaryRpc(
      channel,
      SyncServiceGrpc.getRequeueFromDLQMethod(),
      request,
      callOptions,
      headers
    )
  }

  /**
   * Skeletal implementation of the helixgitpx.sync.v1.SyncService service based on Kotlin coroutines.
   */
  public abstract class SyncServiceCoroutineImplBase(
    coroutineContext: CoroutineContext = EmptyCoroutineContext,
  ) : AbstractCoroutineServerImpl(coroutineContext) {
    /**
     * Returns the response to an RPC for helixgitpx.sync.v1.SyncService.ListJobs.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun listJobs(request: ListJobsRequest): ListJobsResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.sync.v1.SyncService.ListJobs is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.sync.v1.SyncService.GetJob.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun getJob(request: GetJobRequest): SyncJob = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.sync.v1.SyncService.GetJob is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.sync.v1.SyncService.RetryJob.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun retryJob(request: RetryJobRequest): SyncJob = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.sync.v1.SyncService.RetryJob is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.sync.v1.SyncService.CancelJob.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun cancelJob(request: CancelJobRequest): Empty = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.sync.v1.SyncService.CancelJob is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.sync.v1.SyncService.ListDLQ.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun listDLQ(request: ListDLQRequest): ListDLQResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.sync.v1.SyncService.ListDLQ is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.sync.v1.SyncService.RequeueFromDLQ.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun requeueFromDLQ(request: RequeueFromDLQRequest): SyncJob = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.sync.v1.SyncService.RequeueFromDLQ is unimplemented"))

    final override fun bindService(): ServerServiceDefinition = builder(getServiceDescriptor())
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = SyncServiceGrpc.getListJobsMethod(),
      implementation = ::listJobs
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = SyncServiceGrpc.getGetJobMethod(),
      implementation = ::getJob
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = SyncServiceGrpc.getRetryJobMethod(),
      implementation = ::retryJob
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = SyncServiceGrpc.getCancelJobMethod(),
      implementation = ::cancelJob
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = SyncServiceGrpc.getListDLQMethod(),
      implementation = ::listDLQ
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = SyncServiceGrpc.getRequeueFromDLQMethod(),
      implementation = ::requeueFromDLQ
    )).build()
  }
}
