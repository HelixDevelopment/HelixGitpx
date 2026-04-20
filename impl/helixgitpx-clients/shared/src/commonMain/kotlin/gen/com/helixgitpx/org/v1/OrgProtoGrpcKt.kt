package com.helixgitpx.org.v1

import com.google.protobuf.Empty
import com.helixgitpx.org.v1.OrgServiceGrpc.getServiceDescriptor
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
 * Holder for Kotlin coroutine-based client and server APIs for helixgitpx.org.v1.OrgService.
 */
public object OrgServiceGrpcKt {
  public const val SERVICE_NAME: String = OrgServiceGrpc.SERVICE_NAME

  @JvmStatic
  public val serviceDescriptor: ServiceDescriptor
    get() = getServiceDescriptor()

  public val createMethod: MethodDescriptor<CreateOrgRequest, Org>
    @JvmStatic
    get() = OrgServiceGrpc.getCreateMethod()

  public val getMethod: MethodDescriptor<GetOrgRequest, Org>
    @JvmStatic
    get() = OrgServiceGrpc.getGetMethod()

  public val listMethod: MethodDescriptor<ListOrgsRequest, ListOrgsResponse>
    @JvmStatic
    get() = OrgServiceGrpc.getListMethod()

  public val updateMethod: MethodDescriptor<UpdateOrgRequest, Org>
    @JvmStatic
    get() = OrgServiceGrpc.getUpdateMethod()

  public val deleteMethod: MethodDescriptor<DeleteOrgRequest, Empty>
    @JvmStatic
    get() = OrgServiceGrpc.getDeleteMethod()

  /**
   * A stub for issuing RPCs to a(n) helixgitpx.org.v1.OrgService service as suspending coroutines.
   */
  @StubFor(OrgServiceGrpc::class)
  public class OrgServiceCoroutineStub @JvmOverloads constructor(
    channel: Channel,
    callOptions: CallOptions = DEFAULT,
  ) : AbstractCoroutineStub<OrgServiceCoroutineStub>(channel, callOptions) {
    override fun build(channel: Channel, callOptions: CallOptions): OrgServiceCoroutineStub = OrgServiceCoroutineStub(channel, callOptions)

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
    public suspend fun create(request: CreateOrgRequest, headers: Metadata = Metadata()): Org = unaryRpc(
      channel,
      OrgServiceGrpc.getCreateMethod(),
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
    public suspend fun `get`(request: GetOrgRequest, headers: Metadata = Metadata()): Org = unaryRpc(
      channel,
      OrgServiceGrpc.getGetMethod(),
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
    public suspend fun list(request: ListOrgsRequest, headers: Metadata = Metadata()): ListOrgsResponse = unaryRpc(
      channel,
      OrgServiceGrpc.getListMethod(),
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
    public suspend fun update(request: UpdateOrgRequest, headers: Metadata = Metadata()): Org = unaryRpc(
      channel,
      OrgServiceGrpc.getUpdateMethod(),
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
    public suspend fun delete(request: DeleteOrgRequest, headers: Metadata = Metadata()): Empty = unaryRpc(
      channel,
      OrgServiceGrpc.getDeleteMethod(),
      request,
      callOptions,
      headers
    )
  }

  /**
   * Skeletal implementation of the helixgitpx.org.v1.OrgService service based on Kotlin coroutines.
   */
  public abstract class OrgServiceCoroutineImplBase(
    coroutineContext: CoroutineContext = EmptyCoroutineContext,
  ) : AbstractCoroutineServerImpl(coroutineContext) {
    /**
     * Returns the response to an RPC for helixgitpx.org.v1.OrgService.Create.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun create(request: CreateOrgRequest): Org = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.org.v1.OrgService.Create is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.org.v1.OrgService.Get.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun `get`(request: GetOrgRequest): Org = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.org.v1.OrgService.Get is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.org.v1.OrgService.List.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun list(request: ListOrgsRequest): ListOrgsResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.org.v1.OrgService.List is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.org.v1.OrgService.Update.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun update(request: UpdateOrgRequest): Org = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.org.v1.OrgService.Update is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.org.v1.OrgService.Delete.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun delete(request: DeleteOrgRequest): Empty = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.org.v1.OrgService.Delete is unimplemented"))

    final override fun bindService(): ServerServiceDefinition = builder(getServiceDescriptor())
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = OrgServiceGrpc.getCreateMethod(),
      implementation = ::create
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = OrgServiceGrpc.getGetMethod(),
      implementation = ::`get`
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = OrgServiceGrpc.getListMethod(),
      implementation = ::list
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = OrgServiceGrpc.getUpdateMethod(),
      implementation = ::update
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = OrgServiceGrpc.getDeleteMethod(),
      implementation = ::delete
    )).build()
  }
}
