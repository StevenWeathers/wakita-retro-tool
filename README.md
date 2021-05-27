![](https://github.com/StevenWeathers/wakita-retro-tool/workflows/Go/badge.svg)
![](https://github.com/StevenWeathers/wakita-retro-tool/workflows/Node.js%20CI/badge.svg)
![](https://github.com/StevenWeathers/wakita-retro-tool/workflows/Docker/badge.svg)
![](https://img.shields.io/docker/cloud/build/stevenweathers/wakita-retro-tool.svg)
[![](https://img.shields.io/docker/pulls/stevenweathers/wakita-retro-tool.svg)](https://hub.docker.com/r/stevenweathers/wakita-retro-tool)
[![](https://goreportcard.com/badge/github.com/stevenweathers/wakita-retro-tool)](https://goreportcard.com/report/github.com/stevenweathers/wakita-retro-tool)

# Wakita Sprint Retrospective Tool

Wakita is an open source agile sprint retrospective tool to make sprint retro's more collaborative in the remote work era.

### **Uses WebSockets and [Svelte](https://svelte.dev/) frontend framework for a truly Reactive UI experience**

![image](https://user-images.githubusercontent.com/846933/116486097-cfb1ef80-a85a-11eb-9769-96a9b5737ad9.png)

# Running in production

## Use latest docker image

```
docker pull stevenweathers/wakita-retro-tool
```

## Use latest released binary

[![](https://img.shields.io/github/v/release/stevenweathers/wakita-retro-tool?include_prereleases)](https://github.com/StevenWeathers/wakita-retro-tool/releases/latest)

# Configuration
Wakita may be configured through environment variables or via a yaml file `config.yaml`
located in one of:

* `/etc/wakita/`
* `$HOME/.config/wakita/`
* Current working directory

### Example yaml configuration file

```
http:
  domain: wakita.dev
db:
  host: localhost
  port: 5432
  user: thor
  pass: odinson
  name: wakita
```

## Required configuration items

For Wakita to work correctly the following configuration items are required:

| Option                     | Environment Variable | Description                                | Default Value           |
| -------------------------- | -------------------- | ------------------------------------------ | ------------------------|
| `http.domain`              | APP_DOMAIN           | The domain/base URL for this instance of Wakita.  Used for creating URLs in emails. | wakita.dev |
| `http.cookie_hashkey`      | COOKIE_HASHKEY       | Secret used to make secure cookies secure. | twister |

### Database configuration

Wakita uses a Postgres database to store all data, the following configuration options exist: 

| Option                     | Environment Variable | Description                                | Default Value           |
| -------------------------- | -------------------- | ------------------------------------------ | ------------------------|
| `db.host`                  | DB_HOST              | Database host name.                        | db |
| `db.port`                  | DB_PORT              | Database port number.                      | 5432 |
| `db.user`                  | DB_USER              | Database user id.                          | thor |
| `db.pass`                  | DB_PASS              | Database user password.                    | odinson |
| `db.name`                  | DB_NAME              | Database instance name.                    | wakita |
| `db.sslmode`               | DB_SSLMODE           | Database SSL Mode (disable, allow, prefer, require, verify-ca, verify-full). | disable |

### SMTP (Mail) server configuration

Wakita sends emails for user registration related activities, the following configuration options exist:

| Option                     | Environment Variable | Description                                | Default Value           |
| -------------------------- | -------------------- | ------------------------------------------ | ------------------------|
| `smtp.host`                | SMTP_HOST            | Smtp server hostname.                      | localhost |
| `smtp.port`                | SMTP_PORT            | Smtp server port number.                   | 25 |
| `smtp.secure`              | SMTP_SECURE          | Set to authenticate with the Smtp server.  | true |
| `smtp.identity`            | SMTP_IDENTITY        | Smtp server authorization identity.  Usually unset. | |
| `smtp.sender`              | SMTP_SENDER          | From address in emails sent by Wakita. | no-reply@wakita.dev |

## Optional configuration items

| Option                     | Environment Variable | Description                                | Default Value           |
| -------------------------- | -------------------- | ------------------------------------------ | ------------------------|
| `http.port`                | PORT                 | Which port to listen for HTTP connections. | 8080 |
| `http.secure_cookie`       | COOKIE_SECURE        | Use secure cookies or not.                 | true |
| `http.backend_cookie_name` | BACKEND_COOKIE_NAME  | The name of the backend cookie utilized for actual auth/validation | userId |
| `http.frontend_cookie_name`| FRONTEND_COOKIE_NAME | The name of the cookie utilized by the UI (purely for convenience not auth) | user |
| `http.path_prefix`         | PATH_PREFIX          | Prefix added to all application urls for shared domain use, in format of `/{prefix}` e.g. `/wakita` | |
| `analytics.enabled`        | ANALYTICS_ENABLED    | Enable/disable google analytics.           | true |
| `analytics.id`             | ANALYTICS_ID         | Google analytics identifier.               | UA-161935945-1 |
| `config.avatar_service`    | CONFIG_AVATAR_SERVICE | Avatar service used, possible values see next paragraph | goadorable |
| `config.toast_timeout`     | CONFIG_TOAST_TIMEOUT | Number of milliseconds before notifications are hidden. | 1000 |
| `config.allow_guests`     | CONFIG_ALLOW_GUESTS | Whether or not to allow guest (anonymous) users. | true |
| `config.allow_registration`     | CONFIG_ALLOW_REGISTRATION | Whether or not to allow user registration (outside Admin). | true |
| `config.default_locale`   | CONFIG_DEFAULT_LOCALE | The default locale (language) for the UI | en |
| `config.allow_external_api`    | CONFIG_ALLOW_EXTERNAL_API | Whether or not to allow External API access | false |
| `config.show_active_countries`    | CONFIG_SHOW_ACTIVE_COUNTRIES | Whether or not to show active countries on landing page | false |
| `config.cleanup_retros_days_old` | CONFIG_CLEANUP_RETROS_DAYS_OLD | How many days back to clean up old retros, e.g. retros older than 180 days. Triggered manually by Admins . | 180 |
| `config.cleanup_guests_days_old` | CONFIG_CLEANUP_GUESTS_DAYS_OLD | How many days back to clean up old guests, e.g. guests older than 180 days.  Triggered manually by Admins. | 180 |
| `auth.method`              | AUTH_METHOD         | Choose `normal` or `ldap` as authentication method.  See separate section on LDAP configuration. | normal |

## Avatar Service configuration

Use the name from table below to configure a service - if not set, `goadorable` is used. Each service provides further options which then can be configured by a user on the profile page. Once a service is configured, drop downs with the different sprites become available. The table shows all supported services and their sprites. In all cases the same ID (`ead26688-5148-4f3c-a35d-1b0117b4f2a9`) has been used creating the avatars.

| Name |           |           |           |           |           |           |           |           |           |
| ---------- | --------- | --------- | --------- | --------- | --------- | --------- | --------- | --------- | --------- |
| `goadorable` (internal)  |           |           |           |           |           |           |           |           |           |
|            | ![image](https://user-images.githubusercontent.com/846933/96212071-e4283d80-0f43-11eb-9f82-ff6c105f8b0a.png) |
| `govatar` (internal) | male | female |  |  |
|            | ![image](https://user-images.githubusercontent.com/846933/96212029-ce1a7d00-0f43-11eb-9e53-8ca13ba9d4b1.png) | ![image](https://user-images.githubusercontent.com/846933/96212031-ceb31380-0f43-11eb-832b-b02c275317a5.png) |  |  |
| `dicebear` | male | female | human | identicon | bottts | avataaars | jdenticon | gridy | code |
|            | ![image](https://avatars.dicebear.com/api/male/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) | ![image](https://avatars.dicebear.com/api/female/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) | ![image](https://avatars.dicebear.com/api/human/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) | ![image](https://avatars.dicebear.com/api/identicon/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) | ![image](https://avatars.dicebear.com/api/bottts/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) | ![image](https://avatars.dicebear.com/api/avataaars/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) | ![image](https://avatars.dicebear.com/api/jdenticon/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) | ![image](https://avatars.dicebear.com/api/gridy/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) | ![image](https://avatars.dicebear.com/api/code/ead26688-5148-4f3c-a35d-1b0117b4f2a9.svg?w=48) |
| `gravatar` | mp | identicon | monsterid | wavatar | retro | robohash | | | |
|            | ![image](https://gravatar.com/avatar/ead26688-5148-4f3c-a35d-1b0117b4f2a9?s=48&d=mp&r=g) | ![image](https://gravatar.com/avatar/ead26688-5148-4f3c-a35d-1b0117b4f2a9?s=48&d=identicon&r=g) | ![image](https://gravatar.com/avatar/ead26688-5148-4f3c-a35d-1b0117b4f2a9?s=48&d=monsterid&r=g) | ![image](https://gravatar.com/avatar/ead26688-5148-4f3c-a35d-1b0117b4f2a9?s=48&d=wavatar&r=g) | ![image](https://gravatar.com/avatar/ead26688-5148-4f3c-a35d-1b0117b4f2a9?s=48&d=retro&r=g) | ![image](https://gravatar.com/avatar/ead26688-5148-4f3c-a35d-1b0117b4f2a9?s=48&d=robohash&r=g) | | | |
| `robohash` | set1 | set2 | set3 | set4 |
|            | ![image](https://robohash.org/ead26688-5148-4f3c-a35d-1b0117b4f2a9.png?set=set1&size=48x48) | ![image](https://robohash.org/ead26688-5148-4f3c-a35d-1b0117b4f2a9.png?set=set2&size=48x48) | ![image](https://robohash.org/ead26688-5148-4f3c-a35d-1b0117b4f2a9.png?set=set3&size=48x48) | ![image](https://robohash.org/ead26688-5148-4f3c-a35d-1b0117b4f2a9.png?set=set4&size=48x48) |

## LDAP Configuration

If `auth.method` is set to `ldap`, then the Create Account function is disabled and authentication
is done using LDAP.  If the LDAP server authenticates a new user successfully, the user 
profile is automatically generated.

The following configuration options are specific to the LDAP authentication method:

| Option                      | Environment Variable | Description                                                        |
| --------------------------- | -------------------- | ------------------------------------------------------------------ |
| `auth.ldap.url`             | AUTH_LDAP_URL        | URL to LDAP server, typically `ldap://host:port`                   |
| `auth.ldap.use_tls`         | AUTH_LDAP_USE_TLS    | Create a TLS connection after establishing the initial connection. |
| `auth.ldap.bindname`        | AUTH_LDAP_BINDNAME   | Bind name / bind DN for connecting to LDAP.  Leave empty for no authentication. |
| `auth.ldap.bindpass`        | AUTH_LDAP_BINDPASS   | Password for the bind.                                             |
| `auth.ldap.basedn`          | AUTH_LDAP_BASEDN     | Base DN for the search for the user.                               |
| `auth.ldap.filter`          | AUTH_LDAP_FILTER     | Filter for searching for the user's login id.  See below.          |
| `auth.ldap.mail_attr`       | AUTH_LDAP_MAIL_ATTR  | The LDAP property containing the user's emil address.              |
| `auth.ldap.cn_attr`         | AUTH_LDAP_CN_ATTR    | The LDAP property containing the user's name.                      |

The default `filter` is `(&(objectClass=posixAccount)(mail=%s))`.  The filter must include a `%s` that will be replaced by the user's login id.
The `mail_attr` configuration option must point to the LDAP attribute containing the user's email address.  The default is `mail`. 
The `cn_attr` configuration option must point to the LDAP attribute containing the user's full name.  The default is `cn`.

On Linux, the parameters may be tested on the command line:

```
ldapsearch -H auth.ldap.url [-Z] -x [-D auth.ldap.bindname -W] -b auth.ldap.basedn 'auth.ldap.filter' dn auth.ldap.mail auth.ldap.cn
```

The `-Z` is only used if `auth.ldap.use_tls` is set, the `-D` and `-W` parameter is only used if `auth.ldap.bindname` is set.

# Developing

## Building and running with Docker (preferred solution)

### Using Docker Compose

```
docker-compose up --build
```

### Using Docker without Compose

This solution will require you to pass environment variables or setup the config file, as well as setup and manage the DB yourself.

```
docker build ./ -f ./build/Dockerfile -t wakita:latest
docker run --publish 8080:8080 --name wakita wakita:latest
```

## Building

To run without docker you will need to first build, then setup the postgres DB,
and pass the user, pass, name, host, and port to the application as environment variables

```
DB_HOST=
DB_PORT=
DB_USER=
DB_PASS=
DB_NAME=
```

### Install dependencies
```
go get
npm install
```

## Build with Make
```
make build
```
### OR manual steps

### Build static assets
```
npm run build
```

### Build for current OS
```
go build
```

## Running with Watch (uses webapp dist files on OS instead of embedded)
```
npm run autobuild
make dev-go
```

# Run Locally

Run the server and visit [http://localhost:8080](http://localhost:8080)

# Adding new Locale's
Using svelte-i18n **Wakita** now supports Locale selection on the UI (Default en-US)

Adding new locale's involves just a couple of steps.

1. First add the locale dictionary json file in ```webapp/public/lang/``` by copying the en.json and just changing the values of all keys
1. Second, the locale will need to be added to the locales list used by switcher component in ```webapp/config.js``` ```locales``` object

# Donations

For those who would like to donate a small amount for my efforts or monthly hosting costs of Wakita.dev I accept paypal.

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://paypal.me/smweathers?locale.x=en_US)
