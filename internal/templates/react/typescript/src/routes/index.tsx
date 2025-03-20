import { lazyLoad } from '../utils/lazyLoad';
import LoadingSpinner from '../components/LoadingSpinner';

export const routes = [
  {
    path: '/',
    component: lazyLoad(() => import('../pages/Home'), <LoadingSpinner />),
  },
  {
    path: '/dashboard',
    component: lazyLoad(() => import('../pages/Dashboard'), <LoadingSpinner />),
  },
  {
    path: '/settings',
    component: lazyLoad(() => import('../pages/Settings'), <LoadingSpinner />),
  },
];