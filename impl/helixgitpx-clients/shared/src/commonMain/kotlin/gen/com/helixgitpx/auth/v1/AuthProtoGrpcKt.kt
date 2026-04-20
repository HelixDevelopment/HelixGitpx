package com.helixgitpx.auth.v1

import com.google.protobuf.Empty
import com.helixgitpx.auth.v1.AuthServiceGrpc.getServiceDescriptor
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
 * Holder for Kotlin coroutine-based client and server APIs for helixgitpx.auth.v1.AuthService.
 */
public object AuthServiceGrpcKt {
  public const val SERVICE_NAME: String = AuthServiceGrpc.SERVICE_NAME

  @JvmStatic
  public val serviceDescriptor: ServiceDescriptor
    get() = getServiceDescriptor()

  public val exchangeOIDCMethod: MethodDescriptor<ExchangeOIDCRequest, Tokens>
    @JvmStatic
    get() = AuthServiceGrpc.getExchangeOIDCMethod()

  public val refreshTokenMethod: MethodDescriptor<RefreshTokenRequest, Tokens>
    @JvmStatic
    get() = AuthServiceGrpc.getRefreshTokenMethod()

  public val issuePATMethod: MethodDescriptor<IssuePATRequest, PAT>
    @JvmStatic
    get() = AuthServiceGrpc.getIssuePATMethod()

  public val revokePATMethod: MethodDescriptor<RevokePATRequest, Empty>
    @JvmStatic
    get() = AuthServiceGrpc.getRevokePATMethod()

  public val listPATsMethod: MethodDescriptor<Empty, ListPATsResponse>
    @JvmStatic
    get() = AuthServiceGrpc.getListPATsMethod()

  public val whoAmIMethod: MethodDescriptor<Empty, User>
    @JvmStatic
    get() = AuthServiceGrpc.getWhoAmIMethod()

  public val enrollTOTPMethod: MethodDescriptor<Empty, EnrollTOTPResponse>
    @JvmStatic
    get() = AuthServiceGrpc.getEnrollTOTPMethod()

  public val enrollFIDO2Method: MethodDescriptor<EnrollFIDO2Request, EnrollFIDO2Response>
    @JvmStatic
    get() = AuthServiceGrpc.getEnrollFIDO2Method()

  public val verifyMFAMethod: MethodDescriptor<VerifyMFARequest, MFAVerification>
    @JvmStatic
    get() = AuthServiceGrpc.getVerifyMFAMethod()

  public val listSessionsMethod: MethodDescriptor<Empty, ListSessionsResponse>
    @JvmStatic
    get() = AuthServiceGrpc.getListSessionsMethod()

  public val revokeSessionMethod: MethodDescriptor<RevokeSessionRequest, Empty>
    @JvmStatic
    get() = AuthServiceGrpc.getRevokeSessionMethod()

  /**
   * A stub for issuing RPCs to a(n) helixgitpx.auth.v1.AuthService service as suspending coroutines.
   */
  @StubFor(AuthServiceGrpc::class)
  public class AuthServiceCoroutineStub @JvmOverloads constructor(
    channel: Channel,
    callOptions: CallOptions = DEFAULT,
  ) : AbstractCoroutineStub<AuthServiceCoroutineStub>(channel, callOptions) {
    override fun build(channel: Channel, callOptions: CallOptions): AuthServiceCoroutineStub = AuthServiceCoroutineStub(channel, callOptions)

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
    public suspend fun exchangeOIDC(request: ExchangeOIDCRequest, headers: Metadata = Metadata()): Tokens = unaryRpc(
      channel,
      AuthServiceGrpc.getExchangeOIDCMethod(),
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
    public suspend fun refreshToken(request: RefreshTokenRequest, headers: Metadata = Metadata()): Tokens = unaryRpc(
      channel,
      AuthServiceGrpc.getRefreshTokenMethod(),
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
    public suspend fun issuePAT(request: IssuePATRequest, headers: Metadata = Metadata()): PAT = unaryRpc(
      channel,
      AuthServiceGrpc.getIssuePATMethod(),
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
    public suspend fun revokePAT(request: RevokePATRequest, headers: Metadata = Metadata()): Empty = unaryRpc(
      channel,
      AuthServiceGrpc.getRevokePATMethod(),
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
    public suspend fun listPATs(request: Empty, headers: Metadata = Metadata()): ListPATsResponse = unaryRpc(
      channel,
      AuthServiceGrpc.getListPATsMethod(),
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
    public suspend fun whoAmI(request: Empty, headers: Metadata = Metadata()): User = unaryRpc(
      channel,
      AuthServiceGrpc.getWhoAmIMethod(),
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
    public suspend fun enrollTOTP(request: Empty, headers: Metadata = Metadata()): EnrollTOTPResponse = unaryRpc(
      channel,
      AuthServiceGrpc.getEnrollTOTPMethod(),
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
    public suspend fun enrollFIDO2(request: EnrollFIDO2Request, headers: Metadata = Metadata()): EnrollFIDO2Response = unaryRpc(
      channel,
      AuthServiceGrpc.getEnrollFIDO2Method(),
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
    public suspend fun verifyMFA(request: VerifyMFARequest, headers: Metadata = Metadata()): MFAVerification = unaryRpc(
      channel,
      AuthServiceGrpc.getVerifyMFAMethod(),
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
    public suspend fun listSessions(request: Empty, headers: Metadata = Metadata()): ListSessionsResponse = unaryRpc(
      channel,
      AuthServiceGrpc.getListSessionsMethod(),
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
    public suspend fun revokeSession(request: RevokeSessionRequest, headers: Metadata = Metadata()): Empty = unaryRpc(
      channel,
      AuthServiceGrpc.getRevokeSessionMethod(),
      request,
      callOptions,
      headers
    )
  }

  /**
   * Skeletal implementation of the helixgitpx.auth.v1.AuthService service based on Kotlin coroutines.
   */
  public abstract class AuthServiceCoroutineImplBase(
    coroutineContext: CoroutineContext = EmptyCoroutineContext,
  ) : AbstractCoroutineServerImpl(coroutineContext) {
    /**
     * Returns the response to an RPC for helixgitpx.auth.v1.AuthService.ExchangeOIDC.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun exchangeOIDC(request: ExchangeOIDCRequest): Tokens = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.auth.v1.AuthService.ExchangeOIDC is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.auth.v1.AuthService.RefreshToken.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun refreshToken(request: RefreshTokenRequest): Tokens = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.auth.v1.AuthService.RefreshToken is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.auth.v1.AuthService.IssuePAT.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun issuePAT(request: IssuePATRequest): PAT = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.auth.v1.AuthService.IssuePAT is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.auth.v1.AuthService.RevokePAT.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun revokePAT(request: RevokePATRequest): Empty = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.auth.v1.AuthService.RevokePAT is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.auth.v1.AuthService.ListPATs.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun listPATs(request: Empty): ListPATsResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.auth.v1.AuthService.ListPATs is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.auth.v1.AuthService.WhoAmI.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun whoAmI(request: Empty): User = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.auth.v1.AuthService.WhoAmI is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.auth.v1.AuthService.EnrollTOTP.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun enrollTOTP(request: Empty): EnrollTOTPResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.auth.v1.AuthService.EnrollTOTP is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.auth.v1.AuthService.EnrollFIDO2.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun enrollFIDO2(request: EnrollFIDO2Request): EnrollFIDO2Response = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.auth.v1.AuthService.EnrollFIDO2 is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.auth.v1.AuthService.VerifyMFA.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun verifyMFA(request: VerifyMFARequest): MFAVerification = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.auth.v1.AuthService.VerifyMFA is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.auth.v1.AuthService.ListSessions.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun listSessions(request: Empty): ListSessionsResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.auth.v1.AuthService.ListSessions is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.auth.v1.AuthService.RevokeSession.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun revokeSession(request: RevokeSessionRequest): Empty = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.auth.v1.AuthService.RevokeSession is unimplemented"))

    final override fun bindService(): ServerServiceDefinition = builder(getServiceDescriptor())
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = AuthServiceGrpc.getExchangeOIDCMethod(),
      implementation = ::exchangeOIDC
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = AuthServiceGrpc.getRefreshTokenMethod(),
      implementation = ::refreshToken
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = AuthServiceGrpc.getIssuePATMethod(),
      implementation = ::issuePAT
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = AuthServiceGrpc.getRevokePATMethod(),
      implementation = ::revokePAT
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = AuthServiceGrpc.getListPATsMethod(),
      implementation = ::listPATs
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = AuthServiceGrpc.getWhoAmIMethod(),
      implementation = ::whoAmI
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = AuthServiceGrpc.getEnrollTOTPMethod(),
      implementation = ::enrollTOTP
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = AuthServiceGrpc.getEnrollFIDO2Method(),
      implementation = ::enrollFIDO2
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = AuthServiceGrpc.getVerifyMFAMethod(),
      implementation = ::verifyMFA
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = AuthServiceGrpc.getListSessionsMethod(),
      implementation = ::listSessions
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = AuthServiceGrpc.getRevokeSessionMethod(),
      implementation = ::revokeSession
    )).build()
  }
}
