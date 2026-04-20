package com.helixgitpx.team.v1

import com.google.protobuf.Empty
import com.helixgitpx.team.v1.TeamServiceGrpc.getServiceDescriptor
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
 * Holder for Kotlin coroutine-based client and server APIs for helixgitpx.team.v1.TeamService.
 */
public object TeamServiceGrpcKt {
  public const val SERVICE_NAME: String = TeamServiceGrpc.SERVICE_NAME

  @JvmStatic
  public val serviceDescriptor: ServiceDescriptor
    get() = getServiceDescriptor()

  public val createMethod: MethodDescriptor<CreateTeamRequest, Team>
    @JvmStatic
    get() = TeamServiceGrpc.getCreateMethod()

  public val getMethod: MethodDescriptor<GetTeamRequest, Team>
    @JvmStatic
    get() = TeamServiceGrpc.getGetMethod()

  public val listMethod: MethodDescriptor<ListTeamsRequest, ListTeamsResponse>
    @JvmStatic
    get() = TeamServiceGrpc.getListMethod()

  public val updateMethod: MethodDescriptor<UpdateTeamRequest, Team>
    @JvmStatic
    get() = TeamServiceGrpc.getUpdateMethod()

  public val deleteMethod: MethodDescriptor<DeleteTeamRequest, Empty>
    @JvmStatic
    get() = TeamServiceGrpc.getDeleteMethod()

  public val addMemberMethod: MethodDescriptor<AddMemberRequest, Membership>
    @JvmStatic
    get() = TeamServiceGrpc.getAddMemberMethod()

  public val removeMemberMethod: MethodDescriptor<RemoveMemberRequest, Empty>
    @JvmStatic
    get() = TeamServiceGrpc.getRemoveMemberMethod()

  public val updateMemberRoleMethod: MethodDescriptor<UpdateMemberRoleRequest, Membership>
    @JvmStatic
    get() = TeamServiceGrpc.getUpdateMemberRoleMethod()

  public val listMembersMethod: MethodDescriptor<ListMembersRequest, ListMembersResponse>
    @JvmStatic
    get() = TeamServiceGrpc.getListMembersMethod()

  /**
   * A stub for issuing RPCs to a(n) helixgitpx.team.v1.TeamService service as suspending coroutines.
   */
  @StubFor(TeamServiceGrpc::class)
  public class TeamServiceCoroutineStub @JvmOverloads constructor(
    channel: Channel,
    callOptions: CallOptions = DEFAULT,
  ) : AbstractCoroutineStub<TeamServiceCoroutineStub>(channel, callOptions) {
    override fun build(channel: Channel, callOptions: CallOptions): TeamServiceCoroutineStub = TeamServiceCoroutineStub(channel, callOptions)

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
    public suspend fun create(request: CreateTeamRequest, headers: Metadata = Metadata()): Team = unaryRpc(
      channel,
      TeamServiceGrpc.getCreateMethod(),
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
    public suspend fun `get`(request: GetTeamRequest, headers: Metadata = Metadata()): Team = unaryRpc(
      channel,
      TeamServiceGrpc.getGetMethod(),
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
    public suspend fun list(request: ListTeamsRequest, headers: Metadata = Metadata()): ListTeamsResponse = unaryRpc(
      channel,
      TeamServiceGrpc.getListMethod(),
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
    public suspend fun update(request: UpdateTeamRequest, headers: Metadata = Metadata()): Team = unaryRpc(
      channel,
      TeamServiceGrpc.getUpdateMethod(),
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
    public suspend fun delete(request: DeleteTeamRequest, headers: Metadata = Metadata()): Empty = unaryRpc(
      channel,
      TeamServiceGrpc.getDeleteMethod(),
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
    public suspend fun addMember(request: AddMemberRequest, headers: Metadata = Metadata()): Membership = unaryRpc(
      channel,
      TeamServiceGrpc.getAddMemberMethod(),
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
    public suspend fun removeMember(request: RemoveMemberRequest, headers: Metadata = Metadata()): Empty = unaryRpc(
      channel,
      TeamServiceGrpc.getRemoveMemberMethod(),
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
    public suspend fun updateMemberRole(request: UpdateMemberRoleRequest, headers: Metadata = Metadata()): Membership = unaryRpc(
      channel,
      TeamServiceGrpc.getUpdateMemberRoleMethod(),
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
    public suspend fun listMembers(request: ListMembersRequest, headers: Metadata = Metadata()): ListMembersResponse = unaryRpc(
      channel,
      TeamServiceGrpc.getListMembersMethod(),
      request,
      callOptions,
      headers
    )
  }

  /**
   * Skeletal implementation of the helixgitpx.team.v1.TeamService service based on Kotlin coroutines.
   */
  public abstract class TeamServiceCoroutineImplBase(
    coroutineContext: CoroutineContext = EmptyCoroutineContext,
  ) : AbstractCoroutineServerImpl(coroutineContext) {
    /**
     * Returns the response to an RPC for helixgitpx.team.v1.TeamService.Create.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun create(request: CreateTeamRequest): Team = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.team.v1.TeamService.Create is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.team.v1.TeamService.Get.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun `get`(request: GetTeamRequest): Team = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.team.v1.TeamService.Get is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.team.v1.TeamService.List.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun list(request: ListTeamsRequest): ListTeamsResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.team.v1.TeamService.List is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.team.v1.TeamService.Update.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun update(request: UpdateTeamRequest): Team = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.team.v1.TeamService.Update is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.team.v1.TeamService.Delete.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun delete(request: DeleteTeamRequest): Empty = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.team.v1.TeamService.Delete is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.team.v1.TeamService.AddMember.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun addMember(request: AddMemberRequest): Membership = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.team.v1.TeamService.AddMember is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.team.v1.TeamService.RemoveMember.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun removeMember(request: RemoveMemberRequest): Empty = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.team.v1.TeamService.RemoveMember is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.team.v1.TeamService.UpdateMemberRole.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun updateMemberRole(request: UpdateMemberRoleRequest): Membership = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.team.v1.TeamService.UpdateMemberRole is unimplemented"))

    /**
     * Returns the response to an RPC for helixgitpx.team.v1.TeamService.ListMembers.
     *
     * If this method fails with a [StatusException], the RPC will fail with the corresponding
     * [io.grpc.Status].  If this method fails with a [java.util.concurrent.CancellationException], the RPC will fail
     * with status `Status.CANCELLED`.  If this method fails for any other reason, the RPC will
     * fail with `Status.UNKNOWN` with the exception as a cause.
     *
     * @param request The request from the client.
     */
    public open suspend fun listMembers(request: ListMembersRequest): ListMembersResponse = throw StatusException(UNIMPLEMENTED.withDescription("Method helixgitpx.team.v1.TeamService.ListMembers is unimplemented"))

    final override fun bindService(): ServerServiceDefinition = builder(getServiceDescriptor())
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = TeamServiceGrpc.getCreateMethod(),
      implementation = ::create
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = TeamServiceGrpc.getGetMethod(),
      implementation = ::`get`
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = TeamServiceGrpc.getListMethod(),
      implementation = ::list
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = TeamServiceGrpc.getUpdateMethod(),
      implementation = ::update
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = TeamServiceGrpc.getDeleteMethod(),
      implementation = ::delete
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = TeamServiceGrpc.getAddMemberMethod(),
      implementation = ::addMember
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = TeamServiceGrpc.getRemoveMemberMethod(),
      implementation = ::removeMember
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = TeamServiceGrpc.getUpdateMemberRoleMethod(),
      implementation = ::updateMemberRole
    ))
      .addMethod(unaryServerMethodDefinition(
      context = this.context,
      descriptor = TeamServiceGrpc.getListMembersMethod(),
      implementation = ::listMembers
    )).build()
  }
}
