package com.helixgitpx.policy.v1

import com.helixgitpx.policy.v1.PolicyBundleServiceGrpc.getServiceDescriptor
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
 * Holder for Kotlin coroutine-based client and server APIs for helixgitpx.policy.v1.PolicyBundleService.
 */
public object PolicyBundleServiceGrpcKt {
  public const val SERVICE_NAME: String = PolicyBundleServiceGrpc.SERVICE_NAME

  @JvmStatic
  public val serviceDescriptor: ServiceDescriptor
    get() = getServiceDescriptor()

  public val listBundlesMethod: MethodDescriptor<ListBundlesRequest, ListBundlesResponse>
    @JvmStatic
    get() = PolicyBundleServiceGrpc.getListBundlesMethod()

  public val getActiveMethod: MethodDescriptor<GetActiveRequest, Bundle>
    @JvmStatic
    get() = PolicyBundleServiceGrpc.getGetActiveMethod()

  public val activateMethod: MethodDescriptor<ActivateRequest, Bundle>
    @JvmStatic
    get() = PolicyBundleServiceGrpc.getActivateMethod()

  public val rollbackMethod: MethodDescriptor<RollbackRequest, Bundle>
    @JvmStatic
    get() = PolicyBundleServiceGrpc.getRollbackMethod()

  /**
   * A stub for issuing RPCs to a(n) helixgitpx.policy.v1.PolicyBundleService service as suspending coroutines.
   */
  @StubFor(PolicyBundleServiceGrpc::class)
  public class PolicyBundleServiceCoroutineStub @JvmOverloads constructor(
    channel: Channel,
    callOptions: CallOptions = DEFAULT,
  ) : AbstractCoroutineStub<PolicyBundleServiceCoroutineStub>(channel, callOptions) {
    override fun build(channel: Channel, callOptions: CallOptions): PolicyBundleServiceCoroutineStub = PolicyBundleServiceCoroutineStub(channel, callOptions)

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
    public suspend fun listBundles(request: ListBundlesRequest, headers: Metadata = Metadata()): ListBundlesResponse = unaryRpc(
      channel,
      PolicyBundleServiceGrpc.getListBundlesMethod(),
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
    public suspend fun getActive(request: GetActiveRequest, headers: Metadata = Metadata()): Bundle = unaryRpc(
      channel,
      PolicyBundleServiceGrpc.getGetActiveMethod(),
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
    public suspend fun activate(request: ActivateRequest, headers: Metadata = Metadata()): Bundle = unaryRpc(
      channel,
      PolicyBundleServiceGrpc.getActivateMethod(),
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
    public suspend fun rollback(request: RollbackRequest, headers: Metadata = Metadata()): Bundle = unaryRpc(
      channel,
      PolicyBundleServiceGrpc.getRollbackMethod(),
      request,
      callOptions,
      headers
    )
  }

  /**
   * Skeletal implementation of the helixgitpx.policy.v1.PolicyBundleService service based on Kotlin coroutines.
   */
  public abstract class PolicyBundleServiceCoroutineImplBase(
    coroutineContext: CoroutineContext = EmptyCoroutineContext,
  ) : AbstractCoroutineServerImpl(coroutineContext) {
    /**
     * Returns the response to an RPC for helixgitpx.policy.v1.PolicyBundleService.ListBundles.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun listBundles(request: ListBundlesRequest): ListBundlesResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.policy.v1.PolicyBundleService.ListBundles is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.policy.v1.PolicyBundleService.GetActive.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun getActive(request: GetActiveRequest): Bundle = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.policy.v1.PolicyBundleService.GetActive is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.policy.v1.PolicyBundleService.Activate.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun activate(request: ActivateRequest): Bundle = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.policy.v1.PolicyBundleService.Activate is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.policy.v1.PolicyBundleService.Rollback.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun rollback(request: RollbackRequest): Bundle = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.policy.v1.PolicyBundleService.Rollback is unimplemented"))

    final override fun bindService(): ServerServiceDefinition = builder(getServiceDescriptor())
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = PolicyBundleServiceGrpc.getListBundlesMethod(),
      implementation = ::listBundles
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = PolicyBundleServiceGrpc.getGetActiveMethod(),
      implementation = ::getActive
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = PolicyBundleServiceGrpc.getActivateMethod(),
      implementation = ::activate
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = PolicyBundleServiceGrpc.getRollbackMethod(),
      implementation = ::rollback
    )).build()
  }
}
