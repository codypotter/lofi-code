import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AccountService } from '../services/account.service';
import { map, tap } from 'rxjs';
import { NGXLogger } from 'ngx-logger';

export const authGuard: CanActivateFn = (route, state) => {
  const acctService = inject(AccountService);
  const router = inject(Router);
  const logger = inject(NGXLogger);
  logger.trace('authGuard: checking if user is logged in')
  return acctService.currentUser.pipe(
    map(user => !!user),
    tap(loggedIn => {
      logger.trace('authGuard: user is logged in:', loggedIn);
      if (!loggedIn) {
        logger.trace('authGuard: user is not logged in, redirecting to /');
        router.navigate(['/']);
      }
    })
  );
};
