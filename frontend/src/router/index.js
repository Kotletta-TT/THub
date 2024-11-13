import { createRouter, createWebHistory } from 'vue-router';
import LoginPage from '../views/LoginPage.vue';
import MainPage from '../views/MainPage.vue';
import TerminalPage from '../views/TerminalPage.vue';

const routes = [
    { path: '/', name: 'Login', component: LoginPage },
    { path: '/main', name: 'Main', component: MainPage },
    { path: '/terminal/:nodeId', name: 'Terminal', component: TerminalPage, props: true },
];

const router = createRouter({
    history: createWebHistory(process.env.BASE_URL),
    routes
});

export default router;
