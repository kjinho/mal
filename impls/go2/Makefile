# export GOPATH := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

#####################

# SOURCES_BASE = src/types/types.go src/readline/readline.go \
# 	       src/reader/reader.go src/printer/printer.go \
# 	       src/env/env.go src/core/core.go

#####################

# SRCS = step0_repl.go # step1_read_print.go step2_eval.go step3_env.go \
#        # step4_if_fn_do.go step5_tco.go step6_file.go step7_quote.go \
#        # step8_macros.go step9_try.go stepA_mal.go
# BINS = $(SRCS:%.go=%)

# #####################

# all: $(BINS)

# # dist: mal

# # mal: $(word $(words $(BINS)),$(BINS))
# # 	cp $< $@

# define dep_template
# $(1): $(SOURCES_BASE) src/$(1)/$(1).go
# 	go build $$@
# endef

# $(foreach b,$(BINS),$(eval $(call dep_template,$(b))))

SRCS0 = readline/readline.go
SRCS1 = reader/reader.go
step0_repl: step0/step0_repl.go $(SRCS0)
	go build step0/step0_repl.go

step1_read_print: step1/step1_read_print.go $(SRCS0) $(SRCS1)
	go build step1/step1_read_print.go

clean:
	rm -f $(BINS) mal
