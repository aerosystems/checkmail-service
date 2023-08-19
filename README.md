# checkmail-service

🍕 A part of microservice infrastructure, who responsible for store and check email domains in black/whitelists

📎 Service has used lightweight HTTP Router [Chi](https://github.com/go-chi/chi), idiomatic ORM library for management PostgreSQL [GORM](https://gorm.io/)

📚 Read & Test with [Swagger Docs](http://localhost:8083/docs/index.html)

🎲 Test or Develop with Postman Collection(just import **postman-collection.json** file)

| CODE   | DESCRIPTION                  |
|--------|------------------------------|
| 422001 | could not read request body  |
| 400001 | email address does not valid |
| 400002 | domain does not valid        |
| 400003 | domain does not exist        |
