import { ApplicationConfig, importProvidersFrom } from '@angular/core';
import { provideRouter } from '@angular/router';

import { routes } from './app.routes';
import { initializeApp, provideFirebaseApp } from '@angular/fire/app';
import { getAuth, provideAuth } from '@angular/fire/auth';
import { getFirestore, provideFirestore } from '@angular/fire/firestore/lite';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { BreadcrumbService } from './services/breadcrumb.service';
import { PostsService } from './services/posts.service';
import { YouTubePlayerModule } from '@angular/youtube-player';
import { environment } from 'src/environments/environment';
import { LoggerModule } from 'ngx-logger';
import { TagsService } from './services/tags.service';
import { UsersService } from './services/users.service';
import { getStorage, provideStorage } from '@angular/fire/storage';
import { TimeagoModule } from 'ngx-timeago';
import { ShareButtonsModule } from 'ngx-sharebuttons/buttons';
import { provideClientHydration } from '@angular/platform-browser';
import { provideHttpClient } from '@angular/common/http';

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    importProvidersFrom(BrowserAnimationsModule),
    importProvidersFrom(LoggerModule.forRoot({ level: environment.logLevel,  })),
    importProvidersFrom(ShareButtonsModule.withConfig({
      include: ['facebook', 'twitter', 'linkedin', 'reddit', 'email', 'tumblr'],
    })),
    provideFirebaseApp(() => initializeApp(environment.firebaseOptions)),
    provideAuth(() => getAuth()),
    provideFirestore(() => getFirestore()),
    provideStorage(() => getStorage()),
    importProvidersFrom(YouTubePlayerModule),
    importProvidersFrom(TimeagoModule.forRoot()),
    BreadcrumbService,
    PostsService,
    TagsService,
    UsersService,
    provideClientHydration(),
    provideHttpClient(),
  ],
};
