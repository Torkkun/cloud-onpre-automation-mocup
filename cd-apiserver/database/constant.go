package database

const (
	AnsibleBasePrefix          = "/ansible/dynamic/"
	etcdIpPrefix               = AnsibleBasePrefix + "ip/"
	etcdHostVarsPrefix         = AnsibleBasePrefix + "hostvars/"
	etcdPlaybookTemplatePrefix = AnsibleBasePrefix + "playbook/template/"
)
