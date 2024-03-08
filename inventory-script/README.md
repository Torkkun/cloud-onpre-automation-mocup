# Ansible用のDynamic Inventory Script
ストアには、etcdで分散可能にする
ansibleではこのetcdからローカルのホスト情報などを保存し、CRUDを動的にAPIサーバーから管理できるようにする。
プロトタイプなので、Scriptで運用する。
ansible Dynamic inventory Document
https://docs.ansible.com/ansible/latest/dev_guide/developing_inventory.html
最近は、インベントリモジュールを勧めている？
なお、テスト用

## debug
go build
(このままであればtestinventoryscという実行ファイルが作成される)
ansible-inventory -i ./main --list

playbookも可能（テストデータ）
