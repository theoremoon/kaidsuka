set(CMAKE_VERBOSE_MAKEFILE TRUE)
find_package(retdec 4.0 REQUIRED COMPONENTS retdec)

add_executable(main main.cpp)
target_link_libraries(main PUBLIC retdec::retdec)
target_include_directories(main PUBLIC retdec::retdec)
