.PHONY: sync-agents-md dev-backend dev-frontend dev test

# 将仓库中所有的 AGENTS.md 符号链接到同目录的 CLAUDE.md，供 Claude Code 使用
sync-agents-md:
	@find . -name 'AGENTS.md' | while read f; do \
		dir=$$(dirname "$$f"); \
		ln -sf AGENTS.md "$$dir/CLAUDE.md"; \
		echo "linked $$dir/AGENTS.md -> $$dir/CLAUDE.md"; \
	done

# 启动后端开发服务（监听 :8080）
dev-backend:
	cd backend && go run main.go

# 启动前端开发服务（监听 :5173）
dev-frontend:
	cd frontend && pnpm dev

# 同时启动前后端（需要支持 & 的 shell，推荐分两个终端用 dev-backend / dev-frontend）
dev:
	@echo "请分别在两个终端运行："
	@echo "  make dev-backend"
	@echo "  make dev-frontend"

# 运行后端单元测试
test:
	cd backend && go test ./...
