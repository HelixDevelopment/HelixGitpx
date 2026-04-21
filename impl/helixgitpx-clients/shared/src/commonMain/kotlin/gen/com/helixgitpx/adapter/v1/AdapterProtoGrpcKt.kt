package com.helixgitpx.adapter.v1

import com.helixgitpx.adapter.v1.AdapterPoolServiceGrpc.getServiceDescriptor
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
 * Holder for Kotlin coroutine-based client and server APIs for helixgitpx.adapter.v1.AdapterPoolService.
 */
public object AdapterPoolServiceGrpcKt {
  public const val SERVICE_NAME: String = AdapterPoolServiceGrpc.SERVICE_NAME

  @JvmStatic
  public val serviceDescriptor: ServiceDescriptor
    get() = getServiceDescriptor()

  public val listProvidersMethod: MethodDescriptor<ListProvidersRequest, ListProvidersResponse>
    @JvmStatic
    get() = AdapterPoolServiceGrpc.getListProvidersMethod()

  public val getProviderHealthMethod: MethodDescriptor<ProviderHealthRequest, ProviderHealth>
    @JvmStatic
    get() = AdapterPoolServiceGrpc.getGetProviderHealthMethod()

  public val rotateTokenMethod: MethodDescriptor<RotateTokenRequest, RotateTokenResponse>
    @JvmStatic
    get() = AdapterPoolServiceGrpc.getRotateTokenMethod()

  /**
   * A stub for issuing RPCs to a(n) helixgitpx.adapter.v1.AdapterPoolService service as suspending coroutines.
   */
  @StubFor(AdapterPoolServiceGrpc::class)
  public class AdapterPoolServiceCoroutineStub @JvmOverloads constructor(
    channel: Channel,
    callOptions: CallOptions = DEFAULT,
  ) : AbstractCoroutineStub<AdapterPoolServiceCoroutineStub>(channel, callOptions) {
    override fun build(channel: Channel, callOptions: CallOptions): AdapterPoolServiceCoroutineStub = AdapterPoolServiceCoroutineStub(channel, callOptions)

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
    public suspend fun listProviders(request: ListProvidersRequest, headers: Metadata = Metadata()): ListProvidersResponse = unaryRpc(
      channel,
      AdapterPoolServiceGrpc.getListProvidersMethod(),
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
    public suspend fun getProviderHealth(request: ProviderHealthRequest, headers: Metadata = Metadata()): ProviderHealth = unaryRpc(
      channel,
      AdapterPoolServiceGrpc.getGetProviderHealthMethod(),
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
    public suspend fun rotateToken(request: RotateTokenRequest, headers: Metadata = Metadata()): RotateTokenResponse = unaryRpc(
      channel,
      AdapterPoolServiceGrpc.getRotateTokenMethod(),
      request,
      callOptions,
      headers
    )
  }

  /**
   * Skeletal implementation of the helixgitpx.adapter.v1.AdapterPoolService service based on Kotlin coroutines.
   */
  public abstract class AdapterPoolServiceCoroutineImplBase(
    coroutineContext: CoroutineContext = EmptyCoroutineContext,
  ) : AbstractCoroutineServerImpl(coroutineContext) {
    /**
     * Returns the response to an RPC for helixgitpx.adapter.v1.AdapterPoolService.ListProviders.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun listProviders(request: ListProvidersRequest): ListProvidersResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.adapter.v1.AdapterPoolService.ListProviders is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.adapter.v1.AdapterPoolService.GetProviderHealth.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun getProviderHealth(request: ProviderHealthRequest): ProviderHealth = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.adapter.v1.AdapterPoolService.GetProviderHealth is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.adapter.v1.AdapterPoolService.RotateToken.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun rotateToken(request: RotateTokenRequest): RotateTokenResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.adapter.v1.AdapterPoolService.RotateToken is unimplemented"))

    final override fun bindService(): ServerServiceDefinition = builder(getServiceDescriptor())
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = AdapterPoolServiceGrpc.getListProvidersMethod(),
      implementation = ::listProviders
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = AdapterPoolServiceGrpc.getGetProviderHealthMethod(),
      implementation = ::getProviderHealth
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = AdapterPoolServiceGrpc.getRotateTokenMethod(),
      implementation = ::rotateToken
    )).build()
  }
}
