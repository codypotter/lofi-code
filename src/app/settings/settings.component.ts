import { Component } from '@angular/core';
import { FormBuilder, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { AccountService } from '../services/account.service';
import { CommonModule } from '@angular/common';
import { NGXLogger } from 'ngx-logger';

@Component({
  selector: 'app-settings',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, FormsModule],
  templateUrl: './settings.component.html',
  styleUrl: './settings.component.scss'
})
export class SettingsComponent {

  successMessage = '';

  errorMessage = '';

  usernameForm = this.fb.group({
    username: ['', [Validators.required, Validators.minLength(4), Validators.maxLength(20)]],
  });

  emailForm = this.fb.group({
    email: ['', [Validators.required, Validators.email]],
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
  ) {}

  validateMatchingPasswords() {
    if (this.passwordForm.get('password')?.value !== this.passwordForm.get('confirmPassword')?.value) {
      this.passwordForm.get('confirmPassword')?.setErrors({ passwordMismatch: true });
    }
  }

  onSaveUsername() {
    this.usernameForm.markAllAsTouched();
    if (this.usernameForm.invalid) {
      this.logger.warn('username form is invalid', this.usernameForm.value);
      return;
    }
    this.logger.warn('submitting username form', this.usernameForm.value);
    this.setSuccessMessage('Username updated');
  }

  onSaveEmail() {
    this.emailForm.markAllAsTouched();
    if (this.emailForm.invalid) {
      return;
    }
    this.logger.warn('submitting email form', this.emailForm.value);
    this.setSuccessMessage('Email updated');
  }

  onSavePassword() {
    this.passwordForm.markAllAsTouched();
    this.validateMatchingPasswords();
    if (this.passwordForm.invalid) {
      return;
    }
    this.logger.warn('submitting password form', this.passwordForm.value);
    this.setSuccessMessage('Password updated');
  }

  onSaveMailingList() {
    this.mailingListForm.markAllAsTouched();
    if (this.mailingListForm.invalid) {
      return;
    }
    this.setSuccessMessage('Mailing list preferences saved');
  }

  setSuccessMessage(msg: string) {
    this.errorMessage = '';
    this.successMessage = msg;
    setTimeout(() => {
      this.successMessage = '';
    }, 3000);
  }

  setErrorMessage(msg: string) {
    this.successMessage = '';
    this.errorMessage = msg;
    setTimeout(() => {
      this.errorMessage = '';
    }, 3000);
  }
}
