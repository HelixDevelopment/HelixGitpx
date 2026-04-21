package com.helixgitpx.search.v1

import com.helixgitpx.search.v1.SearchServiceGrpc.getServiceDescriptor
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
 * Holder for Kotlin coroutine-based client and server APIs for helixgitpx.search.v1.SearchService.
 */
public object SearchServiceGrpcKt {
  public const val SERVICE_NAME: String = SearchServiceGrpc.SERVICE_NAME

  @JvmStatic
  public val serviceDescriptor: ServiceDescriptor
    get() = getServiceDescriptor()

  public val searchMethod: MethodDescriptor<SearchRequest, SearchResponse>
    @JvmStatic
    get() = SearchServiceGrpc.getSearchMethod()

  public val codeSearchMethod: MethodDescriptor<CodeSearchRequest, CodeSearchResponse>
    @JvmStatic
    get() = SearchServiceGrpc.getCodeSearchMethod()

  public val indexStatusMethod: MethodDescriptor<IndexStatusRequest, IndexStatusResponse>
    @JvmStatic
    get() = SearchServiceGrpc.getIndexStatusMethod()

  /**
   * A stub for issuing RPCs to a(n) helixgitpx.search.v1.SearchService service as suspending coroutines.
   */
  @StubFor(SearchServiceGrpc::class)
  public class SearchServiceCoroutineStub @JvmOverloads constructor(
    channel: Channel,
    callOptions: CallOptions = DEFAULT,
  ) : AbstractCoroutineStub<SearchServiceCoroutineStub>(channel, callOptions) {
    override fun build(channel: Channel, callOptions: CallOptions): SearchServiceCoroutineStub = SearchServiceCoroutineStub(channel, callOptions)

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
    public suspend fun search(request: SearchRequest, headers: Metadata = Metadata()): SearchResponse = unaryRpc(
      channel,
      SearchServiceGrpc.getSearchMethod(),
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
    public suspend fun codeSearch(request: CodeSearchRequest, headers: Metadata = Metadata()): CodeSearchResponse = unaryRpc(
      channel,
      SearchServiceGrpc.getCodeSearchMethod(),
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
    public suspend fun indexStatus(request: IndexStatusRequest, headers: Metadata = Metadata()): IndexStatusResponse = unaryRpc(
      channel,
      SearchServiceGrpc.getIndexStatusMethod(),
      request,
      callOptions,
      headers
    )
  }

  /**
   * Skeletal implementation of the helixgitpx.search.v1.SearchService service based on Kotlin coroutines.
   */
  public abstract class SearchServiceCoroutineImplBase(
    coroutineContext: CoroutineContext = EmptyCoroutineContext,
  ) : AbstractCoroutineServerImpl(coroutineContext) {
    /**
     * Returns the response to an RPC for helixgitpx.search.v1.SearchService.Search.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun search(request: SearchRequest): SearchResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.search.v1.SearchService.Search is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.search.v1.SearchService.CodeSearch.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun codeSearch(request: CodeSearchRequest): CodeSearchResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.search.v1.SearchService.CodeSearch is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.search.v1.SearchService.IndexStatus.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun indexStatus(request: IndexStatusRequest): IndexStatusResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.search.v1.SearchService.IndexStatus is unimplemented"))

    final override fun bindService(): ServerServiceDefinition = builder(getServiceDescriptor())
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = SearchServiceGrpc.getSearchMethod(),
      implementation = ::search
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = SearchServiceGrpc.getCodeSearchMethod(),
      implementation = ::codeSearch
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = SearchServiceGrpc.getIndexStatusMethod(),
      implementation = ::indexStatus
    )).build()
  }
}
