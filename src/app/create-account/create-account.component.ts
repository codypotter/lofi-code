import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { RouterModule } from '@angular/router';

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

  constructor(private fb: FormBuilder) { }

  onCreateAccount() {
    this.form.markAllAsTouched();
    this.validateMatchingPasswords();
    console.warn('Form errors:', {
      usernameErrors: this.form.get('username')?.errors,
      emailErrors: this.form.get('email')?.errors,
      passwordErrors: this.form.get('password')?.errors,
      confirmPasswordErrors: this.form.get('confirmPassword')?.errors
    });
    console.warn('isValid', this.form.valid)
    console.warn('onCreateAccount', this.form.value);
  }

  onCancel() {
    console.warn('onCancel');
  }

  validateMatchingPasswords() {
    if (this.form.get('password')?.value !== this.form.get('confirmPassword')?.value) {
      this.form.get('confirmPassword')?.setErrors({ passwordMismatch: true });
    }
  }
}
