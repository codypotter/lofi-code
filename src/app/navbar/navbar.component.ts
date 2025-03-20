import { Component, OnDestroy, OnInit } from '@angular/core';

import { CreateAccountComponent } from '../create-account/create-account.component';
import { AccountService } from '../services/account.service';
import { RouterModule } from '@angular/router';
import { LoginComponent } from '../login/login.component';
import { Subject, takeUntil } from 'rxjs';
import { User } from '@angular/fire/auth';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [CreateAccountComponent, RouterModule, LoginComponent],
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.scss']
})
export class NavbarComponent implements OnInit, OnDestroy {
  user: User | null = null;

  showCreateAccount = false;

  showLogin = false;

  isNavbarBurgerActive = false;

  destroy$ = new Subject<void>();

  constructor(private accountService: AccountService) {}

  ngOnInit(): void {
    this.accountService.currentUser.pipe(takeUntil(this.destroy$)).subscribe((user) => {
      this.user = user;
    });

    this.accountService.showLogin.pipe(takeUntil(this.destroy$)).subscribe((showLogin) => {
      this.showLogin = showLogin;
    });

    this.accountService.showCreateAccount.pipe(takeUntil(this.destroy$)).subscribe((showCreateAccount) => {
      this.showCreateAccount = showCreateAccount;
    });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  onLogOut() {
    this.accountService.logout()
  }

  toggleLogin() {
    this.accountService.setShowLogin(!this.showLogin);
  }

  toggleCreateAccount() {
    this.accountService.setShowCreateAccount(!this.showCreateAccount);
  }
}
