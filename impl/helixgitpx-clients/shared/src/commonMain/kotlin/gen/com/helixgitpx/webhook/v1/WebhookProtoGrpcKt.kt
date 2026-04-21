package com.helixgitpx.webhook.v1

import com.google.protobuf.Empty
import com.helixgitpx.webhook.v1.WebhookGatewayServiceGrpc.getServiceDescriptor
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
 * Holder for Kotlin coroutine-based client and server APIs for helixgitpx.webhook.v1.WebhookGatewayService.
 */
public object WebhookGatewayServiceGrpcKt {
  public const val SERVICE_NAME: String = WebhookGatewayServiceGrpc.SERVICE_NAME

  @JvmStatic
  public val serviceDescriptor: ServiceDescriptor
    get() = getServiceDescriptor()

  public val listDeliveriesMethod: MethodDescriptor<ListDeliveriesRequest, ListDeliveriesResponse>
    @JvmStatic
    get() = WebhookGatewayServiceGrpc.getListDeliveriesMethod()

  public val getDeliveryMethod: MethodDescriptor<GetDeliveryRequest, Delivery>
    @JvmStatic
    get() = WebhookGatewayServiceGrpc.getGetDeliveryMethod()

  public val replayMethod: MethodDescriptor<ReplayRequest, Delivery>
    @JvmStatic
    get() = WebhookGatewayServiceGrpc.getReplayMethod()

  public val rotateSecretMethod: MethodDescriptor<RotateSecretRequest, Empty>
    @JvmStatic
    get() = WebhookGatewayServiceGrpc.getRotateSecretMethod()

  /**
   * A stub for issuing RPCs to a(n) helixgitpx.webhook.v1.WebhookGatewayService service as suspending coroutines.
   */
  @StubFor(WebhookGatewayServiceGrpc::class)
  public class WebhookGatewayServiceCoroutineStub @JvmOverloads constructor(
    channel: Channel,
    callOptions: CallOptions = DEFAULT,
  ) : AbstractCoroutineStub<WebhookGatewayServiceCoroutineStub>(channel, callOptions) {
    override fun build(channel: Channel, callOptions: CallOptions): WebhookGatewayServiceCoroutineStub = WebhookGatewayServiceCoroutineStub(channel, callOptions)

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
    public suspend fun listDeliveries(request: ListDeliveriesRequest, headers: Metadata = Metadata()): ListDeliveriesResponse = unaryRpc(
      channel,
      WebhookGatewayServiceGrpc.getListDeliveriesMethod(),
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
    public suspend fun getDelivery(request: GetDeliveryRequest, headers: Metadata = Metadata()): Delivery = unaryRpc(
      channel,
      WebhookGatewayServiceGrpc.getGetDeliveryMethod(),
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
    public suspend fun replay(request: ReplayRequest, headers: Metadata = Metadata()): Delivery = unaryRpc(
      channel,
      WebhookGatewayServiceGrpc.getReplayMethod(),
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
    public suspend fun rotateSecret(request: RotateSecretRequest, headers: Metadata = Metadata()): Empty = unaryRpc(
      channel,
      WebhookGatewayServiceGrpc.getRotateSecretMethod(),
      request,
      callOptions,
      headers
    )
  }

  /**
   * Skeletal implementation of the helixgitpx.webhook.v1.WebhookGatewayService service based on Kotlin coroutines.
   */
  public abstract class WebhookGatewayServiceCoroutineImplBase(
    coroutineContext: CoroutineContext = EmptyCoroutineContext,
  ) : AbstractCoroutineServerImpl(coroutineContext) {
    /**
     * Returns the response to an RPC for helixgitpx.webhook.v1.WebhookGatewayService.ListDeliveries.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun listDeliveries(request: ListDeliveriesRequest): ListDeliveriesResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.webhook.v1.WebhookGatewayService.ListDeliveries is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.webhook.v1.WebhookGatewayService.GetDelivery.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun getDelivery(request: GetDeliveryRequest): Delivery = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.webhook.v1.WebhookGatewayService.GetDelivery is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.webhook.v1.WebhookGatewayService.Replay.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun replay(request: ReplayRequest): Delivery = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.webhook.v1.WebhookGatewayService.Replay is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.webhook.v1.WebhookGatewayService.RotateSecret.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun rotateSecret(request: RotateSecretRequest): Empty = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.webhook.v1.WebhookGatewayService.RotateSecret is unimplemented"))

    final override fun bindService(): ServerServiceDefinition = builder(getServiceDescriptor())
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = WebhookGatewayServiceGrpc.getListDeliveriesMethod(),
      implementation = ::listDeliveries
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = WebhookGatewayServiceGrpc.getGetDeliveryMethod(),
      implementation = ::getDelivery
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = WebhookGatewayServiceGrpc.getReplayMethod(),
      implementation = ::replay
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = WebhookGatewayServiceGrpc.getRotateSecretMethod(),
      implementation = ::rotateSecret
    )).build()
  }
}
