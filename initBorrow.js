use admin
db.createUser(
    {
        user: "bomongo",
        pwd: "bomongo",
        roles: [ { role: "userAdminAnyDatabase", db: "admin" }, "readWriteAnyDatabase" ]
    }
);