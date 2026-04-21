package com.helixgitpx.ai.v1

import com.google.protobuf.Empty
import com.helixgitpx.ai.v1.AIServiceGrpc.getServiceDescriptor
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
import io.grpc.kotlin.ClientCalls.serverStreamingRpc
import io.grpc.kotlin.ClientCalls.unaryRpc
import io.grpc.kotlin.ServerCalls.serverStreamingServerMethodDefinition
import io.grpc.kotlin.ServerCalls.unaryServerMethodDefinition
import io.grpc.kotlin.StubFor
import kotlin.String
import kotlin.coroutines.CoroutineContext
import kotlin.coroutines.EmptyCoroutineContext
import kotlin.jvm.JvmOverloads
import kotlin.jvm.JvmStatic
import kotlinx.coroutines.flow.Flow

/**
 * Holder for Kotlin coroutine-based client and server APIs for helixgitpx.ai.v1.AIService.
 */
public object AIServiceGrpcKt {
  public const val SERVICE_NAME: String = AIServiceGrpc.SERVICE_NAME

  @JvmStatic
  public val serviceDescriptor: ServiceDescriptor
    get() = getServiceDescriptor()

  public val summarizeMethod: MethodDescriptor<SummarizeRequest, SummarizeResponse>
    @JvmStatic
    get() = AIServiceGrpc.getSummarizeMethod()

  public val proposeConflictResolutionMethod:
      MethodDescriptor<ProposeConflictResolutionRequest, ConflictProposal>
    @JvmStatic
    get() = AIServiceGrpc.getProposeConflictResolutionMethod()

  public val suggestLabelMethod: MethodDescriptor<SuggestLabelRequest, SuggestLabelResponse>
    @JvmStatic
    get() = AIServiceGrpc.getSuggestLabelMethod()

  public val chatMethod: MethodDescriptor<ChatRequest, ChatMessage>
    @JvmStatic
    get() = AIServiceGrpc.getChatMethod()

  public val feedbackMethod: MethodDescriptor<FeedbackRequest, Empty>
    @JvmStatic
    get() = AIServiceGrpc.getFeedbackMethod()

  /**
   * A stub for issuing RPCs to a(n) helixgitpx.ai.v1.AIService service as suspending coroutines.
   */
  @StubFor(AIServiceGrpc::class)
  public class AIServiceCoroutineStub @JvmOverloads constructor(
    channel: Channel,
    callOptions: CallOptions = DEFAULT,
  ) : AbstractCoroutineStub<AIServiceCoroutineStub>(channel, callOptions) {
    override fun build(channel: Channel, callOptions: CallOptions): AIServiceCoroutineStub = AIServiceCoroutineStub(channel, callOptions)

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
    public suspend fun summarize(request: SummarizeRequest, headers: Metadata = Metadata()): SummarizeResponse = unaryRpc(
      channel,
      AIServiceGrpc.getSummarizeMethod(),
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
    public suspend fun proposeConflictResolution(request: ProposeConflictResolutionRequest, headers: Metadata = Metadata()): ConflictProposal = unaryRpc(
      channel,
      AIServiceGrpc.getProposeConflictResolutionMethod(),
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
    public suspend fun suggestLabel(request: SuggestLabelRequest, headers: Metadata = Metadata()): SuggestLabelResponse = unaryRpc(
      channel,
      AIServiceGrpc.getSuggestLabelMethod(),
      request,
      callOptions,
      headers
    )

    /**
     * Returns a [Flow] that, when collected, executes this RPC and emits responses from the
     * server as they arrive.  That flow finishes normally if the server closes its response with
     * [`Status.OK`][io.grpc.Status], and fails by throwing a [StatusException] otherwise.  If
     * collecting the flow downstream fails exceptionally (including via cancellation), the RPC
     * is cancelled with that exception as a cause.
     *
     * @param request The request message to send to the server.
     *
     * @param headers Metadata to attach to the request.  Most users will not need this.
     *
     * @return A flow that, when collected, emits the responses from the server.
     */
    public fun chat(request: ChatRequest, headers: Metadata = Metadata()): Flow<ChatMessage> = serverStreamingRpc(
      channel,
      AIServiceGrpc.getChatMethod(),
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
    public suspend fun feedback(request: FeedbackRequest, headers: Metadata = Metadata()): Empty = unaryRpc(
      channel,
      AIServiceGrpc.getFeedbackMethod(),
      request,
      callOptions,
      headers
    )
  }

  /**
   * Skeletal implementation of the helixgitpx.ai.v1.AIService service based on Kotlin coroutines.
   */
  public abstract class AIServiceCoroutineImplBase(
    coroutineContext: CoroutineContext = EmptyCoroutineContext,
  ) : AbstractCoroutineServerImpl(coroutineContext) {
    /**
     * Returns the response to an RPC for helixgitpx.ai.v1.AIService.Summarize.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun summarize(request: SummarizeRequest): SummarizeResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.ai.v1.AIService.Summarize is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.ai.v1.AIService.ProposeConflictResolution.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun proposeConflictResolution(request: ProposeConflictResolutionRequest): ConflictProposal = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.ai.v1.AIService.ProposeConflictResolution is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.ai.v1.AIService.SuggestLabel.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun suggestLabel(request: SuggestLabelRequest): SuggestLabelResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.ai.v1.AIService.SuggestLabel is unimplemented"))

    /**
     * Returns a [Flow] of responses to an RPC for helixgitpx.ai.v1.AIService.Chat.
     *
     * If creating or collecting the returned flow fails with a [StatusException], the RPC
     * will fail with the corresponding [io.grpc.Status].  If it fails with a
     * [java.util.concurrent.CancellationException], the RPC will fail with status `Status.CANCELLED`.  If creating
     * or collecting the returned flow fails for any other reason, the RPC will fail with
     * `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open fun chat(request: ChatRequest): Flow<ChatMessage> = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.ai.v1.AIService.Chat is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.ai.v1.AIService.Feedback.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun feedback(request: FeedbackRequest): Empty = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.ai.v1.AIService.Feedback is unimplemented"))

    final override fun bindService(): ServerServiceDefinition = builder(getServiceDescriptor())
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = AIServiceGrpc.getSummarizeMethod(),
      implementation = ::summarize
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = AIServiceGrpc.getProposeConflictResolutionMethod(),
      implementation = ::proposeConflictResolution
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = AIServiceGrpc.getSuggestLabelMethod(),
      implementation = ::suggestLabel
    ))
      .addMethod(serverStreamingServerMethodDefinition(
      context = this.context,
      descriptor = AIServiceGrpc.getChatMethod(),
      implementation = ::chat
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = AIServiceGrpc.getFeedbackMethod(),
      implementation = ::feedback
    )).build()
  }
}
