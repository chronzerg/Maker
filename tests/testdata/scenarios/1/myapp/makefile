ifndef isMaker
$(shell make -C .. myapp)
else

moduleType      = exec
moduleDeps      = mylib
moduleCompFlags = -X{0} -wasted
moduleLinkFlags = -redflag -nope

endif
