import { Routes } from '@angular/router';
import { LoginComponent } from './login/login.component';
import { AuthCallbackComponent } from './auth-callback/auth-callback.component';
import { OrgsComponent } from './orgs/orgs.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { ReposComponent } from './repos/repos.component';
import { PrsComponent } from './prs/prs.component';
import { IssuesComponent } from './issues/issues.component';
import { ConflictsComponent } from './conflicts/conflicts.component';
import { SettingsComponent } from './settings/settings.component';
import { SearchComponent } from './search/search.component';
import { authGuard } from './core/auth.guard';

export const routes: Routes = [
  { path: '', redirectTo: '/dashboard', pathMatch: 'full' },
  { path: 'login', component: LoginComponent },
  { path: 'auth/callback', component: AuthCallbackComponent },
  { path: 'dashboard', component: DashboardComponent, canActivate: [authGuard] },
  { path: 'orgs', component: OrgsComponent, canActivate: [authGuard] },
  { path: 'repos', component: ReposComponent, canActivate: [authGuard] },
  { path: 'prs', component: PrsComponent, canActivate: [authGuard] },
  { path: 'issues', component: IssuesComponent, canActivate: [authGuard] },
  { path: 'conflicts', component: ConflictsComponent, canActivate: [authGuard] },
  { path: 'search', component: SearchComponent, canActivate: [authGuard] },
  { path: 'settings', component: SettingsComponent, canActivate: [authGuard] },
];
