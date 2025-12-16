# Lilac

Bash-based инструмент для наливки серверов, записи образов системы, виртуальных машин, генерации прошивок.
Предназначен для замены ansible и подобных инструментов.


## Суть явления

Шаблонизированый bash с шаблонизироваными конфигами. Синтаксис шаблонов go text/template.
В bash'е неудобные условия и циклы по структурам данных, зато в шаблонизаторе go это сделано удобно.


## Почему не ansible

В Ansible нужно запоминать промежуточный декларативный YAML-based DSL, который не подходит для описания циклов и ветвлений. Шаблонизация конфигов в Ansible сделана хорошо, декларативность, кажется что мешает.


## Установка

```bash
go install github.com/s-tyryshkin/lilac/cmd/lilac@latest
```


## Как пользоваться

Фрагмент для настройки openssh-сервера и выписывания TLS-сертификата. Положите в папку bash-скрипт, данные для шаблонизации и шаблоны конфигов.

`script.sh`
```
#!/usr/bin/env bash

apt update -y
apt upgrade -y
apt install -y openssh-server docker.io

{{ if .openssh.enabled }}
    systemctl enable sshd
    render /etc/ssh/sshd_config
{{ end }}

{{ if .tls.enabled }}
    systemctl enable docker
    systemctl start docker

    if [ ! -e "/etc/letsencrypt/live/www.example.com/privkey.pem" ]; then
        docker run \
            -p 80:80 \
            --rm \
            --name certbot \
            -v /etc/letsencrypt:/etc/letsencrypt \
            -v /var/lib/letsencrypt:/var/lib/letsencrypt \
            certbot/certbot \
            certonly \
                --noninteractive \
                --agree-tos \
                -m contact@example.com \
                -d www.example.com \
                --standalone
    fi
{{ end }}
```

`values.yaml`
```yaml
ssh:
    enabled: true
    port: 2222

tls:
    enabled: true
```

`etc/ssh/sshd_config`
```
Port {{ .ssh.port }}
```

Для полноценной настройки сервера потребуется описать одним или несколькими bash-скриптами процесс настройки, если вы сделаете скрипты идемпотентными, то не придётся каждый раз переналивать сервер с нуля при отладке в случае ошибки.



## Как работает

Конфигурация сервера описывается на шаблонизированом `bash'е`.
Bash расширен директивой `render`. Например, `render /etc/ssh/sshd_confg`.
Это означает, что в bash-скрипте `render` будет заменён на `cat > /etc/ssh/sshd_config << 'EOF`.
Конфиг тоже будет шаблонизирован параметрами из yaml'а.

Для того, чтобы одни и те же скрипты и конфиги могли использоваться для настройки разных систем
`lilac` поддерживает множественные `--values values.yaml`. Это позволяет заменять часть параметров,
варьирующихся между серверами, виртуальными машинами и т.д.

---

Например, на серверах различаются версии запускаемого контейнера.

`values.node-1.example.com.yaml`
```yaml
image: "postgresql:17"
```

`values.node-2.example.com.yaml`
```yaml
image: "postgresql:18"
```

Для предпросмотра скрипта выполните.

```bash
lilac render --values values.yaml --values values.node-1.example.com.yaml --script script.sh
```

Для применения настроек выполните.

```bash
lilac apply remote --host node-1.example.com.yaml --values values.yaml --values values.node-1.example.com.yaml --script script.sh
lilac apply remote --host node-2.example.com.yaml --values values.yaml --values values.node-2.example.com.yaml --script script.sh
```
