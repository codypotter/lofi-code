import { Component } from '@angular/core';
import { FormBuilder, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { AccountService } from '../services/account.service';
import { CommonModule } from '@angular/common';
import { NGXLogger } from 'ngx-logger';
import { Router, RouterModule, RouterState } from '@angular/router';

@Component({
  selector: 'app-settings',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, FormsModule, RouterModule],
  templateUrl: './settings.component.html',
  styleUrl: './settings.component.scss'
})
export class SettingsComponent {

  successMessage = '';

  errorMessage = '';

  usernameForm = this.fb.group({
    username: [this.accountService.getCurrentUserInfo()?.displayName ?? '', [Validators.required, Validators.minLength(4), Validators.maxLength(20)]],
  });

  emailForm = this.fb.group({
    email: [this.accountService.getCurrentUserInfo()?.email ?? '', [Validators.required, Validators.email]],
  });

  passwordForm = this.fb.group({
    password: ['', [Validators.required, Validators.pattern(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,20}$/)]],
    confirmPassword: ['', [Validators.required]],
  });

  mailingListForm = this.fb.group({
    mailingList: [false]
  });

  constructor(
    private logger: NGXLogger,
    private fb: FormBuilder,
    private accountService: AccountService,
    private router: Router,
  ) { }

  validateMatchingPasswords() {
    this.logger.trace('validating matching passwords');
    if (this.passwordForm.get('password')?.value !== this.passwordForm.get('confirmPassword')?.value) {
      this.logger.trace('passwords do not match');
      this.passwordForm.get('confirmPassword')?.setErrors({ passwordMismatch: true });
    }
  }

  onSaveUsername() {
    this.logger.trace('saving username');
    this.usernameForm.markAllAsTouched();
    if (this.usernameForm.invalid) {
      this.logger.warn('username form is invalid', this.usernameForm.value);
      return;
    }
    this.logger.trace('submitting username form', this.usernameForm.value);
    this.accountService.updateUsername(this.usernameForm.get('username')!.value!).subscribe({
      next: () => this.setSuccessMessage('Username updated'),
      error: (err) => this.setErrorMessage(err.message),
    });
  }

  onSaveEmail() {
    this.logger.trace('saving email');
    this.emailForm.markAllAsTouched();
    if (this.emailForm.invalid) {
      return;
    }

    this.accountService.updateEmail(this.emailForm.get('email')!.value!).subscribe({
      next: () => this.setSuccessMessage('An email has been sent to verify the new email.'),
      error: (err) => {
        switch (err.code) {
          case 'auth/invalid-email':
            this.setErrorMessage('Invalid email');
            break;
          case 'auth/email-already-in-use':
            this.setErrorMessage('Email already in use');
            break;
          case 'auth/requires-recent-login':
            this.setErrorMessage('Please log in again to change email. Logging out...');
            setTimeout(() => this.accountService.logout(), 3000);
            break;
          default:
            this.setErrorMessage(err.message);
        }
        if (err.code === 'auth/operation-not-allowed') {
          this.setErrorMessage('Please verify the new email before changing email.');
        } else {
          this.setErrorMessage(err.message);
        }
      },
    });
  }

  onSavePassword() {
    this.logger.trace('saving password');
    this.passwordForm.markAllAsTouched();
    this.validateMatchingPasswords();
    if (this.passwordForm.invalid) {
      return;
    }
    this.accountService.updatePassword(this.passwordForm.get('password')!.value!).subscribe({
      next: () => this.setSuccessMessage('Password updated'),
      error: (err) => {
        switch (err.code) {
          case 'auth/weak-password':
            this.setErrorMessage('Password is too weak');
            break;
          case 'auth/requires-recent-login':
            this.setErrorMessage('Please log in again to change password. Logging out...');
            setTimeout(() => {
              this.accountService.logout();
              this.router.navigate(['/'])
            }, 3000);
            break;
          default:
            this.setErrorMessage(err.message);
        }
      },
    });
  }

  onSaveMailingList() {
    this.logger.trace('saving mailing list preferences');
    this.mailingListForm.markAllAsTouched();
    if (this.mailingListForm.invalid) {
      return;
    }
    this.setSuccessMessage('Mailing list preferences saved');
  }

  setSuccessMessage(msg: string) {
    this.logger.trace('setting success message');
    this.errorMessage = '';
    this.successMessage = msg;
    setTimeout(() => {
      this.successMessage = '';
      this.logger.trace('success message cleared');
    }, 3000);
  }

  setErrorMessage(msg: string) {
    this.logger.trace('setting error message');
    this.successMessage = '';
    this.errorMessage = msg;
    setTimeout(() => {
      this.errorMessage = '';
      this.logger.trace('error message cleared');
    }, 3000);
  }
}
