package com.helixgitpx.upstream.v1

import com.google.protobuf.Empty
import com.helixgitpx.upstream.v1.UpstreamServiceGrpc.getServiceDescriptor
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
 * Holder for Kotlin coroutine-based client and server APIs for helixgitpx.upstream.v1.UpstreamService.
 */
public object UpstreamServiceGrpcKt {
  public const val SERVICE_NAME: String = UpstreamServiceGrpc.SERVICE_NAME

  @JvmStatic
  public val serviceDescriptor: ServiceDescriptor
    get() = getServiceDescriptor()

  public val createMethod: MethodDescriptor<CreateUpstreamRequest, Upstream>
    @JvmStatic
    get() = UpstreamServiceGrpc.getCreateMethod()

  public val getMethod: MethodDescriptor<GetUpstreamRequest, Upstream>
    @JvmStatic
    get() = UpstreamServiceGrpc.getGetMethod()

  public val listMethod: MethodDescriptor<Empty, ListUpstreamsResponse>
    @JvmStatic
    get() = UpstreamServiceGrpc.getListMethod()

  public val updateMethod: MethodDescriptor<UpdateUpstreamRequest, Upstream>
    @JvmStatic
    get() = UpstreamServiceGrpc.getUpdateMethod()

  public val deleteMethod: MethodDescriptor<DeleteUpstreamRequest, Empty>
    @JvmStatic
    get() = UpstreamServiceGrpc.getDeleteMethod()

  public val bindMethod: MethodDescriptor<BindRequest, Binding>
    @JvmStatic
    get() = UpstreamServiceGrpc.getBindMethod()

  public val unbindMethod: MethodDescriptor<UnbindRequest, Empty>
    @JvmStatic
    get() = UpstreamServiceGrpc.getUnbindMethod()

  public val listBindingsMethod: MethodDescriptor<ListBindingsRequest, ListBindingsResponse>
    @JvmStatic
    get() = UpstreamServiceGrpc.getListBindingsMethod()

  /**
   * A stub for issuing RPCs to a(n) helixgitpx.upstream.v1.UpstreamService service as suspending coroutines.
   */
  @StubFor(UpstreamServiceGrpc::class)
  public class UpstreamServiceCoroutineStub @JvmOverloads constructor(
    channel: Channel,
    callOptions: CallOptions = DEFAULT,
  ) : AbstractCoroutineStub<UpstreamServiceCoroutineStub>(channel, callOptions) {
    override fun build(channel: Channel, callOptions: CallOptions): UpstreamServiceCoroutineStub = UpstreamServiceCoroutineStub(channel, callOptions)

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
    public suspend fun create(request: CreateUpstreamRequest, headers: Metadata = Metadata()): Upstream = unaryRpc(
      channel,
      UpstreamServiceGrpc.getCreateMethod(),
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
    public suspend fun `get`(request: GetUpstreamRequest, headers: Metadata = Metadata()): Upstream = unaryRpc(
      channel,
      UpstreamServiceGrpc.getGetMethod(),
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
    public suspend fun list(request: Empty, headers: Metadata = Metadata()): ListUpstreamsResponse = unaryRpc(
      channel,
      UpstreamServiceGrpc.getListMethod(),
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
    public suspend fun update(request: UpdateUpstreamRequest, headers: Metadata = Metadata()): Upstream = unaryRpc(
      channel,
      UpstreamServiceGrpc.getUpdateMethod(),
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
    public suspend fun delete(request: DeleteUpstreamRequest, headers: Metadata = Metadata()): Empty = unaryRpc(
      channel,
      UpstreamServiceGrpc.getDeleteMethod(),
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
    public suspend fun bind(request: BindRequest, headers: Metadata = Metadata()): Binding = unaryRpc(
      channel,
      UpstreamServiceGrpc.getBindMethod(),
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
    public suspend fun unbind(request: UnbindRequest, headers: Metadata = Metadata()): Empty = unaryRpc(
      channel,
      UpstreamServiceGrpc.getUnbindMethod(),
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
    public suspend fun listBindings(request: ListBindingsRequest, headers: Metadata = Metadata()): ListBindingsResponse = unaryRpc(
      channel,
      UpstreamServiceGrpc.getListBindingsMethod(),
      request,
      callOptions,
      headers
    )
  }

  /**
   * Skeletal implementation of the helixgitpx.upstream.v1.UpstreamService service based on Kotlin coroutines.
   */
  public abstract class UpstreamServiceCoroutineImplBase(
    coroutineContext: CoroutineContext = EmptyCoroutineContext,
  ) : AbstractCoroutineServerImpl(coroutineContext) {
    /**
     * Returns the response to an RPC for helixgitpx.upstream.v1.UpstreamService.Create.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun create(request: CreateUpstreamRequest): Upstream = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.upstream.v1.UpstreamService.Create is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.upstream.v1.UpstreamService.Get.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun `get`(request: GetUpstreamRequest): Upstream = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.upstream.v1.UpstreamService.Get is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.upstream.v1.UpstreamService.List.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun list(request: Empty): ListUpstreamsResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.upstream.v1.UpstreamService.List is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.upstream.v1.UpstreamService.Update.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun update(request: UpdateUpstreamRequest): Upstream = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.upstream.v1.UpstreamService.Update is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.upstream.v1.UpstreamService.Delete.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun delete(request: DeleteUpstreamRequest): Empty = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.upstream.v1.UpstreamService.Delete is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.upstream.v1.UpstreamService.Bind.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun bind(request: BindRequest): Binding = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.upstream.v1.UpstreamService.Bind is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.upstream.v1.UpstreamService.Unbind.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun unbind(request: UnbindRequest): Empty = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.upstream.v1.UpstreamService.Unbind is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.upstream.v1.UpstreamService.ListBindings.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun listBindings(request: ListBindingsRequest): ListBindingsResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.upstream.v1.UpstreamService.ListBindings is unimplemented"))

    final override fun bindService(): ServerServiceDefinition = builder(getServiceDescriptor())
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = UpstreamServiceGrpc.getCreateMethod(),
      implementation = ::create
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = UpstreamServiceGrpc.getGetMethod(),
      implementation = ::`get`
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = UpstreamServiceGrpc.getListMethod(),
      implementation = ::list
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = UpstreamServiceGrpc.getUpdateMethod(),
      implementation = ::update
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = UpstreamServiceGrpc.getDeleteMethod(),
      implementation = ::delete
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = UpstreamServiceGrpc.getBindMethod(),
      implementation = ::bind
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = UpstreamServiceGrpc.getUnbindMethod(),
      implementation = ::unbind
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = UpstreamServiceGrpc.getListBindingsMethod(),
      implementation = ::listBindings
    )).build()
  }
}
