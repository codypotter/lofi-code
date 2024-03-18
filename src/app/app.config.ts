import { ApplicationConfig, importProvidersFrom } from '@angular/core';
import { provideRouter } from '@angular/router';

import { routes } from './app.routes';
import { initializeApp, provideFirebaseApp } from '@angular/fire/app';
import { getAuth, provideAuth } from '@angular/fire/auth';
import { getFirestore, provideFirestore } from '@angular/fire/firestore/lite';
import { NgxSpinnerModule } from 'ngx-spinner';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { BreadcrumbService } from './services/breadcrumb.service';
import { PostsService } from './services/posts.service';
import { YouTubePlayerModule } from '@angular/youtube-player';
import { HttpClientModule } from '@angular/common/http';
import { environment } from 'src/environments/environment';
import { LoggerModule } from 'ngx-logger';
import { TagsService } from './services/tags.service';
import { UsersService } from './services/users.service';

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    importProvidersFrom(BrowserAnimationsModule, HttpClientModule),
    importProvidersFrom(NgxSpinnerModule.forRoot({ type: 'ball-scale-multiple' })),
    importProvidersFrom(LoggerModule.forRoot({ level: environment.logLevel,  })),
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
    importProvidersFrom(YouTubePlayerModule),
    BreadcrumbService,
    PostsService,
    TagsService,
    UsersService,
  ],
};
