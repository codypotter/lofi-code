import { Component, EventEmitter, Output } from '@angular/core';

import { FormBuilder, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { AccountService } from '../services/account.service';
import { NGXLogger } from 'ngx-logger';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [ReactiveFormsModule, FormsModule],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent {
  form = this.fb.group({
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required, Validators.pattern(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,20}$/)]],
  });

  errorMessage = '';

  @Output() done = new EventEmitter<void>();

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
    this.accountService.login(this.form.get('email')!.value!, this.form.get('password')!.value!).subscribe((res) => {
      this.logger.debug('Login response:', res);
      this.logger.debug('currentUserInfo:', this.accountService.getCurrentUserInfo());
      this.done.emit();
    });
  }

  loginWithGitHub() {
    this.accountService.loginWithGitHub().subscribe({
      next: () => {
        this.logger.debug('currentUserInfo:', this.accountService.getCurrentUserInfo());
        this.done.emit();
      },
      error: (err) => {
        this.logger.error('Error logging in with GitHub:', err);
        if (err.code === 'auth/account-exists-with-different-credential') {
          this.errorMessage = 'An account already exists with the same email address but different sign-in credentials. Sign in using a provider associated with this email address.';
        }
      }
    });
  }

  onCancel() {
    this.done.emit();
  }
}
