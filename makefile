##       ##  ########  ##   ##   #######  ######
####   ####  ##    ##  ##  ##    ##       ##   ##
###########  ########  ######    #####    ######
##  ###  ##  ##    ##  ##   ##   ##       ##   ##
##       ##  ##    ##  ##    ##  #######  ##    ##

##################################################
# https://github.com/janderland/Maker



# GLOBAL CONFIG
###############

# Directory where all files generated by the build go.
buildDir ?= _build


# C++ Compiler Executable
# This executable is used to convert C++ source files into
# C++ object files.
cxx ?= g++


# C++ Compiler Flags
# cxxFlags     : used in every invocation of the compiler.
# cxxFlagsComp : used only during compilation.
# cxxFlagsLink : used only during linking.
#
# NOTE: Flags responsible for header-dependency information
# are already included by this makefile. They shouldn't be
# added here.
cxxFlags ?=
cxxFlagsComp ?=
cxxFlagsLink ?=


# Path to the ar utility.
ar ?= ar


# Flags passed to the ar utility.
arFlags ?= rcs


# Path to the rm utility.
rm ?=/bin/rm


# File extensions for the input files
sourceExt ?= cpp


# File extensions for the generated executables and
# static libraries.
execExt ?= out
slibExt ?= lib


# Set this variable to a non-empty string to turn off
# verbose output.
verbose ?= true



# GLOBAL PROCESSING
###################

# Let included makefiles know this is Maker.
isMaker=true


# Empty the .SUFFIXES variable to turn off
# almost all the builtin rules.
.SUFFIXES:


# If "verbose" is empty, then don't print the
# command being invoked by make.
ifeq ($(verbose),)
.SILENT:
endif



# RULE MACROS
#############

# Define Static Lib Rule
# 1 - Output file
# 2 - Input files
# 3 - Compile Flags
# 4 - Package Flags
slibRule=$(eval $(call slibRuleTempl,$1,$2,$3,$4))
define slibRuleTempl
$1: cxxFlagsCompExtra = $3
$1: $2
	@echo "Packaging $$@"
	mkdir -p $$(dir $$@)
	$$(ar) $$(arFlags) $4 $$@ $$^
endef


# Define Executable Rule
# 1 - Output file
# 2 - Input files
# 3 - Compile Flags
# 4 - Link Flags
execRule=$(eval $(call execRuleTempl,$1,$2,$3,$4))
define execRuleTempl
$1: cxxFlagsCompExtra = $3
$1: $2
	@echo "Linking $$@"
	mkdir -p $$(dir $$@)
	$$(cxx) $$(cxxFlags) $$(cxxFlagsLink) $4 $$^ -o $$@
endef


# Define Alias Rule
# 1 - Alias Name
# 2 - Real Name
aliasRule=$(eval $(call aliasRuleTempl,$1,$2))
define aliasRuleTempl
.PHONY: $1
$1: $2
endef


# Define Run Rule
# 1 - Rule Name
# 2 - Executable
runRule=$(eval $(call runRuleTempl,$1,$2))
define runRuleTempl
.PHONY: $1
$1: $2
	@echo "Running $$<"
	$$<
endef



# LOGGING MACROS
################

# Log Debug Message
# 1 - Message
debug=$(if $(verbose),$(info $1))



# METADATA MACROS
#################

# Module Path
# 1 - Module Makefile Path
modulePath=$(patsubst %/,%,$(dir $1))


# Check Path For 'At' Symbol ('@' isn't allowed in Maker paths)
# 1 - Module Path
checkPathForAt=$(if $(findstring @,$1),$(error $1 contains the @ symbol))


# Runner Name
# 1 - Module Path
rname=run@$1


# Input Source Files
# 1 - Module Path
sources=$(shell find $1 -iname '*.$(sourceExt)')


# Output Object Files
# 1 - Module Path
objects=$(addprefix $(buildDir)/,$(patsubst %.cpp,%.obj,$(call sources,$1)))


# Output File Path
# 1 - Module Path
# 2 - Type
file=$(buildDir)/$1.$(call fileExt,$2)


# Output File Extension
# 1 - Type
fileExt=$(strip \
$(if $(filter exec,$1),$(execExt),\
$(if $(filter slib,$1),$(slibExt),\
$$$(error No file extension defined for type $1)\
))\
)


# File Rule Macro Name
# 1 - Type
fileRuleName=$(strip \
$(if $(filter exec,$1),execRule,\
$(if $(filter slib,$1),slibRule,\
$$$(error No file rule defined for type $1)\
))\
)


# Metadata Print Template
# 1 - Input Path
# 2 - Type
# 3 - Dependencies
# 4 - Compile Flags
# 5 - Link/Package Flags
define formatMetadata
Type          : $2
Dependencies  : $3
Compile Flags : $4
Linker Flags  : $5
Source Files  : $(call sources,$1)
Object Files  : $(call objects,$1)
Output File   : $(call file,$1,$2)
endef


# Define Module Metadata
# 1 - Input Path
# 2 - Type
# 3 - Dependencies
# 4 - Compile Flags
# 5 - Link/Package Flags
module=$(eval $(call moduleTempl,$(strip $1),$(strip $2),$(strip $3),$(strip $4),$(strip $5)))
define moduleTempl
$(call debug,Defining module for $1)
$(call checkPathForAt,$1)
$(call debug,$(call formatMetadata,$1,$2,$3,$4,$5))
$(call debug,)

$(1)Type=$2
$(1)Deps=$3
$(1)CFlags=$4
$(1)LFlags=$5
targets+=$1
endef


# Dependency Files
# 1 - Dependency input paths
depFiles=$(foreach v,$1,$(call file,$v,$($(v)Type)))


# Define Module Rules
# 1 - Input Path
rules=$(eval $(call rulesTempl,$1,$($(1)Type),$($(1)Deps),$($(1)CFlags),$($(1)LFlags)))
define rulesTempl
$(call $(call fileRuleName,$2),$(call file,$1,$2),$(call objects,$1) $(call depFiles,$3),$4,$5)
$(call aliasRule,$1,$(call file,$1,$2))
$(if $(filter exec,$2),$(call runRule,$(call rname,$1),$(call file,$1,$2)))
endef



# HEADER DEP MACROS
###################

headerDepFlags=-MMD -MF $(basename $@).dep

headerDepFiles=$(shell if [ -d $(buildDir) ]; then find $(buildDir) -iname *.dep; fi)



# LOAD METADATA
###############

$(call debug,Configuration)
$(call debug,=============)

moduleFiles=$(shell find . -iname makefile -mindepth 1 | cut -c3-)
$(foreach f,$(moduleFiles),\
 $(eval include $f)\
 $(call module,$(call modulePath,$f),$(moduleType),$(moduleDeps),$(moduleCompFlags),$(moduleLinkFlags)))



# STATIC RULES
##############

.PHONY: all clean

all: $(targets)

$(buildDir)/%.obj: %.cpp
	$(info,Compiling $<)
	mkdir -p $(dir $@)
	$(cxx) -c $(headerDepFlags) $(cxxFlags) $(cxxFlagsComp) $(cxxFlagsCompExtra) $< -o $@

clean:
	$(info,Cleaning)
	$(rm) -rf $(buildDir)



# DYNAMIC RULES
###############

ifndef targets
$(error No module definitions found)
else
$(call debug,Targets: $(targets))
$(call debug,)

$(foreach t,$(targets),$(call rules,$t))

include $(headerDepFiles)

$(call debug,Compilation)
$(call debug,===========)
endif
