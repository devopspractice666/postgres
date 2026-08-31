ansible-playbook playbooks/deploy.yml                             Полный деплой с копией конфига
playbooks/deploy.yml --tags start                                 Просто запуск контейнеров 
ansible-playbook playbooks/deploy.yml --tags stop                 Остановка контейнеров