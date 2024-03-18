import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CreateAccountComponent } from '../create-account/create-account.component';
import { User } from 'firebase/auth';
import { AccountService } from '../services/account.service';
import { RouterLink } from '@angular/router';
import { LoginComponent } from '../login/login.component';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [CommonModule, CreateAccountComponent, RouterLink, LoginComponent],
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.scss']
})
export class NavbarComponent implements OnInit {
  user: User | null = null;

  showCreateAccount = false;

  showLogin = false;

  isNavbarBurgerActive = false;

  constructor(private accountService: AccountService) {}

  ngOnInit(): void {
    this.accountService.currentUser.subscribe((user) => {
      this.user = user;
    })
  }

  onLogOut() {
    this.accountService.logout()
  }
}
