import { Routes } from '@angular/router';

import { Dashboard } from './components/dashboard';
import { VideoList } from './components/video-list';
import { Download } from './components/download';

export const JAYE_ROUTES: Routes = [
    { path: '', component: Dashboard, pathMatch: 'full' },
    { path: 'videos', component: VideoList },
    { path: 'download', component: Download },
    { path: '**', redirectTo: '' },
];
