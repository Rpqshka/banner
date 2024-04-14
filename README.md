# Описание

Необходимо реализовать сервис, который позволяет показывать пользователям баннеры, в зависимости от требуемой фичи и тега пользователя, а также управлять баннерами и связанными с ними тегами и фичами.

После того, как я прочитал задания я первым делом решил сделать систему аунтетификации пользователя, тк в дальнейшем это поможет как и в решении этой задачи, так и при масштабировании данного сервиса.

Пользователь будет отправлять запросы на API используя jwt-токены, в которых мы будем хранить информацию о нем, в том числе и его роль (admin/user).

Проект разбит на 3 слоя:

* handler - обработчик API;
* repository - работа с данными;
* service - бизнес-логика


# Установка проекта

Clone repo

```bash
  git clone https://github.com/Rpqshka/banner.git

  cd banner
```

Build

```bash
  make build
```


Run app
```bash
  make run
```

Restart app

```bash
  make restart
```

Run tests

```bash
  make test
```
## Usage/Examples

[![postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/24093475-bb87d8b6-1303-4380-b2aa-af26b47c66d5?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D24093475-bb87d8b6-1303-4380-b2aa-af26b47c66d5%26entityType%3Dcollection%26workspaceId%3D782890c1-b8de-44f1-989c-8fe0dcfcd622/)

```
https://app.getpostman.com/run-collection/24093475-bb87d8b6-1303-4380-b2aa-af26b47c66d5?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D24093475-bb87d8b6-1303-4380-b2aa-af26b47c66d5%26entityType%3Dcollection%26workspaceId%3D782890c1-b8de-44f1-989c-8fe0dcfcd622
```
## Tests

Для тестов интеграционных тестов создаются отдельный сервер и база данных, на которых используются тестовые данные для проверки методов поиска баннера.

Перед началом теста регистрируются 2 новых пользователя с ролями admin и user соответвсенно, после чего выполняется логин для получения jwt-токенов, с помощью которых мы и сможем определить роль пользователя, отправившего запрос.
После того как jwt-токены получены, записываем информацию о трех баннерах в таблицу banners нашей тестовой бд. Далее с помощью 4 тест кейсов проверяем правильность работы метода.

TestGetInactiveBannerByUser - кейс проверяет возможность получения скрытого баннера обычным пользователем

TestGetInactiveBannerByAdmin - кейс проверяет возможность получения скрытого баннера админом

TestGetBannerByUser - кейс проверяет получения баннера по введенным tag_id и feature_id для обычного пользователя

TestGetBannerByAdmin - кейс проверяет получения баннера по введенным tag_id и feature_id для админа

```bash
make test
```

Результат теста:

```
--- PASS: TestBannerSuite (4.25s)
    --- PASS: TestBannerSuite/TestGetBannerByAdmin (0.00s)
    --- PASS: TestBannerSuite/TestGetBannerByUser (0.00s)
    --- PASS: TestBannerSuite/TestGetInactiveBannerByAdmin (0.00s)
    --- PASS: TestBannerSuite/TestGetInactiveBannerByUser (0.00s)
PASS
ok      banner/tests    4.313s
```
## Контакты

Telegram : @rpqshka

email: rpqshka@gmail.com

