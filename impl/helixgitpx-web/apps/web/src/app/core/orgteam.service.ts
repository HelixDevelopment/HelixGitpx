import { Injectable } from '@angular/core';
import { createConnectTransport } from '@connectrpc/connect-web';
import { createClient } from '@connectrpc/connect';
import { OrgService } from '../../../../../libs/proto/src/helixgitpx/org/v1/org_connect';
import { TeamService } from '../../../../../libs/proto/src/helixgitpx/team/v1/team_connect';

@Injectable({ providedIn: 'root' })
export class OrgTeamApi {
  readonly orgs = createClient(OrgService, createConnectTransport({
    baseUrl: 'https://orgteam.helix.local',
    credentials: 'include',
  }));

  readonly teams = createClient(TeamService, createConnectTransport({
    baseUrl: 'https://orgteam.helix.local',
    credentials: 'include',
  }));
}
