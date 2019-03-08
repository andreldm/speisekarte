const express = require('express');
const bodyParser = require('body-parser');
const fs = require("fs");
const rq = require('request-promise');

const app = express();

app.use(express.static('public'));
app.use(bodyParser.json());

const envvar = process.env.SPEISEKARTE_CONFIG;
const config = JSON.parse(fs.readFileSync(envvar || 'config.json'));
const port = config.port || 3000;
const timeout = config.timeout || 10000;
const title = config.title || "Speisekarte";

const compact = config.compact || false;
const categories = config.categories;
const services = config.services;
const servicesMap = new Map(services.map(s => [s.key, s]));

app.get('/config', (req, res) => {
    res.send({
        compact: compact,
        categories: categories,
        services: services,
        title: title
    });
});

app.get('/service/:key/info', (req, res) => {
    const service = servicesMap.get(req.params.key);
    const info = {};

    rq({uri: service.url, timeout: timeout, simple: false })
        .then(() => info.status = 'online')
        .catch(() => info.status = 'offline')
        .finally(() => res.send(info));
});

app.listen(port, () => console.log(`Speisekarte ist ready zu serve at http://localhost:${port}`));
