package com.helixgitpx.collab.v1

import com.google.protobuf.Empty
import com.helixgitpx.collab.v1.CollabServiceGrpc.getServiceDescriptor
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
import io.grpc.kotlin.ClientCalls.bidiStreamingRpc
import io.grpc.kotlin.ClientCalls.unaryRpc
import io.grpc.kotlin.ServerCalls.bidiStreamingServerMethodDefinition
import io.grpc.kotlin.ServerCalls.unaryServerMethodDefinition
import io.grpc.kotlin.StubFor
import kotlin.String
import kotlin.coroutines.CoroutineContext
import kotlin.coroutines.EmptyCoroutineContext
import kotlin.jvm.JvmOverloads
import kotlin.jvm.JvmStatic
import kotlinx.coroutines.flow.Flow

/**
 * Holder for Kotlin coroutine-based client and server APIs for helixgitpx.collab.v1.CollabService.
 */
public object CollabServiceGrpcKt {
  public const val SERVICE_NAME: String = CollabServiceGrpc.SERVICE_NAME

  @JvmStatic
  public val serviceDescriptor: ServiceDescriptor
    get() = getServiceDescriptor()

  public val openDocMethod: MethodDescriptor<OpenDocRequest, Document>
    @JvmStatic
    get() = CollabServiceGrpc.getOpenDocMethod()

  public val applyOpsMethod: MethodDescriptor<DocOp, DocOp>
    @JvmStatic
    get() = CollabServiceGrpc.getApplyOpsMethod()

  public val listPresenceMethod: MethodDescriptor<ListPresenceRequest, ListPresenceResponse>
    @JvmStatic
    get() = CollabServiceGrpc.getListPresenceMethod()

  public val closeDocMethod: MethodDescriptor<CloseDocRequest, Empty>
    @JvmStatic
    get() = CollabServiceGrpc.getCloseDocMethod()

  /**
   * A stub for issuing RPCs to a(n) helixgitpx.collab.v1.CollabService service as suspending coroutines.
   */
  @StubFor(CollabServiceGrpc::class)
  public class CollabServiceCoroutineStub @JvmOverloads constructor(
    channel: Channel,
    callOptions: CallOptions = DEFAULT,
  ) : AbstractCoroutineStub<CollabServiceCoroutineStub>(channel, callOptions) {
    override fun build(channel: Channel, callOptions: CallOptions): CollabServiceCoroutineStub = CollabServiceCoroutineStub(channel, callOptions)

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
    public suspend fun openDoc(request: OpenDocRequest, headers: Metadata = Metadata()): Document = unaryRpc(
      channel,
      CollabServiceGrpc.getOpenDocMethod(),
      request,
      callOptions,
      headers
    )

    /**
     * Returns a [Flow] that, when collected, executes this RPC and emits responses from the
     * server as they arrive.  That flow finishes normally if the server closes its response with
     * [`Status.OK`][io.grpc.Status], and fails by throwing a [StatusException] otherwise.  If
     * collecting the flow downstream fails exceptionally (including via cancellation), the RPC
     * is cancelled with that exception as a cause.
     *
     * The [Flow] of requests is collected once each time the [Flow] of responses is
     * collected. If collection of the [Flow] of responses completes normally or
     * exceptionally before collection of `requests` completes, the collection of
     * `requests` is cancelled.  If the collection of `requests` completes
     * exceptionally for any other reason, then the collection of the [Flow] of responses
     * completes exceptionally for the same reason and the RPC is cancelled with that reason.
     *
     * @param requests A [Flow] of request messages.
     *
     * @param headers Metadata to attach to the request.  Most users will not need this.
     *
     * @return A flow that, when collected, emits the responses from the server.
     */
    public fun applyOps(requests: Flow<DocOp>, headers: Metadata = Metadata()): Flow<DocOp> = bidiStreamingRpc(
      channel,
      CollabServiceGrpc.getApplyOpsMethod(),
      requests,
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
    public suspend fun listPresence(request: ListPresenceRequest, headers: Metadata = Metadata()): ListPresenceResponse = unaryRpc(
      channel,
      CollabServiceGrpc.getListPresenceMethod(),
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
    public suspend fun closeDoc(request: CloseDocRequest, headers: Metadata = Metadata()): Empty = unaryRpc(
      channel,
      CollabServiceGrpc.getCloseDocMethod(),
      request,
      callOptions,
      headers
    )
  }

  /**
   * Skeletal implementation of the helixgitpx.collab.v1.CollabService service based on Kotlin coroutines.
   */
  public abstract class CollabServiceCoroutineImplBase(
    coroutineContext: CoroutineContext = EmptyCoroutineContext,
  ) : AbstractCoroutineServerImpl(coroutineContext) {
    /**
     * Returns the response to an RPC for helixgitpx.collab.v1.CollabService.OpenDoc.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun openDoc(request: OpenDocRequest): Document = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.collab.v1.CollabService.OpenDoc is unimplemented"))

    /**
     * Returns a [Flow] of responses to an RPC for helixgitpx.collab.v1.CollabService.ApplyOps.
     *
     * If creating or collecting the returned flow fails with a [StatusException], the RPC
     * will fail with the corresponding [io.grpc.Status].  If it fails with a
     * [java.util.concurrent.CancellationException], the RPC will fail with status `Status.CANCELLED`.  If creating
     * or collecting the returned flow fails for any other reason, the RPC will fail with
     * `Status.UNKNOWN` with the exception as a cause.
     *
     * @param requests A [Flow] of requests from the client.  This flow can be
     *        collected only once and throws [java.lang.IllegalStateException] on attempts to collect
     *        it more than once.
     */
    public open fun applyOps(requests: Flow<DocOp>): Flow<DocOp> = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.collab.v1.CollabService.ApplyOps is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.collab.v1.CollabService.ListPresence.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun listPresence(request: ListPresenceRequest): ListPresenceResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.collab.v1.CollabService.ListPresence is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.collab.v1.CollabService.CloseDoc.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun closeDoc(request: CloseDocRequest): Empty = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.collab.v1.CollabService.CloseDoc is unimplemented"))

    final override fun bindService(): ServerServiceDefinition = builder(getServiceDescriptor())
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = CollabServiceGrpc.getOpenDocMethod(),
      implementation = ::openDoc
    ))
      .addMethod(bidiStreamingServerMethodDefinition(
      context = this.context,
      descriptor = CollabServiceGrpc.getApplyOpsMethod(),
      implementation = ::applyOps
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = CollabServiceGrpc.getListPresenceMethod(),
      implementation = ::listPresence
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = CollabServiceGrpc.getCloseDocMethod(),
      implementation = ::closeDoc
    )).build()
  }
}
