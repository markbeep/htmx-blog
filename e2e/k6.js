import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    vus: 50, // number of virtual users
    duration: '30s', // test duration
};

const BASE_URL = 'https://markc.su';
const PAGES = [
    '/',
    '/posts',
    '/about',
    '/polyring',
    '/posts/build_an_overkill_server',
    '/posts/fusing-go-into-python',
    '/posts/getting_started',
    '/posts/golang-interfaces',
    '/posts/hangman',
    '/posts/learning_nim',
    '/posts/nix-direnv-setup',
    '/posts/nix-package-manager',
    '/posts/nix',
    '/posts/raycasting',
    '/posts/react_games',
    '/posts/threadpool_c',
];

export default function () {
    for (let page of PAGES) {
        let res = http.get(`${BASE_URL}${page}`);
        check(res, {
            'is status 200': (r) => r.status === 200,
            'duration was <= 200ms': (r) => r.timings.duration <= 200,
        });
        sleep(1); // Simulate think time for a real user
    }
}
