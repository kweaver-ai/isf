#包含prep.cmake(定义编译时的选项)
INCLUDE($ENV{CMAKE}/common/prep.cmake)

#设置版本号，和向下兼容性
CMAKE_MINIMUM_REQUIRED(VERSION 2.8)
CMAKE_POLICY(VERSION 2.8)

#project name :EACHTTPSERVER
PROJECT(EACHTTPSERVER)
SET(NEED_XPCOM TRUE)
SET(NEED_CFL TRUE)
SET(NEED_SYSWRAP TRUE)
SET(NEED_NCUTIL TRUE)
SET(NEED_BOOST TRUE)
SET(NEED_EDATAAPI TRUE)
SET(NEED_LIBEVENT TRUE)
SET(ENABLE_C++11 TRUE)
SET(NEED_EHTTPSERVER TRUE)
SET(NEED_EHTTPCLIENT TRUE)
SET(NEED_ETHRIFT TRUE)
SET(NEED_OPENSSL TRUE)
SET(NEED_ERRCODECONVERTER TRUE)
SET(NEED_BOOST_REGEX TRUE)
SET(NEED_THRIFTAPI TRUE)
SET(NEED_BRPC TRUE)

SET(ERRGEN_FILES $ENV{SHARESERVER}/src/eachttpserver/err/eachttpservererr.csv)

#target name:eachttpserver
SET(TARGET_NAME eachttpserver)

#TARGET_TYPE
#LIBRARY | ARCHIVE | RUNTIME | COMPONENT
#动态库  | 静态库  | 可执行     |  组件

SET(TARGET_TYPE LIBRARY)

SET(CUSTOM_IDL_PATH $ENV{SHARESERVER}/src/acsprocessor/public)

# 依赖的idl文件
SET(CUSTOM_IDL_FILES $ENV{DATAENGINEIOC}/evfsioc/ncIEVFSNameIOC.idl)

SET(PROTO_FILES $ENV{SHARESERVER}/src/eachttpserver/proto/eachttpserver.proto
                $ENV{SHARESERVER}/src/eachttpserver/proto/eacinnerhttpserver.proto)

SET(CUSTOM_IDL_FILES )

INCLUDE_DIRECTORIES($ENV{DATAENGINEIOC})
SOURCE_DIRS(auth
            share
            message
            sysconfig
            third
            third/org
            third/dep
            third/user
            third/message
            protocol)

ADD_DEFINITIONS("-I$ENV{SHARESERVER}/src")
ADD_DEFINITIONS("-I$ENV{SHARESERVER}")

ADD_DEFINITIONS("-DHAVE_INTTYPES_H")
ADD_DEFINITIONS("-DHAVE_NETINET_IN_H")

INCLUDE($ENV{CMAKE}/common/target.cmake)
