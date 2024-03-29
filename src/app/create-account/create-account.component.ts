import { Component } from '@angular/core';

import { FormBuilder, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { AccountService } from '../services/account.service';
import { NGXLogger } from 'ngx-logger';
import { faGithub } from '@fortawesome/free-brands-svg-icons';
import { faEnvelope, faLock, faUser } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';

@Component({
  selector: 'app-create-account',
  standalone: true,
  imports: [ReactiveFormsModule, FormsModule, RouterModule, FontAwesomeModule],
  templateUrl: './create-account.component.html',
  styleUrls: ['./create-account.component.scss']
})
export class CreateAccountComponent {

  faGithub = faGithub;
  faEnvelope = faEnvelope;
  faLock = faLock;
  faUser = faUser;

  form = this.fb.group({
    terms: [false, [Validators.requiredTrue]],
    mailingList: [false],
    username: ['', [Validators.required, Validators.minLength(4), Validators.maxLength(20)]],
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required, Validators.pattern(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,20}$/)]],
    confirmPassword: ['', [Validators.required]]
  });

  errorMessage = '';

  constructor(
    private logger: NGXLogger,
    private fb: FormBuilder,
    private accountService: AccountService,
  ) { }

  onCreateAccount() {
    this.logger.trace('creating account');
    this.form.markAllAsTouched();
    this.validateMatchingPasswords();
    if (this.form.invalid) {
      return;
    }
    this.accountService.createAccount(
      this.form.get('email')!.value!,
      this.form.get('password')!.value!,
      this.form.get('username')!.value!,
      this.form.get('mailingList')!.value!,
    ).subscribe({
      next: (res) => {
        this.logger.debug('Login response:', res);
        this.accountService.setShowCreateAccount(false);
      },
      error: (err) => {
        switch (err.code) {
          case 'auth/email-already-in-use':
            this.errorMessage = 'An account already exists with that email address.';
            break;
          case 'auth/invalid-email':
            this.errorMessage = 'Invalid email address.';
            break;
          case 'auth/operation-not-allowed':
            this.errorMessage = 'Account creation is not allowed.';
            break;
          case 'auth/weak-password':
            this.errorMessage = 'Password is too weak.';
            break;
          case 'auth/too-many-requests':
            this.errorMessage = 'Too many requests. Try again later.';
            break;
          case 'auth/network-request-failed':
            this.errorMessage = 'Network request failed. Try again later.';
            break;
          case 'auth/user-disabled':
            this.errorMessage = 'This account has been disabled.';
            break;
          default:
            this.errorMessage = 'An error occurred while creating the account.';
        }
      },
    });
  }

  onCancel() {
    this.logger.trace('cancelling');
    this.accountService.setShowCreateAccount(false);
  }

  validateMatchingPasswords() {
    this.logger.trace('validating matching passwords');
    if (this.form.get('password')?.value !== this.form.get('confirmPassword')?.value) {
      this.logger.trace('passwords do not match');
      this.form.get('confirmPassword')?.setErrors({ passwordMismatch: true });
    }
  }

  loginWithGitHub() {
    this.logger.trace('logging in with GitHub');
    if (this.form.get('terms')?.value === false) {
      this.form.get('terms')?.markAsTouched();
      return;
    }
    this.accountService.loginWithGitHub().subscribe({
      next: () => {
        this.logger.debug('Logged in with GitHub');
        this.accountService.setShowCreateAccount(false);
      },
      error: (err) => {
        this.logger.debug('Error logging in with GitHub:', err);
        switch (err.code) {
          case 'auth/account-exists-with-different-credential':
            this.errorMessage = "An account already exists with the same email address but different sign-in credentials. Sign in using a provider that you've used to sign in to this account.";
            break;
          default:
            this.errorMessage = "An error occurred while logging in with GitHub.";
        }
      }
    });
  }

  onLogin() {
    this.accountService.setShowCreateAccount(false);
    this.accountService.setShowLogin(true);
  }
}
