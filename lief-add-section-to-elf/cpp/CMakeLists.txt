
# LIEF static library
set(LIB_LIEF "/usr/lib/libLIEF.a")


add_executable(main main.cpp)
target_include_directories(main PUBLIC "/usr/include")

# Link the executable with LIEF
target_link_libraries(main PUBLIC ${LIB_LIEF})
add_dependencies(main LIEF)
