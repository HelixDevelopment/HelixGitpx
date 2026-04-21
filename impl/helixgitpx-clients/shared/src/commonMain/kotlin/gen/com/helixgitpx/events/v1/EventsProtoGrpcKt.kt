package com.helixgitpx.events.v1

import com.helixgitpx.events.v1.LiveEventsServiceGrpc.getServiceDescriptor
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
import io.grpc.kotlin.ClientCalls.serverStreamingRpc
import io.grpc.kotlin.ClientCalls.unaryRpc
import io.grpc.kotlin.ServerCalls.serverStreamingServerMethodDefinition
import io.grpc.kotlin.ServerCalls.unaryServerMethodDefinition
import io.grpc.kotlin.StubFor
import kotlin.String
import kotlin.coroutines.CoroutineContext
import kotlin.coroutines.EmptyCoroutineContext
import kotlin.jvm.JvmOverloads
import kotlin.jvm.JvmStatic
import kotlinx.coroutines.flow.Flow

/**
 * Holder for Kotlin coroutine-based client and server APIs for helixgitpx.events.v1.LiveEventsService.
 */
public object LiveEventsServiceGrpcKt {
  public const val SERVICE_NAME: String = LiveEventsServiceGrpc.SERVICE_NAME

  @JvmStatic
  public val serviceDescriptor: ServiceDescriptor
    get() = getServiceDescriptor()

  public val subscribeMethod: MethodDescriptor<SubscribeRequest, Event>
    @JvmStatic
    get() = LiveEventsServiceGrpc.getSubscribeMethod()

  public val resumeMethod: MethodDescriptor<ResumeRequest, Event>
    @JvmStatic
    get() = LiveEventsServiceGrpc.getResumeMethod()

  public val publishMethod: MethodDescriptor<PublishRequest, PublishResponse>
    @JvmStatic
    get() = LiveEventsServiceGrpc.getPublishMethod()

  /**
   * A stub for issuing RPCs to a(n) helixgitpx.events.v1.LiveEventsService service as suspending coroutines.
   */
  @StubFor(LiveEventsServiceGrpc::class)
  public class LiveEventsServiceCoroutineStub @JvmOverloads constructor(
    channel: Channel,
    callOptions: CallOptions = DEFAULT,
  ) : AbstractCoroutineStub<LiveEventsServiceCoroutineStub>(channel, callOptions) {
    override fun build(channel: Channel, callOptions: CallOptions): LiveEventsServiceCoroutineStub = LiveEventsServiceCoroutineStub(channel, callOptions)

    /**
     * Returns a [Flow] that, when collected, executes this RPC and emits responses from the
     * server as they arrive.  That flow finishes normally if the server closes its response with
     * [`Status.OK`][io.grpc.Status], and fails by throwing a [StatusException] otherwise.  If
     * collecting the flow downstream fails exceptionally (including via cancellation), the RPC
     * is cancelled with that exception as a cause.
     *
     * @param request The request message to send to the server.
     *
     * @param headers Metadata to attach to the request.  Most users will not need this.
     *
     * @return A flow that, when collected, emits the responses from the server.
     */
    public fun subscribe(request: SubscribeRequest, headers: Metadata = Metadata()): Flow<Event> = serverStreamingRpc(
      channel,
      LiveEventsServiceGrpc.getSubscribeMethod(),
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
     * @param request The request message to send to the server.
     *
     * @param headers Metadata to attach to the request.  Most users will not need this.
     *
     * @return A flow that, when collected, emits the responses from the server.
     */
    public fun resume(request: ResumeRequest, headers: Metadata = Metadata()): Flow<Event> = serverStreamingRpc(
      channel,
      LiveEventsServiceGrpc.getResumeMethod(),
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
    public suspend fun publish(request: PublishRequest, headers: Metadata = Metadata()): PublishResponse = unaryRpc(
      channel,
      LiveEventsServiceGrpc.getPublishMethod(),
      request,
      callOptions,
      headers
    )
  }

  /**
   * Skeletal implementation of the helixgitpx.events.v1.LiveEventsService service based on Kotlin coroutines.
   */
  public abstract class LiveEventsServiceCoroutineImplBase(
    coroutineContext: CoroutineContext = EmptyCoroutineContext,
  ) : AbstractCoroutineServerImpl(coroutineContext) {
    /**
     * Returns a [Flow] of responses to an RPC for helixgitpx.events.v1.LiveEventsService.Subscribe.
     *
     * If creating or collecting the returned flow fails with a [StatusException], the RPC
     * will fail with the corresponding [io.grpc.Status].  If it fails with a
     * [java.util.concurrent.CancellationException], the RPC will fail with status `Status.CANCELLED`.  If creating
     * or collecting the returned flow fails for any other reason, the RPC will fail with
     * `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open fun subscribe(request: SubscribeRequest): Flow<Event> = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.events.v1.LiveEventsService.Subscribe is unimplemented"))

    /**
     * Returns a [Flow] of responses to an RPC for helixgitpx.events.v1.LiveEventsService.Resume.
     *
     * If creating or collecting the returned flow fails with a [StatusException], the RPC
     * will fail with the corresponding [io.grpc.Status].  If it fails with a
     * [java.util.concurrent.CancellationException], the RPC will fail with status `Status.CANCELLED`.  If creating
     * or collecting the returned flow fails for any other reason, the RPC will fail with
     * `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open fun resume(request: ResumeRequest): Flow<Event> = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.events.v1.LiveEventsService.Resume is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.events.v1.LiveEventsService.Publish.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun publish(request: PublishRequest): PublishResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.events.v1.LiveEventsService.Publish is unimplemented"))

    final override fun bindService(): ServerServiceDefinition = builder(getServiceDescriptor())
      .addMethod(serverStreamingServerMethodDefinition(
      context = this.context,
      descriptor = LiveEventsServiceGrpc.getSubscribeMethod(),
      implementation = ::subscribe
    ))
      .addMethod(serverStreamingServerMethodDefinition(
      context = this.context,
      descriptor = LiveEventsServiceGrpc.getResumeMethod(),
      implementation = ::resume
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = LiveEventsServiceGrpc.getPublishMethod(),
      implementation = ::publish
    )).build()
  }
}
