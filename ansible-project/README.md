# Запустить всё (БД + приложение)
ansible-playbook playbooks/site.yml --ask-vault-pass

# Только БД
ansible-playbook playbooks/db.yml --ask-vault-pass

# Только приложение
ansible-playbook playbooks/myapp.yml --ask-vault-pass

# Остановить всё
ansible-playbook playbooks/site.yml --tags stop --ask-vault-pass

# Только БД
ansible-playbook playbooks/db.yml --tags stop --ask-vault-pass

# Только приложение
ansible-playbook playbooks/myapp.yml --tags stop --ask-vault-pass