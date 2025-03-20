import { Routes } from '@angular/router';
import { authGuard } from './guards/auth.guard';
import { postResolver } from './resolvers/post.resolver';
import { postsResolver } from './services/posts.resolver';

export const routes: Routes = [
    {
        path: '', loadComponent: () => import('./home/home.component').then(m => m.HomeComponent)
    },
    {
        path: 'posts', loadComponent: () => import('./posts/posts.component').then(m => m.PostsComponent),
        resolve: {
            posts: postsResolver,
        },
    },
    {
        path: 'posts/:slug', loadComponent: () => import('./post/post.component').then(m => m.PostComponent),
        resolve: {
            post: postResolver,
        },
    },
    {
        path: 'tos', loadComponent: () => import('./tos/tos.component').then(m => m.TosComponent),
    },
    {
        path: 'privacy-policy', loadComponent: () => import('./privacy-policy/privacy-policy.component').then(m => m.PrivacyPolicyComponent),
    },
    {
        path: 'settings',
        loadComponent: () => import('./settings/settings.component').then(m => m.SettingsComponent),
        canActivate: [authGuard]
    }
];
