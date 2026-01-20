
#project name :ACSUT
PROJECT(ACSUT)

CMAKE_MINIMUM_REQUIRED(VERSION 2.8)

IF(NOT WIN32)
	ADD_DEFINITIONS("-fprofile-arcs -ftest-coverage")
	LINK_LIBRARIES(gcov)
ENDIF()

ADD_SUBDIRECTORY(acssharemgntut)
ADD_SUBDIRECTORY(acsdbut)
ADD_SUBDIRECTORY(acsprocessorut)
ADD_SUBDIRECTORY(evfsacsut)
