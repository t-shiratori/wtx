package config

const DefaultConfigTOML = `
[worktree]
root_dir = ""
default_base_branch = "main"

[[copy]]
from = ""
to = ""

[hooks]
pre_create = []
post_create = []
post_copy = []
`
