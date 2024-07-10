RED   = \033[31m
RESET = \033[0m

define message
    @echo -e "$(RED)$(1)$(RESET)"
endef

conf2dev:
	$(call message,<1> UPDATE parkour)
	cd parkour; git reset --hard; git clean -df; git checkout dev; git pull
	$(call message,<2> UPDATE paoku_res)
	cd paoku_res; git reset --hard; git clean -df; git svn rebase
	$(call message,<3> COPY paoku_res/DEV/*.xlsx --> parkour/conf/xlsx/)
	/usr/bin/cp -f paoku_res/DEV/* parkour/conf/xlsx
	$(call message,<4> EXPORT parkour/conf/xlsx/)
	cd parkour && ./conf/conf.exe && gofmt -w conf && goimports -w conf
	$(call message,<4> COMMIT)
	cd parkour && git add . && git commit -m "配置文件更新：$$(git status -s conf/json)"


conf2test:
	$(call message,<1> UPDATE parkour)
	cd parkour; git reset --hard; git clean -df; git checkout build_test; git pull
	$(call message,<2> UPDATE paoku_res)
	cd paoku_res; git reset --hard; git clean -df; git svn rebase
	$(call message,<3> COPY paoku_res/DEV/*.xlsx --> parkour/conf/xlsx/)
	/usr/bin/cp -f paoku_res/DEV/* parkour/conf/xlsx
	$(call message,<4> EXPORT parkour/conf/xlsx/)
	cd parkour && ./conf/conf.exe && gofmt -w conf && goimports -w conf
	$(call message,<4> COMMIT)
	cd parkour && git add . && git commit -m "配置文件更新：$$(git status -s conf/json)"


init:
	@(! test -e parkour)   || (echo -e "$(RED)ERROR:$(RESET) parkour existed"   && exit -1)
	@(! test -e paoku_res) || (echo -e "$(RED)ERROR:$(RESET) paoku_res existed" && exit -1)
	$(call message,<1> CLONE 服务端仓库)
	git clone http://gitlab.mrmb.xyz:8099/paoku_server/parkour.git
	$(call message,<2> CLONE 配置文件仓库)
	git svn clone -r HEAD https://pc-42server/svn/paoku/res paoku_res
	cp conf.exe parkour/conf/

switchtodev:
	cd parkour; [[ $(git branch --show-current) == dev ]] || git checkout dev

test:
	$(call message,<1> CLONE 服务端仓库)
