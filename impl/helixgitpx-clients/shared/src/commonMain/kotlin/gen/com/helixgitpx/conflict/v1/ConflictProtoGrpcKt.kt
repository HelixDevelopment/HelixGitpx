package com.helixgitpx.conflict.v1

import com.helixgitpx.conflict.v1.ConflictServiceGrpc.getServiceDescriptor
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
 * Holder for Kotlin coroutine-based client and server APIs for helixgitpx.conflict.v1.ConflictService.
 */
public object ConflictServiceGrpcKt {
  public const val SERVICE_NAME: String = ConflictServiceGrpc.SERVICE_NAME

  @JvmStatic
  public val serviceDescriptor: ServiceDescriptor
    get() = getServiceDescriptor()

  public val listMethod: MethodDescriptor<ListConflictsRequest, ListConflictsResponse>
    @JvmStatic
    get() = ConflictServiceGrpc.getListMethod()

  public val getMethod: MethodDescriptor<GetConflictRequest, Conflict>
    @JvmStatic
    get() = ConflictServiceGrpc.getGetMethod()

  public val proposeResolutionMethod: MethodDescriptor<ProposeResolutionRequest, Resolution>
    @JvmStatic
    get() = ConflictServiceGrpc.getProposeResolutionMethod()

  public val acceptResolutionMethod: MethodDescriptor<AcceptResolutionRequest, Conflict>
    @JvmStatic
    get() = ConflictServiceGrpc.getAcceptResolutionMethod()

  public val rejectResolutionMethod: MethodDescriptor<RejectResolutionRequest, Conflict>
    @JvmStatic
    get() = ConflictServiceGrpc.getRejectResolutionMethod()

  /**
   * A stub for issuing RPCs to a(n) helixgitpx.conflict.v1.ConflictService service as suspending coroutines.
   */
  @StubFor(ConflictServiceGrpc::class)
  public class ConflictServiceCoroutineStub @JvmOverloads constructor(
    channel: Channel,
    callOptions: CallOptions = DEFAULT,
  ) : AbstractCoroutineStub<ConflictServiceCoroutineStub>(channel, callOptions) {
    override fun build(channel: Channel, callOptions: CallOptions): ConflictServiceCoroutineStub = ConflictServiceCoroutineStub(channel, callOptions)

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
    public suspend fun list(request: ListConflictsRequest, headers: Metadata = Metadata()): ListConflictsResponse = unaryRpc(
      channel,
      ConflictServiceGrpc.getListMethod(),
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
    public suspend fun `get`(request: GetConflictRequest, headers: Metadata = Metadata()): Conflict = unaryRpc(
      channel,
      ConflictServiceGrpc.getGetMethod(),
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
    public suspend fun proposeResolution(request: ProposeResolutionRequest, headers: Metadata = Metadata()): Resolution = unaryRpc(
      channel,
      ConflictServiceGrpc.getProposeResolutionMethod(),
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
    public suspend fun acceptResolution(request: AcceptResolutionRequest, headers: Metadata = Metadata()): Conflict = unaryRpc(
      channel,
      ConflictServiceGrpc.getAcceptResolutionMethod(),
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
    public suspend fun rejectResolution(request: RejectResolutionRequest, headers: Metadata = Metadata()): Conflict = unaryRpc(
      channel,
      ConflictServiceGrpc.getRejectResolutionMethod(),
      request,
      callOptions,
      headers
    )
  }

  /**
   * Skeletal implementation of the helixgitpx.conflict.v1.ConflictService service based on Kotlin coroutines.
   */
  public abstract class ConflictServiceCoroutineImplBase(
    coroutineContext: CoroutineContext = EmptyCoroutineContext,
  ) : AbstractCoroutineServerImpl(coroutineContext) {
    /**
     * Returns the response to an RPC for helixgitpx.conflict.v1.ConflictService.List.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun list(request: ListConflictsRequest): ListConflictsResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.conflict.v1.ConflictService.List is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.conflict.v1.ConflictService.Get.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun `get`(request: GetConflictRequest): Conflict = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.conflict.v1.ConflictService.Get is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.conflict.v1.ConflictService.ProposeResolution.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun proposeResolution(request: ProposeResolutionRequest): Resolution = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.conflict.v1.ConflictService.ProposeResolution is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.conflict.v1.ConflictService.AcceptResolution.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun acceptResolution(request: AcceptResolutionRequest): Conflict = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.conflict.v1.ConflictService.AcceptResolution is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.conflict.v1.ConflictService.RejectResolution.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun rejectResolution(request: RejectResolutionRequest): Conflict = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.conflict.v1.ConflictService.RejectResolution is unimplemented"))

    final override fun bindService(): ServerServiceDefinition = builder(getServiceDescriptor())
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = ConflictServiceGrpc.getListMethod(),
      implementation = ::list
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = ConflictServiceGrpc.getGetMethod(),
      implementation = ::`get`
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = ConflictServiceGrpc.getProposeResolutionMethod(),
      implementation = ::proposeResolution
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = ConflictServiceGrpc.getAcceptResolutionMethod(),
      implementation = ::acceptResolution
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = ConflictServiceGrpc.getRejectResolutionMethod(),
      implementation = ::rejectResolution
    )).build()
  }
}
