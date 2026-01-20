#project name :ACS
PROJECT(acs)

CMAKE_MINIMUM_REQUIRED(VERSION 2.8)

ADD_SUBDIRECTORY(acsdb)
ADD_SUBDIRECTORY(drivenadapter)
ADD_SUBDIRECTORY(acssharemgnt)
ADD_SUBDIRECTORY(acsprocessor)
ADD_SUBDIRECTORY(acsdatastore)
ADD_SUBDIRECTORY(teacpserver)
ADD_SUBDIRECTORY(eachttpserver)
ADD_SUBDIRECTORY(evfsacs)
ADD_SUBDIRECTORY(eacpdaemon)
ADD_SUBDIRECTORY(entrydocacs)
