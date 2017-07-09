Maker
=====
Maker is a simple C++ build system implemented in GNU Make. It was created to simplify makefile configuration.

## How to install?
Simply copy the makefile from this repo in your project.

## How to configure?
Maker looks for a 'config.mkr' file in the same directory and reads it as configuration. Here is what the configuration file might look like:

```Makefile
# Project-wide compilation flags
cxxFlags = -Wall -std=c++11;

# Define an executable module.
$(call exec, source/executable_code,

    # Dependencies
    source/static_lib_code
    source/other_module,


    # Module Compilation Flags
    -I path/to/include
    -I other/include/path

    -D THIS_MACRO=1234
    -D THAT_MACRO heyo,


    # Module Linking Flags
    -l ncurses
    -l pthread
);

# Define some static library modules
$(call slib, source/static_lib_code,, -O3);
$(call slib, source/other_module);
```

Configuration is done via global configuration and target definitions:
 - Global configuration consists of several variables which are defined at the top of the Maker makefile. Their values can be changed in the Makefile or preferably redefined in the config.mkr file (as seen with the `cxxFlags` variable in the example above). For more information on these variables, see the top of the Maker makefile where the variables are well documented.

 - Target definitions are done via calls to Maker functions (as seen with the `$(call exec...)` and `$(call slib...)` statements in the example above). All targets are currently created as a result of module definition. What a module is and what the arguments to said functions mean is discussed in the following sections.

### Module Targets
With Maker, build targets are defined via module definitions. A module is a group of source files which are built together using the same compile flags, and linked or packaged into some kind of output file. Maker provides two types of module output files: executable or static library. To define a module, there is a function call for each output type:

```Makefile
# Define an executable module
$(call exec, source path, dependencies, compile flags, link flags);

# Define a static library module
$(call slib, source path, dependencies, compile flags, package flags);
```

The module function arguments, as named above, have the following meanings:
- _source path_: The path to the folder containing all the source files for this module. All '.cpp' and '.h' files in this folder and its subfolders will be included in this module. What extensions are searched for is configurable in the global configuration. Note: source paths cannot contain the '@' symbol. This is because this symbol is reserved for usage by the "run targets" (see next section).

- _dependencies_: The source paths of other module this module depends on. Executable modules can depend on any number of static library modules, causing their output files to be built (if necessary) and added to the input of the executable module's link stage. Currently, this argument isn't used by static library modules, but will be in future versions.

- _compile flags_: The flags to pass to the compiler when compiling a source file.

- _link/package flags_: The flags passed to the utility (g++ or ar) during the linking or packaging stage, for executable or static library modules respectively.

To build a given module, Maker defines a target named after the module's source path. For instance, given the configuration at the start of this readme, one could invoke make as follows:

```bash
make source/executable_code
make source/static_lib_code
make source/other_module
```

### Run Targets
Whenever an executable module is defined, a run target for said module is also defined. When executed, this target runs the module's output executable, building it first if it's out of date. The target takes the form of `run@path/to/module`. For instance, given the configuration at the start of this readme, one could invoke make as follows:

```bash
make run@source/executable_code
```
