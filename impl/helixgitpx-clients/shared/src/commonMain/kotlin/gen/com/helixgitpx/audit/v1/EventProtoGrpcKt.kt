package com.helixgitpx.audit.v1

import com.helixgitpx.audit.v1.AuditServiceGrpc.getServiceDescriptor
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
 * Holder for Kotlin coroutine-based client and server APIs for helixgitpx.audit.v1.AuditService.
 */
public object AuditServiceGrpcKt {
  public const val SERVICE_NAME: String = AuditServiceGrpc.SERVICE_NAME

  @JvmStatic
  public val serviceDescriptor: ServiceDescriptor
    get() = getServiceDescriptor()

  public val listEventsMethod: MethodDescriptor<ListEventsRequest, ListEventsResponse>
    @JvmStatic
    get() = AuditServiceGrpc.getListEventsMethod()

  public val getEventMethod: MethodDescriptor<GetEventRequest, AuditEvent>
    @JvmStatic
    get() = AuditServiceGrpc.getGetEventMethod()

  public val verifyMerkleRootMethod:
      MethodDescriptor<VerifyMerkleRootRequest, VerifyMerkleRootResponse>
    @JvmStatic
    get() = AuditServiceGrpc.getVerifyMerkleRootMethod()

  /**
   * A stub for issuing RPCs to a(n) helixgitpx.audit.v1.AuditService service as suspending coroutines.
   */
  @StubFor(AuditServiceGrpc::class)
  public class AuditServiceCoroutineStub @JvmOverloads constructor(
    channel: Channel,
    callOptions: CallOptions = DEFAULT,
  ) : AbstractCoroutineStub<AuditServiceCoroutineStub>(channel, callOptions) {
    override fun build(channel: Channel, callOptions: CallOptions): AuditServiceCoroutineStub = AuditServiceCoroutineStub(channel, callOptions)

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
    public suspend fun listEvents(request: ListEventsRequest, headers: Metadata = Metadata()): ListEventsResponse = unaryRpc(
      channel,
      AuditServiceGrpc.getListEventsMethod(),
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
    public suspend fun getEvent(request: GetEventRequest, headers: Metadata = Metadata()): AuditEvent = unaryRpc(
      channel,
      AuditServiceGrpc.getGetEventMethod(),
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
    public suspend fun verifyMerkleRoot(request: VerifyMerkleRootRequest, headers: Metadata = Metadata()): VerifyMerkleRootResponse = unaryRpc(
      channel,
      AuditServiceGrpc.getVerifyMerkleRootMethod(),
      request,
      callOptions,
      headers
    )
  }

  /**
   * Skeletal implementation of the helixgitpx.audit.v1.AuditService service based on Kotlin coroutines.
   */
  public abstract class AuditServiceCoroutineImplBase(
    coroutineContext: CoroutineContext = EmptyCoroutineContext,
  ) : AbstractCoroutineServerImpl(coroutineContext) {
    /**
     * Returns the response to an RPC for helixgitpx.audit.v1.AuditService.ListEvents.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun listEvents(request: ListEventsRequest): ListEventsResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.audit.v1.AuditService.ListEvents is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.audit.v1.AuditService.GetEvent.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun getEvent(request: GetEventRequest): AuditEvent = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.audit.v1.AuditService.GetEvent is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.audit.v1.AuditService.VerifyMerkleRoot.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun verifyMerkleRoot(request: VerifyMerkleRootRequest): VerifyMerkleRootResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.audit.v1.AuditService.VerifyMerkleRoot is unimplemented"))

    final override fun bindService(): ServerServiceDefinition = builder(getServiceDescriptor())
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = AuditServiceGrpc.getListEventsMethod(),
      implementation = ::listEvents
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = AuditServiceGrpc.getGetEventMethod(),
      implementation = ::getEvent
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = AuditServiceGrpc.getVerifyMerkleRootMethod(),
      implementation = ::verifyMerkleRoot
    )).build()
  }
}
