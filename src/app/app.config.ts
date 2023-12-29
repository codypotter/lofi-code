import { ApplicationConfig, importProvidersFrom } from '@angular/core';
import { provideRouter } from '@angular/router';

import { routes } from './app.routes';
import { initializeApp, provideFirebaseApp } from '@angular/fire/app';
import { getAuth, provideAuth } from '@angular/fire/auth';
import { getFirestore, provideFirestore } from '@angular/fire/firestore';

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    importProvidersFrom(
      provideFirebaseApp(() => initializeApp({
        apiKey: "AIzaSyBkk1syx5c5qVmPKFoLrOXkDwctSozfeL4",
        authDomain: "lofi-code.firebaseapp.com",
        projectId: "lofi-code",
        storageBucket: "lofi-code.appspot.com",
        messagingSenderId: "489391062135",
        appId: "1:489391062135:web:06311425fa89d443119a72",
        measurementId: "G-G4Z07Y8ZC3"
      })),
    ),
    importProvidersFrom(provideAuth(() => getAuth())),
    importProvidersFrom(provideFirestore(() => getFirestore())),
  ]
};
