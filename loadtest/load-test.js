import http from 'k6/http';

export const options = {
    scenarios: {
        gateway_load: { 
            executor: 'constant-arrival-rate',
            rate: 500,
            timeUnit: '1s',
            duration: '200s',
            preAllocatedVUs: 10,
            maxVUs: 50,
        },
    },
};

export default function () {
    http.get(`${__ENV.BASE_URL}/api/products`);
}
