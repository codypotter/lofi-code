import { Component, EventEmitter, Output } from '@angular/core';

import { FormBuilder, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { AccountService } from '../services/account.service';

@Component({
  selector: 'app-create-account',
  standalone: true,
  imports: [ReactiveFormsModule, FormsModule, RouterModule],
  templateUrl: './create-account.component.html',
  styleUrls: ['./create-account.component.scss']
})
export class CreateAccountComponent {
  form = this.fb.group({
    terms: [false, [Validators.requiredTrue]],
    mailingList: [false],
    username: ['', [Validators.required, Validators.minLength(4), Validators.maxLength(20)]],
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required, Validators.pattern(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,20}$/)]],
    confirmPassword: ['', [Validators.required]]
  });

  errorMessage = '';

  @Output() done = new EventEmitter<void>();

  constructor(private fb: FormBuilder, private accountService: AccountService) { }

  onCreateAccount() {
    this.form.markAllAsTouched();
    this.validateMatchingPasswords();
    if (this.form.invalid) {
      return;
    }
    this.accountService.createAccount(this.form.get('email')!.value!, this.form.get('password')!.value!, this.form.get('username')!.value!).subscribe((res) => {
      console.log('Login response:', res);
      console.log('currentUserInfo:', this.accountService.getCurrentUserInfo());
      this.done.emit();
    });
  }

  onCancel() {
    this.done.emit();
  }

  validateMatchingPasswords() {
    if (this.form.get('password')?.value !== this.form.get('confirmPassword')?.value) {
      this.form.get('confirmPassword')?.setErrors({ passwordMismatch: true });
    }
  }

  loginWithGitHub() {
    if (this.form.get('terms')?.value === false) {
      this.form.get('terms')?.markAsTouched();
      return;
    }
    this.accountService.loginWithGitHub().subscribe({
      next: () => {
        console.log('currentUserInfo:', this.accountService.getCurrentUserInfo());
        this.done.emit();
      },
      error: (err) => {
        console.error('Error logging in with GitHub:', err);
        if (err.code === 'auth/account-exists-with-different-credential') {
          this.errorMessage = "An account already exists with the same email address but different sign-in credentials. Sign in using a provider that you've used to sign in to this account.";
        }
      }
    });
  }
}
