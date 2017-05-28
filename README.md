Maker
=====

## What is it?
Maker is a simple C++ build system implemented in GNU Make. It was created to simplify makefile configuration.

## How to install?
Simply copy the makefile from this repo in your project.

## How to configure?
Maker scans for any ".mk" files in the same directory and reads them as configuration. Here is what a configuration file might look like:

```Makefile
# Define an executable module.
$(call exec, source/executable_code,   \
                                       \
    # Dependencies                     \
    source/static_lib_code             \
    source/other_module,               \
                                       \
    # Compile Flags                    \
    $(call includes,                   \
        path/to/include                \
        other/include/path             \
    )                                  \
    $(call define, THIS_MACRO, 1234)   \
    $(call define, THAT_MACRO, heyo),  \
                                       \
    # Link Flags                       \
    $(call libraries,                  \
        ncurses,                       \
        pthread                        \
    )                                  \
)

# Define some static library modules
$(call slib, source/static_lib_code,,-O3)
$(call slib, source/other_module)
```

### Defining Modules
With Maker, targets are defined via module definitions. A module is a group of source files which are built together using the same compile flags, and linked or packaged into some kind of output file. Maker provides two types of module output file: executable or static library. To define a module, there is a function call for each output type:

```Makefile
# Define an executable module
$(call exec, source path, dependencies, compile flags, link flags)

# Define a static library module
$(call slib, source path, dependencies, compile flags, package flags)
```

Modules are configured via global configuration and arguments passed into the module functions. The global configuration is located at the top of the makefile and is well documented with comments. The module function arguments are documented below the following sections.

#### source path
The path to the folder containing all the source files for this module. All '.cpp' and '.h' files in this folder and its subfolders will be included in this module. What extensions are searched for is configurable in the global configuration.

#### dependencies
The source paths of other module this module depends on. Executable modules can depend on any number of static library modules, causing their output files to be built (if necessary) and added to the input of the executable module's link stage.

Currently, this argument isn't used by static library modules, but will be in future versions.

#### compile flags
The flags to pass to the compiler when compiling a source file. As shown in the example, there are several functions provided for semantics.

```Makefile
# -Ipath/to/include -Iother/include/path
$(call includes,       \
    path/to/include    \
    other/include/path \
)

# -DMY_NUMBER=10
$(call define, MY_NUMBER, 10)
```

#### link/package flags
The flags passed to the compiler during linking or packaging stage, for executable or static library modules respectively. As shown in the example, there are several functions provided for semantics.

```Makefile
# -lncurses -lpthread
$(call libraries, \
    ncurses,      \
    pthread       \
)
```

