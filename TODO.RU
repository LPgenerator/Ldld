* Пофиксить выдачу IP адресов в апи, потому что ключ всегда IP, и второй перезаписывает первый
* Фиксировать IP адреса с учетом кол-ва интерфейсов
* Передача аргументов в шаблон (поможет для lvm бекенда в плане указания веса fs | lxc-create --template ubuntu --name c1 --bdev loop -- --release trusty)

* Довести до ума команду Migrate
* Backends support: lvm/overlayfs
* Md5sum hashes (on commits, and /Release.json files)
* Commit: hash/message/date/author
* Reset base repository and Rollback on client
* Compression support for images: gzip/xz
* Tests + Coverage
* Multiple repository support such as git with .git directory on project directory
* Run provision script
* Веб-морда
* Providers support: freebsd/solaris
* LXC/OVZ support
* Добавить возможность указать all вместо имени контейнера all (ldld stop all, ldld start all, ldld autostart all 1, ldld pull all)
* Добавить возможность указать через запятую контейнеры (ldld stop web1,web2..n)
* Добавить возможность работы в bash (ldld stop {web1,web2})
* Добавить поддержку bash-completion

Конечный смысл проекта:
0. Создаем git репозиторий
1. Пишем Bash/Ansible скрипты
2. Добавляем файл .ldl (указываем неймспейс, бустрап файл)
3. Разворачиваем репозиторий на нужном хосте
4. Запускаем provision
5. Выполняем commit и отправляем образ в репозиторий используя push
6. На клиенте Делаем пул нужных образов запускаем/обновляем то что необходимо
