import { Component, EventEmitter, Output } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { AccountService } from '../services/account.service';

@Component({
  selector: 'app-create-account',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, FormsModule, RouterModule],
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

  @Output() cancel = new EventEmitter<void>();

  constructor(private fb: FormBuilder, private accountService: AccountService) { }

  onCreateAccount() {
    this.form.markAllAsTouched();
    this.validateMatchingPasswords();
    console.warn('Form errors:', {
      usernameErrors: this.form.get('username')?.errors,
      emailErrors: this.form.get('email')?.errors,
      passwordErrors: this.form.get('password')?.errors,
      confirmPasswordErrors: this.form.get('confirmPassword')?.errors
    });
    this.accountService.createAccount(this.form.get('email')!.value!, this.form.get('password')!.value!, this.form.get('username')!.value!).subscribe((res) => {
      console.log('Login response:', res);
      console.log('currentUserInfo:', this.accountService.getCurrentUserInfo());
    });
  }

  onCancel() {
    this.cancel.emit();
  }

  validateMatchingPasswords() {
    if (this.form.get('password')?.value !== this.form.get('confirmPassword')?.value) {
      this.form.get('confirmPassword')?.setErrors({ passwordMismatch: true });
    }
  }

  loginWithGitHub() {
    this.accountService.loginWithGitHub().subscribe((res) => {
      console.log('Login with GitHub response:', res);
    });
  }
}
