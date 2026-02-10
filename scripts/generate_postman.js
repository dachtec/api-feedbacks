#!/usr/bin/env node
const fs = require('fs');

const BASE = '{{base_url}}';
const KEY = '{{api_key}}';

function url(path, query) {
    const parts = path.split('/').filter(Boolean);
    const raw = query
        ? `${BASE}/${path}?${Object.entries(query).map(([k, v]) => `${k}=${v}`).join('&')}`
        : `${BASE}/${path}`;
    const obj = { raw, host: [BASE], path: parts };
    if (query) obj.query = Object.entries(query).map(([key, value]) => ({ key, value: String(value) }));
    return obj;
}

function authHeader() {
    return [{ key: 'X-API-Key', value: KEY }];
}

function jsonHeaders() {
    return [{ key: 'Content-Type', value: 'application/json' }, ...authHeader()];
}

function req(name, method, urlObj, headers, body, tests) {
    const item = { name, request: { method, header: headers, url: urlObj } };
    if (body !== null) item.request.body = { mode: 'raw', raw: typeof body === 'string' ? body : JSON.stringify(body, null, 4) };
    if (tests) item.event = [{ listen: 'test', script: { exec: tests } }];
    return item;
}

function t(code) { return code.split('\n').map(l => l.trim()).filter(Boolean); }

const collection = {
    info: {
        name: 'API Feedbacks - Completa',
        description: 'Colección completa para probar la API de Feedbacks.\n\n## Setup\n1. Levantar la API: `docker compose up -d --build`\n2. Usar las variables `base_url` y `api_key` configuradas.\n\n## Flujo recomendado\n1. Health Check / Ready\n2. Crear varios Feedbacks\n3. Listar con filtros\n4. Obtener por ID\n5. Actualizar\n6. Probar escenarios de error',
        schema: 'https://schema.getpostman.com/json/collection/v2.1.0/collection.json'
    },
    variable: [
        { key: 'base_url', value: 'http://localhost:8080', type: 'string' },
        { key: 'api_key', value: 'my-secret-api-key', type: 'string' },
        { key: 'feedback_id', value: '', type: 'string' },
        { key: 'feedback_id_2', value: '', type: 'string' }
    ],
    item: []
};

// 01 - Health & Readiness
collection.item.push({
    name: '01 - Health & Readiness',
    item: [
        req('Health Check', 'GET', url('health'), [], null, t(`
      pm.test('Status 200', () => pm.response.to.have.status(200));
      pm.test('Status is ok', () => {
        const body = pm.response.json();
        pm.expect(body.status).to.eql('ok');
      });
    `)),
        req('Readiness Check', 'GET', url('ready'), [], null, t(`
      pm.test('Status 200', () => pm.response.to.have.status(200));
      pm.test('Status is ready', () => {
        const body = pm.response.json();
        pm.expect(body.status).to.eql('ready');
      });
    `))
    ]
});

// 02 - Autenticación (Excepciones)
collection.item.push({
    name: '02 - Autenticación (Excepciones)',
    item: [
        req('Sin API Key → 401', 'GET', url('api/v1/feedbacks'), [], null, t(`
      pm.test('Status 401', () => pm.response.to.have.status(401));
      pm.test('Error UNAUTHORIZED', () => {
        const body = pm.response.json();
        pm.expect(body.success).to.be.false;
        pm.expect(body.error.code).to.eql('UNAUTHORIZED');
        pm.expect(body.error.message).to.eql('missing API key');
      });
    `)),
        req('API Key Inválida → 401', 'GET', url('api/v1/feedbacks'),
            [{ key: 'X-API-Key', value: 'wrong-api-key' }], null, t(`
      pm.test('Status 401', () => pm.response.to.have.status(401));
      pm.test('Error UNAUTHORIZED', () => {
        const body = pm.response.json();
        pm.expect(body.success).to.be.false;
        pm.expect(body.error.code).to.eql('UNAUTHORIZED');
        pm.expect(body.error.message).to.eql('invalid API key');
      });
    `)),
        req('API Key vacía → 401', 'GET', url('api/v1/feedbacks'),
            [{ key: 'X-API-Key', value: '' }], null, t(`
      pm.test('Status 401', () => pm.response.to.have.status(401));
      pm.test('Error UNAUTHORIZED', () => {
        const body = pm.response.json();
        pm.expect(body.success).to.be.false;
        pm.expect(body.error.code).to.eql('UNAUTHORIZED');
      });
    `))
    ]
});

// 03 - Crear Feedback (Happy Path)
const feedbackTypes = [
    { type: 'bug', userId: 'u-001', rating: 2, comment: 'El botón de pago no responde en Safari', saveTo: 'feedback_id' },
    { type: 'sugerencia', userId: 'u-002', rating: 4, comment: 'Sería útil poder exportar reportes a PDF', saveTo: 'feedback_id_2' },
    { type: 'elogio', userId: 'u-001', rating: 5, comment: 'Excelente experiencia, la plataforma es muy intuitiva', saveTo: null },
    { type: 'duda', userId: 'u-003', rating: 3, comment: '¿Cómo puedo cambiar mi método de pago?', saveTo: null },
    { type: 'queja', userId: 'u-004', rating: 1, comment: 'El tiempo de carga es inaceptable', saveTo: null }
];

collection.item.push({
    name: '03 - Crear Feedback (Happy Path)',
    item: feedbackTypes.map(ft => {
        const saveScript = ft.saveTo
            ? `pm.collectionVariables.set('${ft.saveTo}', body.data.feedback_id);`
            : '';
        return req(`Crear Feedback tipo ${ft.type}`, 'POST', url('api/v1/feedbacks'), jsonHeaders(),
            { user_id: ft.userId, feedback_type: ft.type, rating: ft.rating, comment: ft.comment },
            t(`
        pm.test('Status 201', () => pm.response.to.have.status(201));
        pm.test('Feedback creado correctamente', () => {
          const body = pm.response.json();
          pm.expect(body.success).to.be.true;
          pm.expect(body.data.user_id).to.eql('${ft.userId}');
          pm.expect(body.data.feedback_type).to.eql('${ft.type}');
          pm.expect(body.data.rating).to.eql(${ft.rating});
          pm.expect(body.data.feedback_id).to.match(/^f-\\d{4}$/);
          pm.expect(body.data.created_at).to.be.a('string');
          pm.expect(body.data.updated_at).to.be.a('string');
          ${saveScript}
        });
      `)
        );
    })
});

// 04 - Crear Feedback (Excepciones)
collection.item.push({
    name: '04 - Crear Feedback (Excepciones)',
    item: [
        req('Body vacío → 400', 'POST', url('api/v1/feedbacks'), jsonHeaders(), {}, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Errores de validación', () => {
        const body = pm.response.json();
        pm.expect(body.success).to.be.false;
      });
    `)),
        req('Sin user_id → 400', 'POST', url('api/v1/feedbacks'), jsonHeaders(),
            { feedback_type: 'bug', rating: 3, comment: 'Falta el user_id' }, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Validation error user_id', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('VALIDATION_ERROR');
        const fields = body.error.details.map(d => d.field);
        pm.expect(fields).to.include('user_id');
      });
    `)),
        req('user_id formato inválido (usr-001) → 400', 'POST', url('api/v1/feedbacks'), jsonHeaders(),
            { user_id: 'usr-001', feedback_type: 'bug', rating: 3, comment: 'Formato incorrecto de user_id' }, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Validation error user_id format', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('VALIDATION_ERROR');
        const fields = body.error.details.map(d => d.field);
        pm.expect(fields).to.include('user_id');
      });
    `)),
        req('user_id formato inválido (texto libre) → 400', 'POST', url('api/v1/feedbacks'), jsonHeaders(),
            { user_id: 'john_doe', feedback_type: 'bug', rating: 3, comment: 'user_id no cumple u-###' }, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Validation error user_id', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('VALIDATION_ERROR');
        const fields = body.error.details.map(d => d.field);
        pm.expect(fields).to.include('user_id');
      });
    `)),
        req('feedback_type inválido → 400', 'POST', url('api/v1/feedbacks'), jsonHeaders(),
            { user_id: 'u-001', feedback_type: 'invalid_type', rating: 3, comment: 'Tipo inválido' }, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Validation error feedback_type', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('VALIDATION_ERROR');
        const fields = body.error.details.map(d => d.field);
        pm.expect(fields).to.include('feedback_type');
      });
    `)),
        req('feedback_type en inglés (suggestion) → 400', 'POST', url('api/v1/feedbacks'), jsonHeaders(),
            { user_id: 'u-001', feedback_type: 'suggestion', rating: 3, comment: 'Tipo en inglés no válido' }, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Validation error feedback_type', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('VALIDATION_ERROR');
        const fields = body.error.details.map(d => d.field);
        pm.expect(fields).to.include('feedback_type');
      });
    `)),
        req('Rating fuera de rango (0) → 400', 'POST', url('api/v1/feedbacks'), jsonHeaders(),
            { user_id: 'u-001', feedback_type: 'bug', rating: 0, comment: 'Rating menor a 1' }, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Validation error rating', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('VALIDATION_ERROR');
        const fields = body.error.details.map(d => d.field);
        pm.expect(fields).to.include('rating');
      });
    `)),
        req('Rating fuera de rango (6) → 400', 'POST', url('api/v1/feedbacks'), jsonHeaders(),
            { user_id: 'u-001', feedback_type: 'bug', rating: 6, comment: 'Rating mayor a 5' }, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Validation error rating', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('VALIDATION_ERROR');
        const fields = body.error.details.map(d => d.field);
        pm.expect(fields).to.include('rating');
      });
    `)),
        req('Rating negativo (-1) → 400', 'POST', url('api/v1/feedbacks'), jsonHeaders(),
            { user_id: 'u-001', feedback_type: 'bug', rating: -1, comment: 'Rating negativo' }, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Validation error rating', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('VALIDATION_ERROR');
        const fields = body.error.details.map(d => d.field);
        pm.expect(fields).to.include('rating');
      });
    `)),
        req('Comentario vacío → 400', 'POST', url('api/v1/feedbacks'), jsonHeaders(),
            { user_id: 'u-001', feedback_type: 'bug', rating: 3, comment: '' }, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Validation error comment', () => {
        const body = pm.response.json();
        const fields = body.error.details.map(d => d.field);
        pm.expect(fields).to.include('comment');
      });
    `)),
        req('Comentario solo espacios → 400', 'POST', url('api/v1/feedbacks'), jsonHeaders(),
            { user_id: 'u-001', feedback_type: 'bug', rating: 3, comment: '   ' }, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Validation error comment', () => {
        const body = pm.response.json();
        const fields = body.error.details.map(d => d.field);
        pm.expect(fields).to.include('comment');
      });
    `)),
        req('JSON malformado → 400', 'POST', url('api/v1/feedbacks'), jsonHeaders(),
            '{ invalid json }', t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Error INVALID_JSON', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('INVALID_JSON');
      });
    `)),
        req('Múltiples errores de validación → 400', 'POST', url('api/v1/feedbacks'), jsonHeaders(),
            { user_id: '', feedback_type: 'nonexistent', rating: 99, comment: '' }, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Multiple validation errors', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('VALIDATION_ERROR');
        pm.expect(body.error.details.length).to.be.greaterThan(1);
      });
    `)),
        req('Comentario excede 2000 caracteres → 400', 'POST', url('api/v1/feedbacks'), jsonHeaders(),
            { user_id: 'u-001', feedback_type: 'bug', rating: 3, comment: 'A'.repeat(2001) }, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Validation error comment length', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('VALIDATION_ERROR');
        const fields = body.error.details.map(d => d.field);
        pm.expect(fields).to.include('comment');
      });
    `))
    ]
});

// 05 - Obtener Feedback por ID
collection.item.push({
    name: '05 - Obtener Feedback por ID',
    item: [
        req('Obtener por ID (Happy Path)', 'GET', url('api/v1/feedbacks/{{feedback_id}}'), authHeader(), null, t(`
      pm.test('Status 200', () => pm.response.to.have.status(200));
      pm.test('Feedback retornado', () => {
        const body = pm.response.json();
        pm.expect(body.success).to.be.true;
        pm.expect(body.data.feedback_id).to.eql(pm.collectionVariables.get('feedback_id'));
        pm.expect(body.data.feedback_id).to.match(/^f-\\d{4}$/);
        pm.expect(body.data.user_id).to.match(/^u-\\d{3}$/);
        pm.expect(body.data).to.have.all.keys('feedback_id','user_id','feedback_type','rating','comment','created_at','updated_at');
      });
    `)),
        req('ID inexistente → 404', 'GET', url('api/v1/feedbacks/f-9999'), authHeader(), null, t(`
      pm.test('Status 404', () => pm.response.to.have.status(404));
      pm.test('Error NOT_FOUND', () => {
        const body = pm.response.json();
        pm.expect(body.success).to.be.false;
        pm.expect(body.error.code).to.eql('NOT_FOUND');
      });
    `)),
        req('ID con formato inválido → 404', 'GET', url('api/v1/feedbacks/abc-invalid'), authHeader(), null, t(`
      pm.test('Status 404', () => pm.response.to.have.status(404));
      pm.test('Error NOT_FOUND', () => {
        const body = pm.response.json();
        pm.expect(body.success).to.be.false;
        pm.expect(body.error.code).to.eql('NOT_FOUND');
      });
    `))
    ]
});

// 06 - Actualizar Feedback (Happy Path)
collection.item.push({
    name: '06 - Actualizar Feedback (Happy Path)',
    item: [
        req('Actualizar rating solamente', 'PATCH', url('api/v1/feedbacks/{{feedback_id}}'), jsonHeaders(),
            { rating: 5 }, t(`
      pm.test('Status 200', () => pm.response.to.have.status(200));
      pm.test('Rating actualizado', () => {
        const body = pm.response.json();
        pm.expect(body.data.rating).to.eql(5);
        pm.expect(body.data.feedback_id).to.eql(pm.collectionVariables.get('feedback_id'));
      });
    `)),
        req('Actualizar comment solamente', 'PATCH', url('api/v1/feedbacks/{{feedback_id}}'), jsonHeaders(),
            { comment: 'Problema resuelto, gracias' }, t(`
      pm.test('Status 200', () => pm.response.to.have.status(200));
      pm.test('Comment actualizado', () => {
        const body = pm.response.json();
        pm.expect(body.data.comment).to.eql('Problema resuelto, gracias');
      });
    `)),
        req('Actualizar feedback_type solamente', 'PATCH', url('api/v1/feedbacks/{{feedback_id}}'), jsonHeaders(),
            { feedback_type: 'elogio' }, t(`
      pm.test('Status 200', () => pm.response.to.have.status(200));
      pm.test('Type actualizado a elogio', () => {
        const body = pm.response.json();
        pm.expect(body.data.feedback_type).to.eql('elogio');
      });
    `)),
        req('Actualizar múltiples campos', 'PATCH', url('api/v1/feedbacks/{{feedback_id}}'), jsonHeaders(),
            { feedback_type: 'sugerencia', rating: 4, comment: 'Actualización final' }, t(`
      pm.test('Status 200', () => pm.response.to.have.status(200));
      pm.test('Campos actualizados', () => {
        const body = pm.response.json();
        pm.expect(body.data.rating).to.eql(4);
        pm.expect(body.data.comment).to.eql('Actualización final');
        pm.expect(body.data.feedback_type).to.eql('sugerencia');
      });
    `)),
        req('Verificar updated_at cambia', 'PATCH', url('api/v1/feedbacks/{{feedback_id}}'), jsonHeaders(),
            { rating: 3 }, t(`
      pm.test('Status 200', () => pm.response.to.have.status(200));
      pm.test('updated_at es reciente', () => {
        const body = pm.response.json();
        const updatedAt = new Date(body.data.updated_at);
        const now = new Date();
        const diffMs = now - updatedAt;
        pm.expect(diffMs).to.be.below(60000);
      });
    `))
    ]
});

// 07 - Actualizar Feedback (Excepciones)
collection.item.push({
    name: '07 - Actualizar Feedback (Excepciones)',
    item: [
        req('Body vacío (sin campos) → 400', 'PATCH', url('api/v1/feedbacks/{{feedback_id}}'), jsonHeaders(),
            {}, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Error EMPTY_UPDATE', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('EMPTY_UPDATE');
      });
    `)),
        req('ID inexistente → 404', 'PATCH', url('api/v1/feedbacks/f-9999'), jsonHeaders(),
            { rating: 5 }, t(`
      pm.test('Status 404', () => pm.response.to.have.status(404));
      pm.test('Error NOT_FOUND', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('NOT_FOUND');
      });
    `)),
        req('Rating inválido en update → 400', 'PATCH', url('api/v1/feedbacks/{{feedback_id}}'), jsonHeaders(),
            { rating: 10 }, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Validation error', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('VALIDATION_ERROR');
      });
    `)),
        req('feedback_type inválido en update → 400', 'PATCH', url('api/v1/feedbacks/{{feedback_id}}'), jsonHeaders(),
            { feedback_type: 'invalid' }, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Validation error', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('VALIDATION_ERROR');
      });
    `)),
        req('Comment vacío en update → 400', 'PATCH', url('api/v1/feedbacks/{{feedback_id}}'), jsonHeaders(),
            { comment: '' }, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Validation error comment', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('VALIDATION_ERROR');
      });
    `)),
        req('JSON malformado en update → 400', 'PATCH', url('api/v1/feedbacks/{{feedback_id}}'), jsonHeaders(),
            'not a json', t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Error INVALID_JSON', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('INVALID_JSON');
      });
    `))
    ]
});

// 08 - Listar Feedbacks con Filtros (Happy Path)
collection.item.push({
    name: '08 - Listar Feedbacks con Filtros (Happy Path)',
    item: [
        req('Listar todos (sin filtros)', 'GET', url('api/v1/feedbacks'), authHeader(), null, t(`
      pm.test('Status 200', () => pm.response.to.have.status(200));
      pm.test('Respuesta con paginación', () => {
        const body = pm.response.json();
        pm.expect(body.success).to.be.true;
        pm.expect(body.data).to.be.an('array');
        pm.expect(body.meta).to.have.property('total');
        pm.expect(body.meta).to.have.property('limit');
        pm.expect(body.meta).to.have.property('offset');
      });
    `)),
        req('Filtrar por user_id', 'GET', url('api/v1/feedbacks', { user_id: 'u-001' }), authHeader(), null, t(`
      pm.test('Status 200', () => pm.response.to.have.status(200));
      pm.test('Todos son del user_id', () => {
        const body = pm.response.json();
        body.data.forEach(f => pm.expect(f.user_id).to.eql('u-001'));
      });
    `)),
        req('Filtrar por feedback_type', 'GET', url('api/v1/feedbacks', { feedback_type: 'bug' }), authHeader(), null, t(`
      pm.test('Status 200', () => pm.response.to.have.status(200));
      pm.test('Todos son tipo bug', () => {
        const body = pm.response.json();
        body.data.forEach(f => pm.expect(f.feedback_type).to.eql('bug'));
      });
    `)),
        req('Filtrar por feedback_type sugerencia', 'GET', url('api/v1/feedbacks', { feedback_type: 'sugerencia' }), authHeader(), null, t(`
      pm.test('Status 200', () => pm.response.to.have.status(200));
      pm.test('Todos son tipo sugerencia', () => {
        const body = pm.response.json();
        body.data.forEach(f => pm.expect(f.feedback_type).to.eql('sugerencia'));
      });
    `)),
        req('Filtrar por rango de rating', 'GET', url('api/v1/feedbacks', { min_rating: '3', max_rating: '5' }), authHeader(), null, t(`
      pm.test('Status 200', () => pm.response.to.have.status(200));
      pm.test('Ratings en rango 3-5', () => {
        const body = pm.response.json();
        body.data.forEach(f => {
          pm.expect(f.rating).to.be.at.least(3);
          pm.expect(f.rating).to.be.at.most(5);
        });
      });
    `)),
        req('Filtrar por rango de fechas', 'GET', url('api/v1/feedbacks', { created_from: '2026-01-01T00:00:00Z', created_to: '2026-12-31T23:59:59Z' }), authHeader(), null, t(`
      pm.test('Status 200', () => pm.response.to.have.status(200));
      pm.test('Respuesta exitosa', () => {
        const body = pm.response.json();
        pm.expect(body.success).to.be.true;
      });
    `)),
        req('Filtrar solo min_rating', 'GET', url('api/v1/feedbacks', { min_rating: '4' }), authHeader(), null, t(`
      pm.test('Status 200', () => pm.response.to.have.status(200));
      pm.test('Ratings >= 4', () => {
        const body = pm.response.json();
        body.data.forEach(f => pm.expect(f.rating).to.be.at.least(4));
      });
    `)),
        req('Filtrar solo max_rating', 'GET', url('api/v1/feedbacks', { max_rating: '2' }), authHeader(), null, t(`
      pm.test('Status 200', () => pm.response.to.have.status(200));
      pm.test('Ratings <= 2', () => {
        const body = pm.response.json();
        body.data.forEach(f => pm.expect(f.rating).to.be.at.most(2));
      });
    `)),
        req('Filtros combinados (user_id + type + rating)', 'GET',
            url('api/v1/feedbacks', { user_id: 'u-001', min_rating: '1', max_rating: '3', feedback_type: 'bug' }),
            authHeader(), null, t(`
      pm.test('Status 200', () => pm.response.to.have.status(200));
      pm.test('Filtros combinados aplicados', () => {
        const body = pm.response.json();
        pm.expect(body.success).to.be.true;
        body.data.forEach(f => {
          pm.expect(f.user_id).to.eql('u-001');
          pm.expect(f.feedback_type).to.eql('bug');
          pm.expect(f.rating).to.be.at.least(1);
          pm.expect(f.rating).to.be.at.most(3);
        });
      });
    `)),
        req('Paginación con limit y offset', 'GET', url('api/v1/feedbacks', { limit: '2', offset: '0' }), authHeader(), null, t(`
      pm.test('Status 200', () => pm.response.to.have.status(200));
      pm.test('Paginación aplicada', () => {
        const body = pm.response.json();
        pm.expect(body.meta.limit).to.eql(2);
        pm.expect(body.meta.offset).to.eql(0);
        pm.expect(body.data.length).to.be.at.most(2);
      });
    `)),
        req('Paginación segunda página', 'GET', url('api/v1/feedbacks', { limit: '2', offset: '2' }), authHeader(), null, t(`
      pm.test('Status 200', () => pm.response.to.have.status(200));
      pm.test('Segunda página', () => {
        const body = pm.response.json();
        pm.expect(body.meta.limit).to.eql(2);
        pm.expect(body.meta.offset).to.eql(2);
      });
    `))
    ]
});

// 09 - Listar Feedbacks (Excepciones)
collection.item.push({
    name: '09 - Listar Feedbacks (Excepciones)',
    item: [
        req('feedback_type inválido en filtro → 400', 'GET', url('api/v1/feedbacks', { feedback_type: 'invalid' }),
            authHeader(), null, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Error INVALID_FILTER', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('INVALID_FILTER');
      });
    `)),
        req('feedback_type en inglés en filtro → 400', 'GET', url('api/v1/feedbacks', { feedback_type: 'suggestion' }),
            authHeader(), null, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Error INVALID_FILTER', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('INVALID_FILTER');
      });
    `)),
        req('min_rating no numérico → 400', 'GET', url('api/v1/feedbacks', { min_rating: 'abc' }),
            authHeader(), null, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Error INVALID_FILTER', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('INVALID_FILTER');
      });
    `)),
        req('max_rating fuera de rango → 400', 'GET', url('api/v1/feedbacks', { max_rating: '10' }),
            authHeader(), null, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Error INVALID_FILTER', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('INVALID_FILTER');
      });
    `)),
        req('min_rating fuera de rango (0) → 400', 'GET', url('api/v1/feedbacks', { min_rating: '0' }),
            authHeader(), null, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Error INVALID_FILTER', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('INVALID_FILTER');
      });
    `)),
        req('created_from formato inválido → 400', 'GET', url('api/v1/feedbacks', { created_from: '2026-01-01' }),
            authHeader(), null, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Error INVALID_FILTER', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('INVALID_FILTER');
      });
    `)),
        req('created_to formato inválido → 400', 'GET', url('api/v1/feedbacks', { created_to: 'not-a-date' }),
            authHeader(), null, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Error INVALID_FILTER', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('INVALID_FILTER');
      });
    `)),
        req('limit negativo → 400', 'GET', url('api/v1/feedbacks', { limit: '-1' }),
            authHeader(), null, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Error INVALID_FILTER', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('INVALID_FILTER');
      });
    `)),
        req('limit no numérico → 400', 'GET', url('api/v1/feedbacks', { limit: 'abc' }),
            authHeader(), null, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Error INVALID_FILTER', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('INVALID_FILTER');
      });
    `)),
        req('offset negativo → 400', 'GET', url('api/v1/feedbacks', { offset: '-5' }),
            authHeader(), null, t(`
      pm.test('Status 400', () => pm.response.to.have.status(400));
      pm.test('Error INVALID_FILTER', () => {
        const body = pm.response.json();
        pm.expect(body.error.code).to.eql('INVALID_FILTER');
      });
    `)),
        req('user_id sin resultados → 200 vacío', 'GET', url('api/v1/feedbacks', { user_id: 'u-999' }),
            authHeader(), null, t(`
      pm.test('Status 200', () => pm.response.to.have.status(200));
      pm.test('Array vacío', () => {
        const body = pm.response.json();
        pm.expect(body.success).to.be.true;
        pm.expect(body.data).to.be.an('array').that.is.empty;
        pm.expect(body.meta.total).to.eql(0);
      });
    `)),
        req('Fechas sin resultados → 200 vacío', 'GET',
            url('api/v1/feedbacks', { created_from: '2020-01-01T00:00:00Z', created_to: '2020-01-02T00:00:00Z' }),
            authHeader(), null, t(`
      pm.test('Status 200', () => pm.response.to.have.status(200));
      pm.test('Array vacío', () => {
        const body = pm.response.json();
        pm.expect(body.success).to.be.true;
        pm.expect(body.data).to.be.an('array').that.is.empty;
      });
    `))
    ]
});

// 10 - Ruta no encontrada
collection.item.push({
    name: '10 - Rutas No Encontradas',
    item: [
        req('Ruta inexistente → 404/405', 'GET', url('api/v1/nonexistent'), authHeader(), null, t(`
      pm.test('Status 404 o 405', () => {
        pm.expect(pm.response.code).to.be.oneOf([404, 405]);
      });
    `)),
        req('DELETE no soportado → 405', 'DELETE', url('api/v1/feedbacks/{{feedback_id}}'), authHeader(), null, t(`
      pm.test('Status 405', () => pm.response.to.have.status(405));
    `)),
        req('PUT no soportado → 405', 'PUT', url('api/v1/feedbacks/{{feedback_id}}'), jsonHeaders(),
            { user_id: 'u-001', feedback_type: 'bug', rating: 3, comment: 'PUT no soportado' }, t(`
      pm.test('Status 405', () => pm.response.to.have.status(405));
    `))
    ]
});

const output = JSON.stringify(collection, null, 4);
const outputPath = __dirname + '/../docs/API_Feedbacks.postman_collection.json';
fs.writeFileSync(outputPath, output + '\n');
console.log(`✅ Colección generada: ${outputPath}`);
console.log(`   Total carpetas: ${collection.item.length}`);
const totalReqs = collection.item.reduce((sum, folder) => sum + folder.item.length, 0);
console.log(`   Total requests: ${totalReqs}`);
