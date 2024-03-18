import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AccountService } from '../services/account.service';
import { map, tap } from 'rxjs';

export const authGuard: CanActivateFn = () => {
  const acctService = inject(AccountService);
  const router = inject(Router);
  return acctService.currentUser.pipe(
    map(user => !!user),
    tap(loggedIn => {
      if (!loggedIn) {
        router.navigate(['/']);
      }
    })
  );
};
