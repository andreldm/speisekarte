const express = require('express');
const bodyParser = require('body-parser');
const fs = require("fs");
const app = express();
const rq = require('request-promise');

app.use(express.static('public'));
app.use(bodyParser.json());

const envvar = process.env.SPEISEKARTE_CONFIG;

const config = JSON.parse(fs.readFileSync(envvar || "config.json"));
const categories = config.categories;
const categoriesMap = new Map(categories.map(s => [s.key, s]));
const services = config.services;
const servicesMap = new Map(services.map(s => [s.key, s]));

app.get('/categories', (req, res) => {
    res.send(categories);
});

app.get('/services', (req, res) => {
    res.send(services);
});

app.get('/service/:key', (req, res) => {
    res.send(servicesMap.get(req.params.key));
});

app.get('/service/:key/info', (req, res) => {
    const service = servicesMap.get(req.params.key);
    const info = {};

    rq({uri: service.url, timeout: 10000 })
        .then(() => info.status = 'online')
        .catch(() => info.status = 'offline')
        .finally(() => res.send(info));
});

app.listen(3000, () => console.log('Speisekarte ist ready zu serve at port 3000'));
