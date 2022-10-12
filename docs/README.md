

Initial Credentials
===================

```yaml
users:
- username: admin
  hash: $2a$12$VcCDgh2NDk07JGN0rjGbM.Ad41qVR/YFJcgHp0UGns5JDymv..TOG ## admin
  isAdmin: true
```

Changing Credentials
==================

Step 1. Create a bcrypt hash
```shell
$ htpasswd -nbBC 12 "" topsecret | tr -d :
$2y$12$KTbZnVgxAIUmnu5W2bRGmuJ/in8A9sHLt2je2lxOriq8TJP0vMk1y
```

Step 2. Edit users.yaml
```yaml
users:
- username: admin
  hash: $2a$12$VcCDgh2NDk07JGN0rjGbM.Ad41qVR/YFJcgHp0UGns5JDymv..TOG ## admin
  isAdmin: true
- username: user1
  hash: $2y$12$KTbZnVgxAIUmnu5W2bRGmuJ/in8A9sHLt2je2lxOriq8TJP0vMk1y ## topsecret
  isAdmin: true
```
