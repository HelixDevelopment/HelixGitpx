import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthSession } from './auth.service';

export const authGuard: CanActivateFn = async () => {
  const auth = inject(AuthSession);
  const router = inject(Router);
  const me = await auth.whoAmI();
  if (!me) {
    router.navigate(['/login']);
    return false;
  }
  return true;
};
