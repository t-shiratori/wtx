package config

const DefaultConfigTOML = `
[worktree]
root_dir = ""
default_base_branch = "main"

[copy]
patterns = []

# For renaming files (optional)
# [[copy.files]]
# from = ".env.example"
# to = ".env"

[hooks]
pre_create = []
post_create = []
post_copy = []
`
