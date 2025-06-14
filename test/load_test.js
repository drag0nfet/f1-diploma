import http from 'k6/http';
import { check, sleep } from 'k6';
import encoding from 'k6/encoding';
import crypto from 'k6/crypto';

const secret = "Wv1%`j9pr]0d[s'_HwX,U|m;6^3>u="; // JWT secret как на сервере

function encodeBase64(str) {
    return encoding.b64encode(str, 'std'); // используем k6/encoding
}

function generateJWT(username, userId, rights) {
    const header = {
        alg: 'HS256',
        typ: 'JWT'
    };

    const payload = {
        username: username,
        user_id: userId,
        rights: rights,
        exp: Math.floor(Date.now() / 1000) + 3600
    };

    const headerEncoded = encodeBase64(JSON.stringify(header));
    const payloadEncoded = encodeBase64(JSON.stringify(payload));
    const data = `${headerEncoded}.${payloadEncoded}`;

    const signature = crypto.hmac('sha256', data, secret, 'base64');

    return `${data}.${signature}`;
}

function uuidv4() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        const r = Math.random() * 16 | 0,
            v = c === 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

export let options = {
    stages: [
        {duration: '10s', target: 180},
        {duration: '60s', target: 180},
        {duration: '10s', target: 0},
    ]
};

const BASE_URL = __ENV.BASE_URL;

const topicIds = [4, 10, 2, 7, 1, 8, 11];

export default function () {
    // 1. Register
    const id = uuidv4();
    const username = `K6TEST_${id}`;
    const password = "password123";
    const email = `user_${id}@example.com`;

    let res1 = http.post(`${BASE_URL}/register`, JSON.stringify({
        username: username,
        password: password,
        email: email
    }), { headers: { 'Content-Type': 'application/json' } });

    check(res1, { 'register 200': (r) => r.status === 200 });

    const userId = Math.floor(1000 + Math.random() * 1000);
    const token = generateJWT(username, userId, 1);

    const authHeaders = {
        'Content-Type': 'application/json',
        'X-Requested-With': 'XMLHttpRequest',
        'Cookie': `auth=${token}`,
    };

    // 2. Load news
    const status = 'ACTIVE';
    const page = 1;
    const limit = 10;

    let res2 = http.get(`${BASE_URL}/load-news-by-status?status=${status}&page=${page}&limit=${limit}`, {
        headers: authHeaders
    });

    check(res2, { 'load-news 200': (r) => r.status === 200 });

    // 3. Send message
    const topicId = topicIds[Math.floor(Math.random() * topicIds.length)];
    const content = 'ТЕСТОВОЕ_SMS';
    let currentReplyId = null;

    if (Math.random() < 0.4) {
        currentReplyId = 40;
    }

    const msgBody = {
        chat_id: String(topicId),
        content: content
    };
    if (currentReplyId !== null) {
        msgBody.reply_id = String(currentReplyId);
    }

    let res3 = http.post(`${BASE_URL}/send-message`, JSON.stringify(msgBody), { headers: {
            'Content-Type': 'application/json',
            'X-Requested-With': 'XMLHttpRequest',
            'Cookie': `auth=${token}`,
        } });

    check(res3, { 'send-message 200': (r) => r.status === 200 });

    sleep(1); // ожидание 1 сек между итерациями для моделирования средней нагрузки
}
