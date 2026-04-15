.PHONY: sync-agents-md

# 将仓库中所有的 AGENTS.md 符号链接到同目录的 CLAUDE.md，供 Claude Code 使用
sync-agents-md:
	@find . -name 'AGENTS.md' | while read f; do \
		dir=$$(dirname "$$f"); \
		ln -sf AGENTS.md "$$dir/CLAUDE.md"; \
		echo "linked $$dir/AGENTS.md -> $$dir/CLAUDE.md"; \
	done
