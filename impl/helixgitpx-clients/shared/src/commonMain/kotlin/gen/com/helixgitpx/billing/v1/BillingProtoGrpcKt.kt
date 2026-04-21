package com.helixgitpx.billing.v1

import com.helixgitpx.billing.v1.BillingServiceGrpc.getServiceDescriptor
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
 * Holder for Kotlin coroutine-based client and server APIs for helixgitpx.billing.v1.BillingService.
 */
public object BillingServiceGrpcKt {
  public const val SERVICE_NAME: String = BillingServiceGrpc.SERVICE_NAME

  @JvmStatic
  public val serviceDescriptor: ServiceDescriptor
    get() = getServiceDescriptor()

  public val getSubscriptionMethod: MethodDescriptor<GetSubscriptionRequest, Subscription>
    @JvmStatic
    get() = BillingServiceGrpc.getGetSubscriptionMethod()

  public val upgradePlanMethod: MethodDescriptor<UpgradePlanRequest, Subscription>
    @JvmStatic
    get() = BillingServiceGrpc.getUpgradePlanMethod()

  public val downgradePlanMethod: MethodDescriptor<DowngradePlanRequest, Subscription>
    @JvmStatic
    get() = BillingServiceGrpc.getDowngradePlanMethod()

  public val cancelSubscriptionMethod: MethodDescriptor<CancelSubscriptionRequest, Subscription>
    @JvmStatic
    get() = BillingServiceGrpc.getCancelSubscriptionMethod()

  public val listInvoicesMethod: MethodDescriptor<ListInvoicesRequest, ListInvoicesResponse>
    @JvmStatic
    get() = BillingServiceGrpc.getListInvoicesMethod()

  public val downloadInvoiceMethod: MethodDescriptor<DownloadInvoiceRequest, InvoiceDownload>
    @JvmStatic
    get() = BillingServiceGrpc.getDownloadInvoiceMethod()

  /**
   * A stub for issuing RPCs to a(n) helixgitpx.billing.v1.BillingService service as suspending coroutines.
   */
  @StubFor(BillingServiceGrpc::class)
  public class BillingServiceCoroutineStub @JvmOverloads constructor(
    channel: Channel,
    callOptions: CallOptions = DEFAULT,
  ) : AbstractCoroutineStub<BillingServiceCoroutineStub>(channel, callOptions) {
    override fun build(channel: Channel, callOptions: CallOptions): BillingServiceCoroutineStub = BillingServiceCoroutineStub(channel, callOptions)

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
    public suspend fun getSubscription(request: GetSubscriptionRequest, headers: Metadata = Metadata()): Subscription = unaryRpc(
      channel,
      BillingServiceGrpc.getGetSubscriptionMethod(),
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
    public suspend fun upgradePlan(request: UpgradePlanRequest, headers: Metadata = Metadata()): Subscription = unaryRpc(
      channel,
      BillingServiceGrpc.getUpgradePlanMethod(),
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
    public suspend fun downgradePlan(request: DowngradePlanRequest, headers: Metadata = Metadata()): Subscription = unaryRpc(
      channel,
      BillingServiceGrpc.getDowngradePlanMethod(),
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
    public suspend fun cancelSubscription(request: CancelSubscriptionRequest, headers: Metadata = Metadata()): Subscription = unaryRpc(
      channel,
      BillingServiceGrpc.getCancelSubscriptionMethod(),
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
    public suspend fun listInvoices(request: ListInvoicesRequest, headers: Metadata = Metadata()): ListInvoicesResponse = unaryRpc(
      channel,
      BillingServiceGrpc.getListInvoicesMethod(),
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
    public suspend fun downloadInvoice(request: DownloadInvoiceRequest, headers: Metadata = Metadata()): InvoiceDownload = unaryRpc(
      channel,
      BillingServiceGrpc.getDownloadInvoiceMethod(),
      request,
      callOptions,
      headers
    )
  }

  /**
   * Skeletal implementation of the helixgitpx.billing.v1.BillingService service based on Kotlin coroutines.
   */
  public abstract class BillingServiceCoroutineImplBase(
    coroutineContext: CoroutineContext = EmptyCoroutineContext,
  ) : AbstractCoroutineServerImpl(coroutineContext) {
    /**
     * Returns the response to an RPC for helixgitpx.billing.v1.BillingService.GetSubscription.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun getSubscription(request: GetSubscriptionRequest): Subscription = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.billing.v1.BillingService.GetSubscription is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.billing.v1.BillingService.UpgradePlan.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun upgradePlan(request: UpgradePlanRequest): Subscription = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.billing.v1.BillingService.UpgradePlan is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.billing.v1.BillingService.DowngradePlan.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun downgradePlan(request: DowngradePlanRequest): Subscription = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.billing.v1.BillingService.DowngradePlan is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.billing.v1.BillingService.CancelSubscription.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun cancelSubscription(request: CancelSubscriptionRequest): Subscription = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.billing.v1.BillingService.CancelSubscription is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.billing.v1.BillingService.ListInvoices.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun listInvoices(request: ListInvoicesRequest): ListInvoicesResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.billing.v1.BillingService.ListInvoices is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.billing.v1.BillingService.DownloadInvoice.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun downloadInvoice(request: DownloadInvoiceRequest): InvoiceDownload = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.billing.v1.BillingService.DownloadInvoice is unimplemented"))

    final override fun bindService(): ServerServiceDefinition = builder(getServiceDescriptor())
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = BillingServiceGrpc.getGetSubscriptionMethod(),
      implementation = ::getSubscription
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = BillingServiceGrpc.getUpgradePlanMethod(),
      implementation = ::upgradePlan
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = BillingServiceGrpc.getDowngradePlanMethod(),
      implementation = ::downgradePlan
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = BillingServiceGrpc.getCancelSubscriptionMethod(),
      implementation = ::cancelSubscription
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = BillingServiceGrpc.getListInvoicesMethod(),
      implementation = ::listInvoices
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = BillingServiceGrpc.getDownloadInvoiceMethod(),
      implementation = ::downloadInvoice
    )).build()
  }
}
