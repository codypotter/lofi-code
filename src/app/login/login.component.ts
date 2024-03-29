import { Component, EventEmitter, OnInit, Output } from '@angular/core';

import { FormBuilder, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { AccountService } from '../services/account.service';
import { NGXLogger } from 'ngx-logger';
import { CommonModule } from '@angular/common';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { faGithub } from '@fortawesome/free-brands-svg-icons';
import { faEnvelope, faLock } from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [ReactiveFormsModule, FormsModule, CommonModule, FontAwesomeModule],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent {

  faGithub = faGithub;
  faEnvelope = faEnvelope;
  faLock = faLock;
  
  form = this.fb.group({
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required, Validators.pattern(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,20}$/)]],
  });

  errorMessage = '';

  constructor(
    private logger: NGXLogger,
    private fb: FormBuilder,
    private accountService: AccountService,
  ) { }

  onLogIn() {
    this.form.markAllAsTouched();
    if (this.form.invalid) {
      return;
    }
    this.accountService.login(this.form.get('email')!.value!, this.form.get('password')!.value!).subscribe({
      next: (res) => {
        this.logger.debug('Login response:', res);
        this.logger.debug('currentUserInfo:', this.accountService.getCurrentUserInfo());
        this.accountService.setShowLogin(false);
      },
      error: (err) => {
        switch(err.code) {
          case 'auth/user-not-found':
            this.setErrorMessage('No user found with that email address.');
            break;
          case 'auth/wrong-password':
            this.setErrorMessage('Incorrect password.');
            break;
          case 'auth/too-many-requests':
            this.setErrorMessage('Too many requests. Try again later.');
            break;
          case 'auth/user-disabled':
            this.setErrorMessage('This account has been disabled.');
            break;
          case 'auth/invalid-email':
            this.setErrorMessage('Invalid email address.');
            break;
          case 'auth/invalid-credential':
            this.setErrorMessage('Invalid credential.');
            break;
          default:
            this.logger.debug('Error logging in:', err);
            this.setErrorMessage('An error occurred while logging in.');
        }
      },
    });
  }

  loginWithGitHub() {
    this.accountService.loginWithGitHub().subscribe({
      next: () => {
        this.logger.debug('currentUserInfo:', this.accountService.getCurrentUserInfo());
        this.accountService.setShowLogin(false);
      },
      error: (err) => {
        this.logger.error('Error logging in with GitHub:', err);
        switch(err.code) {
          case 'auth/account-exists-with-different-credential':
            this.setErrorMessage('An account already exists with the same email address but different sign-in credentials. Sign in using a provider associated with this email address.');
            break;
          default:
            this.setErrorMessage('An error occurred while logging in with GitHub.');
        }
      }
    });
  }

  onCancel() {
    this.accountService.setShowLogin(false);
  }

  setErrorMessage(message: string) {
    this.logger.trace('setting error message', message);
    this.errorMessage = message;
    setTimeout(() => this.errorMessage = '', 5000);
  }

  onCreateAccount() {
    this.accountService.setShowLogin(false);
    this.accountService.setShowCreateAccount(true);
  }
}
