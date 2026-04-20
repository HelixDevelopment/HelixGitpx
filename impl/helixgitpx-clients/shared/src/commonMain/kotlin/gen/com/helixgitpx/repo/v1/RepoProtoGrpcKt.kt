package com.helixgitpx.repo.v1

import com.google.protobuf.Empty
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
import com.helixgitpx.repo.v1.RefServiceGrpc.getServiceDescriptor as refServiceGrpcGetServiceDescriptor
import com.helixgitpx.repo.v1.RepoServiceGrpc.getServiceDescriptor as repoServiceGrpcGetServiceDescriptor

/**
 * Holder for Kotlin coroutine-based client and server APIs for helixgitpx.repo.v1.RepoService.
 */
public object RepoServiceGrpcKt {
  public const val SERVICE_NAME: String = RepoServiceGrpc.SERVICE_NAME

  @JvmStatic
  public val serviceDescriptor: ServiceDescriptor
    get() = repoServiceGrpcGetServiceDescriptor()

  public val createMethod: MethodDescriptor<CreateRepoRequest, Repo>
    @JvmStatic
    get() = RepoServiceGrpc.getCreateMethod()

  public val getMethod: MethodDescriptor<GetRepoRequest, Repo>
    @JvmStatic
    get() = RepoServiceGrpc.getGetMethod()

  public val listMethod: MethodDescriptor<ListReposRequest, ListReposResponse>
    @JvmStatic
    get() = RepoServiceGrpc.getListMethod()

  public val updateMethod: MethodDescriptor<UpdateRepoRequest, Repo>
    @JvmStatic
    get() = RepoServiceGrpc.getUpdateMethod()

  public val deleteMethod: MethodDescriptor<DeleteRepoRequest, Empty>
    @JvmStatic
    get() = RepoServiceGrpc.getDeleteMethod()

  /**
   * A stub for issuing RPCs to a(n) helixgitpx.repo.v1.RepoService service as suspending coroutines.
   */
  @StubFor(RepoServiceGrpc::class)
  public class RepoServiceCoroutineStub @JvmOverloads constructor(
    channel: Channel,
    callOptions: CallOptions = DEFAULT,
  ) : AbstractCoroutineStub<RepoServiceCoroutineStub>(channel, callOptions) {
    override fun build(channel: Channel, callOptions: CallOptions): RepoServiceCoroutineStub = RepoServiceCoroutineStub(channel, callOptions)

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
    public suspend fun create(request: CreateRepoRequest, headers: Metadata = Metadata()): Repo = unaryRpc(
      channel,
      RepoServiceGrpc.getCreateMethod(),
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
    public suspend fun `get`(request: GetRepoRequest, headers: Metadata = Metadata()): Repo = unaryRpc(
      channel,
      RepoServiceGrpc.getGetMethod(),
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
    public suspend fun list(request: ListReposRequest, headers: Metadata = Metadata()): ListReposResponse = unaryRpc(
      channel,
      RepoServiceGrpc.getListMethod(),
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
    public suspend fun update(request: UpdateRepoRequest, headers: Metadata = Metadata()): Repo = unaryRpc(
      channel,
      RepoServiceGrpc.getUpdateMethod(),
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
    public suspend fun delete(request: DeleteRepoRequest, headers: Metadata = Metadata()): Empty = unaryRpc(
      channel,
      RepoServiceGrpc.getDeleteMethod(),
      request,
      callOptions,
      headers
    )
  }

  /**
   * Skeletal implementation of the helixgitpx.repo.v1.RepoService service based on Kotlin coroutines.
   */
  public abstract class RepoServiceCoroutineImplBase(
    coroutineContext: CoroutineContext = EmptyCoroutineContext,
  ) : AbstractCoroutineServerImpl(coroutineContext) {
    /**
     * Returns the response to an RPC for helixgitpx.repo.v1.RepoService.Create.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun create(request: CreateRepoRequest): Repo = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.repo.v1.RepoService.Create is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.repo.v1.RepoService.Get.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun `get`(request: GetRepoRequest): Repo = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.repo.v1.RepoService.Get is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.repo.v1.RepoService.List.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun list(request: ListReposRequest): ListReposResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.repo.v1.RepoService.List is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.repo.v1.RepoService.Update.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun update(request: UpdateRepoRequest): Repo = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.repo.v1.RepoService.Update is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.repo.v1.RepoService.Delete.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun delete(request: DeleteRepoRequest): Empty = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.repo.v1.RepoService.Delete is unimplemented"))

    final override fun bindService(): ServerServiceDefinition = builder(repoServiceGrpcGetServiceDescriptor())
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = RepoServiceGrpc.getCreateMethod(),
      implementation = ::create
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = RepoServiceGrpc.getGetMethod(),
      implementation = ::`get`
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = RepoServiceGrpc.getListMethod(),
      implementation = ::list
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = RepoServiceGrpc.getUpdateMethod(),
      implementation = ::update
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = RepoServiceGrpc.getDeleteMethod(),
      implementation = ::delete
    )).build()
  }
}

/**
 * Holder for Kotlin coroutine-based client and server APIs for helixgitpx.repo.v1.RefService.
 */
public object RefServiceGrpcKt {
  public const val SERVICE_NAME: String = RefServiceGrpc.SERVICE_NAME

  @JvmStatic
  public val serviceDescriptor: ServiceDescriptor
    get() = refServiceGrpcGetServiceDescriptor()

  public val listMethod: MethodDescriptor<ListRefsRequest, ListRefsResponse>
    @JvmStatic
    get() = RefServiceGrpc.getListMethod()

  public val protectMethod: MethodDescriptor<ProtectRefRequest, BranchProtection>
    @JvmStatic
    get() = RefServiceGrpc.getProtectMethod()

  public val unprotectMethod: MethodDescriptor<UnprotectRefRequest, Empty>
    @JvmStatic
    get() = RefServiceGrpc.getUnprotectMethod()

  /**
   * A stub for issuing RPCs to a(n) helixgitpx.repo.v1.RefService service as suspending coroutines.
   */
  @StubFor(RefServiceGrpc::class)
  public class RefServiceCoroutineStub @JvmOverloads constructor(
    channel: Channel,
    callOptions: CallOptions = DEFAULT,
  ) : AbstractCoroutineStub<RefServiceCoroutineStub>(channel, callOptions) {
    override fun build(channel: Channel, callOptions: CallOptions): RefServiceCoroutineStub = RefServiceCoroutineStub(channel, callOptions)

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
    public suspend fun list(request: ListRefsRequest, headers: Metadata = Metadata()): ListRefsResponse = unaryRpc(
      channel,
      RefServiceGrpc.getListMethod(),
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
    public suspend fun protect(request: ProtectRefRequest, headers: Metadata = Metadata()): BranchProtection = unaryRpc(
      channel,
      RefServiceGrpc.getProtectMethod(),
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
    public suspend fun unprotect(request: UnprotectRefRequest, headers: Metadata = Metadata()): Empty = unaryRpc(
      channel,
      RefServiceGrpc.getUnprotectMethod(),
      request,
      callOptions,
      headers
    )
  }

  /**
   * Skeletal implementation of the helixgitpx.repo.v1.RefService service based on Kotlin coroutines.
   */
  public abstract class RefServiceCoroutineImplBase(
    coroutineContext: CoroutineContext = EmptyCoroutineContext,
  ) : AbstractCoroutineServerImpl(coroutineContext) {
    /**
     * Returns the response to an RPC for helixgitpx.repo.v1.RefService.List.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun list(request: ListRefsRequest): ListRefsResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.repo.v1.RefService.List is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.repo.v1.RefService.Protect.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun protect(request: ProtectRefRequest): BranchProtection = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.repo.v1.RefService.Protect is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.repo.v1.RefService.Unprotect.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun unprotect(request: UnprotectRefRequest): Empty = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.repo.v1.RefService.Unprotect is unimplemented"))

    final override fun bindService(): ServerServiceDefinition = builder(refServiceGrpcGetServiceDescriptor())
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = RefServiceGrpc.getListMethod(),
      implementation = ::list
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = RefServiceGrpc.getProtectMethod(),
      implementation = ::protect
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = RefServiceGrpc.getUnprotectMethod(),
      implementation = ::unprotect
    )).build()
  }
}
