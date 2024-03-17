import { Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { PostsComponent } from './posts/posts.component';
import { PostComponent } from './post/post.component';
import { TosComponent } from './tos/tos.component';
import { PrivacyPolicyComponent } from './privacy-policy/privacy-policy.component';

export const routes: Routes = [
    {
        path: '', component: HomeComponent
    },
    {
        path: 'posts', component: PostsComponent,
    },
    {
        path: 'posts/:slug', component: PostComponent,
    },
    {
        path: 'tos', component: TosComponent,
    },
    {
        path: 'privacy-policy', component: PrivacyPolicyComponent,
    }
];
